package assistant

import (
	"fmt"
	"strings"
)

// TriageRequest 分诊请求
type TriageRequest struct {
	Symptom string `json:"symptom" binding:"required"`
}

// TriageResponse 分诊响应
type TriageResponse struct {
	Symptom      string       `json:"symptom"`
	Advice       string       `json:"advice"`        // DeepSeek 生成或降级的文本建议
	Depts        []TriageDept `json:"depts"`          // 匹配的科室列表
	KnowledgeRef string       `json:"knowledgeRef"`   // 匹配的知识库分类
	Urgency      string       `json:"urgency"`        // high/medium/low
	Mode         string       `json:"mode"`           // "llm" | "keyword"
	Disclaimer   string       `json:"disclaimer"`
}

// TriageDept 推荐科室
type TriageDept struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Service 就诊助手业务服务
type Service struct {
	retriever  *Retriever
	llm        *LLMClient
	deptRes    *DepartmentResolver
	cfg        *Config
}

// NewService 创建就诊助手服务
func NewService(cfg *Config) *Service {
	return &Service{
		retriever: NewRetriever(cfg),
		llm:       NewLLMClient(cfg),
		deptRes:   NewDepartmentResolver("http://his-user:8082"), // Docker 内部 hostname
		cfg:       cfg,
	}
}

// triageContext 检索阶段收集的上下文
type triageContext struct {
	deptTypes    []string
	category     string
	urgency      string
	contextParts []string
}

// Triage 执行分诊：检索 → 科室匹配 → LLM 生成建议
func (s *Service) Triage(req *TriageRequest) (*TriageResponse, error) {
	symptom := strings.TrimSpace(req.Symptom)
	if symptom == "" {
		return nil, fmt.Errorf("症状描述不能为空")
	}

	// 1. RAG 检索 + 收集上下文
	results := s.retriever.Search(symptom)
	ctx := s.collectTriageContext(results)

	// 2. 匹配本院科室
	matchedDepts, err := s.deptRes.MatchDeptTypes(ctx.deptTypes)
	if err != nil {
		matchedDepts = fallbackDepts(ctx.deptTypes)
	}

	// 3. 生成建议
	deptNames := extractDeptNames(matchedDepts)
	kContext := strings.Join(ctx.contextParts, "\n")
	advice, mode := s.generateAdvice(symptom, kContext, ctx.category, ctx.urgency, deptNames)

	// 4. 构建响应
	triageDepts := make([]TriageDept, len(matchedDepts))
	for i, d := range matchedDepts {
		triageDepts[i] = TriageDept{ID: d.ID, Name: d.Name}
	}

	return &TriageResponse{
		Symptom:      symptom,
		Advice:       advice,
		Depts:        triageDepts,
		KnowledgeRef: ctx.category,
		Urgency:      ctx.urgency,
		Mode:         mode,
		Disclaimer:   "⚠️ 本建议仅供参考，不能替代专业医疗诊断。如症状严重请及时就医。",
	}, nil
}

// collectTriageContext 从检索结果提取科室类型、分类、紧急度
func (s *Service) collectTriageContext(results []SearchResult) triageContext {
	ctx := triageContext{
		deptTypes: []string{"内科"},
		category:  "未知（请进一步描述症状）",
		urgency:   "medium",
	}
	if len(results) == 0 {
		return ctx
	}

	ctx.deptTypes = nil
	seenTypes := make(map[string]bool)
	seenCategories := make(map[string]bool)

	for _, res := range results {
		if res.Chunk.Category != "" && !seenCategories[res.Chunk.Category] {
			ctx.contextParts = append(ctx.contextParts,
				fmt.Sprintf("- %s (紧急度: %s)", res.Chunk.Category, res.Chunk.Urgency))
			seenCategories[res.Chunk.Category] = true
		}
		for _, dt := range res.Chunk.DeptTypes {
			if !seenTypes[dt] {
				ctx.deptTypes = append(ctx.deptTypes, dt)
				seenTypes[dt] = true
			}
		}
		if ctx.urgency == "" || (res.Chunk.Urgency == "high" && ctx.urgency != "high") {
			ctx.urgency = res.Chunk.Urgency
		}
		if ctx.category == "" {
			ctx.category = res.Chunk.Category
		}
	}

	if len(ctx.deptTypes) == 0 {
		ctx.deptTypes = []string{"内科"}
		ctx.category = "未知（请进一步描述症状）"
	}
	if ctx.urgency == "" {
		ctx.urgency = "medium"
	}

	return ctx
}

// generateAdvice 根据可用性选择 LLM 或关键词模式生成建议
func (s *Service) generateAdvice(symptom, kContext, category, urgency string, deptNames []string) (string, string) {
	if !s.llm.Available() {
		return buildKeywordAdvice(symptom, category, urgency, deptNames), "keyword"
	}
	llmAdvice, err := s.llm.GenerateTriageAdvice(symptom, kContext, deptNames)
	if err == nil && llmAdvice != "" {
		return llmAdvice, "llm"
	}
	return buildKeywordAdvice(symptom, category, urgency, deptNames), "keyword"
}

// extractDeptNames 提取科室名称列表
func extractDeptNames(depts []Department) []string {
	names := make([]string, len(depts))
	for i, d := range depts {
		names[i] = d.Name
	}
	return names
}

// buildKeywordAdvice 关键词模式下生成文本建议（无需 LLM）
func buildKeywordAdvice(symptom, category, urgency string, deptNames []string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("根据您描述的「%s」，初步匹配到以下信息：\n\n", symptom))
	sb.WriteString(fmt.Sprintf("**匹配分类**：%s\n", category))

	if len(deptNames) > 0 {
		sb.WriteString(fmt.Sprintf("**推荐科室**：%s\n", strings.Join(deptNames, " / ")))
	} else {
		sb.WriteString("**推荐科室**：内科（建议进一步描述症状以获得更精准推荐）\n")
	}

	if urgency == "high" {
		sb.WriteString("\n**⚠️ 注意事项**：该症状紧急度较高，建议尽快就医或前往急诊科。\n")
	} else if urgency == "medium" {
		sb.WriteString("\n**注意事项**：建议近期挂号就诊。如症状加重请及时就医。\n")
	}

	sb.WriteString("\n**免责声明**：⚠️ 本建议仅供参考，不能替代专业医疗诊断。如症状严重请及时就医。")
	return sb.String()
}

// fallbackDepts 科室服务不可用时的兜底科室列表
func fallbackDepts(deptTypes []string) []Department {
	depts := make([]Department, 0, len(deptTypes))
	for _, dt := range deptTypes {
		depts = append(depts, Department{
			ID:   "",
			Name: dt,
		})
	}
	return depts
}
