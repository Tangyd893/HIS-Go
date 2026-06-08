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

// mergeInto 将搜索结果合并到 dest map 中，markBoth 为 true 时标记双路命中
func mergeInto(dest map[string]SearchResult, src []SearchResult, markBoth bool) {
	for _, res := range src {
		key := res.Chunk.SourceID
		existing, ok := dest[key]
		if !ok || res.Score > existing.Score {
			if markBoth && ok {
				res.MatchType = "both"
			}
			dest[key] = res
		}
	}
}

// topKResult 按分数降序排序并截取 Top-K
func topKResult(results map[string]SearchResult, k int) []SearchResult {
	if k <= 0 {
		k = 5
	}
	sorted := make([]SearchResult, 0, len(results))
	for _, res := range results {
		sorted = append(sorted, res)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score > sorted[j].Score
	})
	if len(sorted) > k {
		sorted = sorted[:k]
	}
	return sorted
}

// Search 双路检索：关键词匹配 + 语义检索（如可用），合并 Top-K
func (r *Retriever) Search(query string) []SearchResult {
	query = strings.TrimSpace(query)
	if query == "" || len(r.chunks) == 0 {
		return nil
	}

	results := make(map[string]SearchResult)

	// 关键词匹配
	mergeInto(results, r.keywordSearch(query), false)

	// 语义检索（如嵌入可用）
	if r.cfg.IsSemanticSearchAvailable() && len(r.embeddings) == len(r.chunks) {
		mergeInto(results, r.semanticSearch(query), true)
	}

	return topKResult(results, r.cfg.TopK)
}

// scoreChunk 计算 query 对单个 chunk 的关键词匹配得分（精确/分词）
func scoreChunk(chunkText, queryLower string) float64 {
	textLower := strings.ToLower(chunkText)
	if strings.Contains(textLower, queryLower) {
		return 1.0
	}
	words := strings.Fields(queryLower)
	matched := 0
	for _, w := range words {
		if len(w) >= 2 && strings.Contains(textLower, w) {
			matched++
		}
	}
	if matched == 0 || len(words) == 0 {
		return 0
	}
	return float64(matched) / float64(len(words))
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
		score := scoreChunk(r.chunks[i].Text, queryLower)
		if score > 0 {
			results = append(results, SearchResult{
				Chunk:     &r.chunks[i],
				Score:     score,
				MatchType: "keyword",
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// semanticSearch 语义检索：使用 BGE 嵌入做余弦相似度（当前用词汇重叠度近似）
func (r *Retriever) semanticSearch(query string) []SearchResult {
	queryLower := strings.ToLower(query)
	var results []SearchResult

	for i := range r.chunks {
		score := scoreChunk(r.chunks[i].Text, queryLower)
		if score > 0.1 {
			// 预留：后续 query embedding API 就绪后启用余弦相似度加权
			_ = r.embeddings
			results = append(results, SearchResult{
				Chunk:     &r.chunks[i],
				Score:     score * 0.8,
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
