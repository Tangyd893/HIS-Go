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

// Triage 执行分诊：检索 → 科室匹配 → LLM 生成建议
func (s *Service) Triage(req *TriageRequest) (*TriageResponse, error) {
	symptom := strings.TrimSpace(req.Symptom)
	if symptom == "" {
		return nil, fmt.Errorf("症状描述不能为空")
	}

	// 1. RAG 检索知识库
	results := s.retriever.Search(symptom)

	var (
		advice       string
		deptTypes    []string
		category     string
		urgency      string
		mode         string
		contextParts []string
	)

	if len(results) > 0 {
		// 从检索结果收集科室类型和分类
		seenTypes := make(map[string]bool)
		seenCategories := make(map[string]bool)
		for _, res := range results {
			if res.Chunk.Category != "" && !seenCategories[res.Chunk.Category] {
				contextParts = append(contextParts, fmt.Sprintf("- %s (紧急度: %s)", res.Chunk.Category, res.Chunk.Urgency))
				seenCategories[res.Chunk.Category] = true
			}
			for _, dt := range res.Chunk.DeptTypes {
				if !seenTypes[dt] {
					deptTypes = append(deptTypes, dt)
					seenTypes[dt] = true
				}
			}
			if urgency == "" || (res.Chunk.Urgency == "high" && urgency != "high") {
				urgency = res.Chunk.Urgency
			}
			if category == "" {
				category = res.Chunk.Category
			}
		}
	}

	if len(deptTypes) == 0 {
		// 无匹配时默认推荐内科
		deptTypes = []string{"内科"}
		category = "未知（请进一步描述症状）"
		urgency = "medium"
	}
	if urgency == "" {
		urgency = "medium"
	}

	// 2. 匹配本院科室
	matchedDepts, err := s.deptRes.MatchDeptTypes(deptTypes)
	if err != nil {
		// 科室服务不可用时降级：返回知识库推荐的科室类型作为兜底
		matchedDepts = fallbackDepts(deptTypes)
	}

	// 3. 尝试 LLM 生成建议
	kContext := strings.Join(contextParts, "\n")
	deptNames := make([]string, len(matchedDepts))
	for i, d := range matchedDepts {
		deptNames[i] = d.Name
	}

	if s.llm.Available() {
		llmAdvice, err := s.llm.GenerateTriageAdvice(symptom, kContext, deptNames)
		if err == nil && llmAdvice != "" {
			advice = llmAdvice
			mode = "llm"
		} else {
			advice = buildKeywordAdvice(symptom, category, urgency, deptNames)
			mode = "keyword"
		}
	} else {
		advice = buildKeywordAdvice(symptom, category, urgency, deptNames)
		mode = "keyword"
	}

	// 4. 构建响应
	triageDepts := make([]TriageDept, len(matchedDepts))
	for i, d := range matchedDepts {
		triageDepts[i] = TriageDept{ID: d.ID, Name: d.Name}
	}

	return &TriageResponse{
		Symptom:      symptom,
		Advice:       advice,
		Depts:        triageDepts,
		KnowledgeRef: category,
		Urgency:      urgency,
		Mode:         mode,
		Disclaimer:   "⚠️ 本建议仅供参考，不能替代专业医疗诊断。如症状严重请及时就医。",
	}, nil
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
