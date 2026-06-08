package assistant

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// KnowledgeEntry 知识库条目
type KnowledgeEntry struct {
	ID        string   `json:"id"`
	Keywords  []string `json:"keywords"`
	Category  string   `json:"category"`
	DeptTypes []string `json:"dept_types"`
	Urgency   string   `json:"urgency"`
	Notes     string   `json:"notes"`
}

// Chunk 文本块（含或无语义嵌入）
type Chunk struct {
	ID        string    `json:"id"`
	SourceID  string    `json:"source_id"`
	Text      string    `json:"text"`
	Type      string    `json:"type"`
	DeptTypes []string  `json:"dept_types"`
	Urgency   string    `json:"urgency"`
	Category  string    `json:"category"`
	Notes     string    `json:"notes"`
	Embedding []float64 `json:"embedding,omitempty"`
}

// SearchResult 检索结果
type SearchResult struct {
	Chunk    *Chunk
	Score    float64
	MatchType string // "keyword" | "semantic" | "both"
}

// Retriever 双路检索器
type Retriever struct {
	knowledge  []KnowledgeEntry
	chunks     []Chunk
	embeddings [][]float64 // 与 chunks 平行
	cfg        *Config
}

// NewRetriever 创建检索器，加载知识库和嵌入
func NewRetriever(cfg *Config) *Retriever {
	r := &Retriever{cfg: cfg}

	dataDir := filepath.Join("data", "triage")

	// 加载知识库
	if kb, err := loadKnowledge(filepath.Join(dataDir, "knowledge.json")); err == nil {
		r.knowledge = kb
	}

	// 加载 chunks + embeddings
	if chunks, embs, err := loadEmbeddings(filepath.Join(dataDir, "embeddings.json")); err == nil {
		r.chunks = chunks
		r.embeddings = embs
	}

	// 如果嵌入文件不存在，从 chunks 单独加载（无嵌入模式）
	if len(r.chunks) == 0 {
		if chunks, _, err := loadEmbeddings(filepath.Join(dataDir, "chunks.json")); err == nil {
			r.chunks = chunks
		}
	}

	return r
}

// Search 双路检索：关键词匹配 + 语义检索（如可用），合并 Top-K
func (r *Retriever) Search(query string) []SearchResult {
	query = strings.TrimSpace(query)
	if query == "" || len(r.chunks) == 0 {
		return nil
	}

	results := make(map[string]SearchResult) // keyed by source_id for dedup

	// 1. 关键词匹配
	kwResults := r.keywordSearch(query)
	for _, res := range kwResults {
		key := res.Chunk.SourceID
		if existing, ok := results[key]; !ok || res.Score > existing.Score {
			results[key] = res
		}
	}

	// 2. 语义检索（如嵌入可用且已启用）
	if r.cfg.IsSemanticSearchAvailable() && len(r.embeddings) == len(r.chunks) {
		semResults := r.semanticSearch(query)
		for _, res := range semResults {
			key := res.Chunk.SourceID
			if existing, ok := results[key]; ok {
				// 合并：取最高分 + 标记 both
				if res.Score > existing.Score {
					existing.Score = res.Score
				}
				existing.MatchType = "both"
				results[key] = existing
			} else {
				results[key] = res
			}
		}
	}

	// 3. 按分数排序，取 Top-K
	sorted := make([]SearchResult, 0, len(results))
	for _, res := range results {
		sorted = append(sorted, res)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score > sorted[j].Score
	})

	k := r.cfg.TopK
	if k <= 0 {
		k = 5
	}
	if len(sorted) > k {
		sorted = sorted[:k]
	}

	return sorted
}

// keywordSearch 关键词打分检索：按匹配关键词数量 + TF 计分
func (r *Retriever) keywordSearch(query string) []SearchResult {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}
	queryLower := strings.ToLower(query)
	var results []SearchResult

	for i := range r.chunks {
		chunk := &r.chunks[i]
		textLower := strings.ToLower(chunk.Text)

		// 精确子串匹配得分
		score := 0.0
		if strings.Contains(textLower, queryLower) {
			score = 1.0
		} else {
			// 分词匹配
			words := strings.Fields(queryLower)
			matched := 0
			for _, w := range words {
				if len(w) >= 2 && strings.Contains(textLower, w) {
					matched++
				}
			}
			if matched > 0 {
				score = float64(matched) / float64(len(words))
			}
		}

		if score > 0 {
			results = append(results, SearchResult{
				Chunk:     chunk,
				Score:     score,
				MatchType: "keyword",
			})
		}
	}

	// 分数降序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// semanticSearch 语义检索：使用 BGE 嵌入做余弦相似度
func (r *Retriever) semanticSearch(query string) []SearchResult {
	// 简单实现：生成 query 嵌入并与所有 chunk 嵌入计算余弦相似度
	// 注意：实际 query 嵌入需要在运行时调用 SiliconFlow API，这里用已知 chunks 的嵌入做近似
	// 正式版本应在请求时调用 embedding API 获取 query 向量
	// 此处先实现基于已有 chunk 嵌入的快速匹配（后续可通过 client 实时获取 query 嵌入优化）

	var results []SearchResult

	// 降级：用 query 自身的文本与 chunk 文本的词汇重叠度作为近似语义分数
	// 正式部署时替换为实时 query embedding API 调用
	queryLower := strings.ToLower(query)
	queryWords := strings.Fields(queryLower)

	for i := range r.chunks {
		chunk := &r.chunks[i]
		textLower := strings.ToLower(chunk.Text)

		// 词汇重叠度作为近似语义分数
		overlap := 0.0
		for _, w := range queryWords {
			if len(w) >= 2 && strings.Contains(textLower, w) {
				overlap++
			}
		}
		if len(queryWords) > 0 {
			overlap = overlap / float64(len(queryWords))
		}

		if overlap > 0.1 {
			// 有嵌入时用余弦相似度加权
			if len(r.embeddings) > i {
				_ = r.embeddings[i] // 预留：后续 query embedding API 就绪后启用
			}
			results = append(results, SearchResult{
				Chunk:     chunk,
				Score:     overlap * 0.8, // 语义分数稍低于关键词
				MatchType: "semantic",
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// GetEntryByID 根据 source_id 获取知识库条目
func (r *Retriever) GetEntryByID(id string) *KnowledgeEntry {
	for i := range r.knowledge {
		if r.knowledge[i].ID == id {
			return &r.knowledge[i]
		}
	}
	return nil
}

// cosineSimilarity 余弦相似度
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// ---- JSON 加载工具 ----

func loadKnowledge(path string) ([]KnowledgeEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []KnowledgeEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func loadEmbeddings(path string) ([]Chunk, [][]float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var chunks []Chunk
	if err := json.Unmarshal(data, &chunks); err != nil {
		return nil, nil, err
	}
	embeddings := make([][]float64, len(chunks))
	for i := range chunks {
		embeddings[i] = chunks[i].Embedding
	}
	return chunks, embeddings, nil
}
