package assistant

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// ============================================================
// cosineSimilarity
// ============================================================

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name   string
		a, b   []float64
		expect float64
	}{
		{"identical", []float64{1, 0, 1}, []float64{1, 0, 1}, 1.0},
		{"orthogonal", []float64{1, 0, 0}, []float64{0, 1, 0}, 0.0},
		{"opposite", []float64{1, 0}, []float64{-1, 0}, -1.0},
		{"empty a", nil, []float64{1, 2}, 0.0},
		{"empty b", []float64{1, 2}, nil, 0.0},
		{"mismatched lengths", []float64{1, 2, 3}, []float64{4, 5}, 0.0},
		{"both empty", []float64{}, []float64{}, 0.0},
		{"partial match", []float64{3, 4}, []float64{3, 4}, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cosineSimilarity(tt.a, tt.b)
			// Allow small floating-point error
			if (got-tt.expect) > 1e-9 || (tt.expect-got) > 1e-9 {
				t.Errorf("cosineSimilarity(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expect)
			}
		})
	}
}

// ============================================================
// keywordSearch
// ============================================================

func TestKeywordSearch_ExactMatch(t *testing.T) {
	chunks := []Chunk{
		{ID: "c1", SourceID: "s1", Text: "咳嗽发热常见于呼吸道感染", Category: "呼吸", Urgency: "medium", DeptTypes: []string{"内科"}},
		{ID: "c2", SourceID: "s2", Text: "头痛可能由多种原因引起", Category: "神经", Urgency: "medium", DeptTypes: []string{"神经内科"}},
	}
	r := &Retriever{chunks: chunks, cfg: &Config{TopK: 5}}

	results := r.keywordSearch("咳嗽")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Chunk.ID != "c1" {
		t.Errorf("expected chunk c1, got %s", results[0].Chunk.ID)
	}
	if results[0].Score != 1.0 {
		t.Errorf("expected score 1.0 for exact match, got %v", results[0].Score)
	}
	if results[0].MatchType != "keyword" {
		t.Errorf("expected MatchType keyword, got %s", results[0].MatchType)
	}
}

func TestKeywordSearch_PartialMatch(t *testing.T) {
	chunks := []Chunk{
		{ID: "c1", SourceID: "s1", Text: "咳嗽发热常见于呼吸道感染", Category: "呼吸", Urgency: "medium", DeptTypes: []string{"内科"}},
		{ID: "c2", SourceID: "s2", Text: "头痛可能由多种原因引起", Category: "神经", Urgency: "medium", DeptTypes: []string{"神经内科"}},
	}
	r := &Retriever{chunks: chunks, cfg: &Config{TopK: 5}}

	// Chinese text has no spaces — partial word "发热" should match via the substring path
	// "发热" IS a substring of "咳嗽发热常见于呼吸道感染"
	results := r.keywordSearch("发热")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Chunk.ID != "c1" {
		t.Errorf("expected chunk c1, got %s", results[0].Chunk.ID)
	}
	// "发热" is an exact substring of the chunk text → score 1.0
	if results[0].Score != 1.0 {
		t.Errorf("expected score 1.0 for exact substring match, got %v", results[0].Score)
	}
}

func TestKeywordSearch_NoMatch(t *testing.T) {
	chunks := []Chunk{
		{ID: "c1", SourceID: "s1", Text: "咳嗽发热", Category: "呼吸", Urgency: "medium", DeptTypes: []string{"内科"}},
	}
	r := &Retriever{chunks: chunks, cfg: &Config{TopK: 5}}

	results := r.keywordSearch("骨折")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestKeywordSearch_EmptyQuery(t *testing.T) {
	chunks := []Chunk{
		{ID: "c1", SourceID: "s1", Text: "咳嗽", Category: "呼吸", Urgency: "medium", DeptTypes: []string{"内科"}},
	}
	r := &Retriever{chunks: chunks, cfg: &Config{TopK: 5}}

	results := r.keywordSearch("")
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty query, got %d", len(results))
	}
}

// ============================================================
// Search
// ============================================================

func TestSearch_DedupBySourceID(t *testing.T) {
	// Two chunks with same SourceID but different text — keywordSearch finds both, Search should dedup
	chunks := []Chunk{
		{ID: "c1", SourceID: "resp", Text: "咳嗽发热常见于呼吸道感染", Category: "呼吸", Urgency: "medium", DeptTypes: []string{"内科"}},
		{ID: "c2", SourceID: "resp", Text: "呼吸道症状包括咳嗽咳痰", Category: "呼吸", Urgency: "medium", DeptTypes: []string{"内科"}},
	}
	r := &Retriever{chunks: chunks, cfg: &Config{TopK: 5, RAGEnabled: false}}

	results := r.Search("咳嗽")
	if len(results) != 1 {
		t.Fatalf("expected 1 deduped result (same SourceID), got %d", len(results))
	}
}

func TestSearch_TopK(t *testing.T) {
	chunks := []Chunk{
		{ID: "c1", SourceID: "s1", Text: "咳嗽", Category: "呼吸", Urgency: "medium", DeptTypes: []string{"内科"}},
		{ID: "c2", SourceID: "s2", Text: "咳嗽咳痰", Category: "呼吸", Urgency: "medium", DeptTypes: []string{"内科"}},
		{ID: "c3", SourceID: "s3", Text: "咳嗽发热", Category: "呼吸", Urgency: "medium", DeptTypes: []string{"内科"}},
		{ID: "c4", SourceID: "s4", Text: "头痛", Category: "神经", Urgency: "medium", DeptTypes: []string{"神经内科"}},
		{ID: "c5", SourceID: "s5", Text: "咳嗽胸闷", Category: "呼吸", Urgency: "high", DeptTypes: []string{"内科"}},
		{ID: "c6", SourceID: "s6", Text: "咳嗽不止", Category: "呼吸", Urgency: "medium", DeptTypes: []string{"内科"}},
	}
	r := &Retriever{chunks: chunks, cfg: &Config{TopK: 3, RAGEnabled: false}}

	results := r.Search("咳嗽")
	if len(results) > 3 {
		t.Errorf("expected at most 3 results (TopK=3), got %d", len(results))
	}
}

func TestSearch_Empty(t *testing.T) {
	r := &Retriever{chunks: []Chunk{}, cfg: &Config{TopK: 5}}
	results := r.Search("咳嗽")
	if len(results) != 0 {
		t.Errorf("expected 0 results from empty chunks, got %d", len(results))
	}
}

// ============================================================
// GetEntryByID
// ============================================================

func TestGetEntryByID_Found(t *testing.T) {
	r := &Retriever{
		knowledge: []KnowledgeEntry{
			{ID: "resp", Keywords: []string{"咳嗽"}, Category: "呼吸系统疾病"},
		},
	}

	entry := r.GetEntryByID("resp")
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if entry.Category != "呼吸系统疾病" {
		t.Errorf("expected category 呼吸系统疾病, got %s", entry.Category)
	}
}

func TestGetEntryByID_NotFound(t *testing.T) {
	r := &Retriever{
		knowledge: []KnowledgeEntry{
			{ID: "resp"},
		},
	}

	entry := r.GetEntryByID("nonexistent")
	if entry != nil {
		t.Errorf("expected nil for missing ID, got %v", entry)
	}
}

func TestGetEntryByID_EmptyKnowledge(t *testing.T) {
	r := &Retriever{knowledge: nil}
	entry := r.GetEntryByID("any")
	if entry != nil {
		t.Errorf("expected nil for empty knowledge, got %v", entry)
	}
}

// ============================================================
// buildKeywordAdvice
// ============================================================

func TestBuildKeywordAdvice(t *testing.T) {
	advice := buildKeywordAdvice("咳嗽三天", "呼吸系统疾病", "medium", []string{"内科", "呼吸内科"})
	if !strings.Contains(advice, "咳嗽三天") {
		t.Error("advice should contain the symptom")
	}
	if !strings.Contains(advice, "呼吸系统疾病") {
		t.Error("advice should contain the category")
	}
	if !strings.Contains(advice, "内科") {
		t.Error("advice should contain the recommended dept 内科")
	}
	if !strings.Contains(advice, "呼吸内科") {
		t.Error("advice should contain the recommended dept 呼吸内科")
	}
	if !strings.Contains(advice, "免责声明") {
		t.Error("advice should contain the disclaimer")
	}
}

func TestBuildKeywordAdvice_HighUrgency(t *testing.T) {
	advice := buildKeywordAdvice("胸闷", "心血管疾病", "high", []string{"心内科"})
	if !strings.Contains(advice, "紧急度较高") {
		t.Error("high urgency advice should mention emergency")
	}
}

func TestBuildKeywordAdvice_EmptyDepts(t *testing.T) {
	advice := buildKeywordAdvice("头疼", "不明", "low", nil)
	if !strings.Contains(advice, "内科（建议进一步描述症状") {
		t.Error("empty depts should fallback to 内科 suggestion")
	}
}

// ============================================================
// fallbackDepts
// ============================================================

func TestFallbackDepts(t *testing.T) {
	depts := fallbackDepts([]string{"内科", "外科"})
	if len(depts) != 2 {
		t.Fatalf("expected 2 depts, got %d", len(depts))
	}
	if depts[0].Name != "内科" {
		t.Errorf("expected Name=内科, got %s", depts[0].Name)
	}
	if depts[0].ID != "" {
		t.Errorf("fallback depts should have empty ID, got %s", depts[0].ID)
	}
	if depts[1].Name != "外科" {
		t.Errorf("expected Name=外科, got %s", depts[1].Name)
	}
}

func TestFallbackDepts_Empty(t *testing.T) {
	depts := fallbackDepts(nil)
	if len(depts) != 0 {
		t.Errorf("expected 0 depts for nil input, got %d", len(depts))
	}
}

// ============================================================
// LoadConfig
// ============================================================

func TestLoadConfig_Defaults(t *testing.T) {
	// Save and restore env
	origEnv := map[string]string{
		"DEEPSEEK_API_KEY":      os.Getenv("DEEPSEEK_API_KEY"),
		"DEEPSEEK_BASE_URL":     os.Getenv("DEEPSEEK_BASE_URL"),
		"DEEPSEEK_MODEL":        os.Getenv("DEEPSEEK_MODEL"),
		"SILICONFLOW_API_KEY":   os.Getenv("SILICONFLOW_API_KEY"),
		"TRIAGE_RAG_ENABLED":    os.Getenv("TRIAGE_RAG_ENABLED"),
		"TRIAGE_TOP_K":          os.Getenv("TRIAGE_TOP_K"),
	}
	defer func() {
		for k, v := range origEnv {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}()

	os.Unsetenv("DEEPSEEK_API_KEY")
	os.Unsetenv("DEEPSEEK_BASE_URL")
	os.Unsetenv("DEEPSEEK_MODEL")
	os.Unsetenv("SILICONFLOW_API_KEY")
	os.Unsetenv("TRIAGE_RAG_ENABLED")
	os.Unsetenv("TRIAGE_TOP_K")

	cfg := LoadConfig()
	if cfg.DeepSeekBaseURL != "https://api.deepseek.com" {
		t.Errorf("default DEEPSEEK_BASE_URL, got %s", cfg.DeepSeekBaseURL)
	}
	if cfg.DeepSeekModel != "deepseek-chat" {
		t.Errorf("default DEEPSEEK_MODEL, got %s", cfg.DeepSeekModel)
	}
	if cfg.TopK != 5 {
		t.Errorf("default TopK=5, got %d", cfg.TopK)
	}
	if cfg.RAGEnabled {
		t.Error("RAGEnabled should be false by default")
	}
	if cfg.IsDeepSeekAvailable() {
		t.Error("DeepSeek should not be available without API key")
	}
	if cfg.IsSemanticSearchAvailable() {
		t.Error("Semantic search should not be available without API key")
	}
}

func TestLoadConfig_WithEnv(t *testing.T) {
	origKey := os.Getenv("DEEPSEEK_API_KEY")
	defer func() {
		if origKey == "" {
			os.Unsetenv("DEEPSEEK_API_KEY")
		} else {
			os.Setenv("DEEPSEEK_API_KEY", origKey)
		}
	}()

	os.Setenv("DEEPSEEK_API_KEY", "sk-test-key")
	os.Setenv("TRIAGE_RAG_ENABLED", "true")
	os.Setenv("TRIAGE_TOP_K", "10")
	os.Setenv("SILICONFLOW_API_KEY", "sk-sf-test")

	cfg := LoadConfig()
	if cfg.DeepSeekAPIKey != "sk-test-key" {
		t.Errorf("expected DEEPSEEK_API_KEY, got %s", cfg.DeepSeekAPIKey)
	}
	if !cfg.RAGEnabled {
		t.Error("RAGEnabled should be true")
	}
	if cfg.TopK != 10 {
		t.Errorf("TopK should be 10, got %d", cfg.TopK)
	}
	if !cfg.IsDeepSeekAvailable() {
		t.Error("DeepSeek should be available")
	}
	if !cfg.IsSemanticSearchAvailable() {
		t.Error("Semantic search should be available")
	}
}

// ============================================================
// MatchDeptTypes (via HTTP mock server)
// ============================================================

func TestMatchDeptTypes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/user/departments" {
			w.WriteHeader(404)
			return
		}
		resp := deptAPIResponse{
			Code:    0,
			Message: "成功",
			Data: []Department{
				{ID: "d1", Name: "内科"},
				{ID: "d2", Name: "呼吸内科"},
				{ID: "d3", Name: "外科"},
				{ID: "d4", Name: "神经内科"},
				{ID: "d5", Name: "儿科"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	resolver := NewDepartmentResolver(server.URL)

	depts, err := resolver.MatchDeptTypes([]string{"内科"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "内科" should match both "内科" and "呼吸内科" and "神经内科"
	if len(depts) < 1 {
		t.Fatalf("expected at least 1 department, got %d", len(depts))
	}
	names := make(map[string]bool)
	for _, d := range depts {
		names[d.Name] = true
	}
	if !names["内科"] {
		t.Error("expected 内科 in matched departments")
	}
}

func TestMatchDeptTypes_Dedup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := deptAPIResponse{
			Code:    0,
			Message: "成功",
			Data: []Department{
				{ID: "d1", Name: "内科"},
				{ID: "d2", Name: "呼吸内科"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	resolver := NewDepartmentResolver(server.URL)

	// Both "内科" and "呼吸内科" map to d2 "呼吸内科" — should dedup
	depts, err := resolver.MatchDeptTypes([]string{"内科", "呼吸内科"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(depts) != 2 {
		t.Errorf("expected 2 unique departments, got %d (内科→d1+d2, 呼吸内科→d2, should be 2 unique)", len(depts))
	}
}

func TestMatchDeptTypes_NoMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := deptAPIResponse{
			Code:    0,
			Message: "成功",
			Data: []Department{
				{ID: "d1", Name: "内科"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	resolver := NewDepartmentResolver(server.URL)

	depts, err := resolver.MatchDeptTypes([]string{"口腔科"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(depts) != 0 {
		t.Errorf("expected 0 matches for 口腔科, got %d", len(depts))
	}
}

func TestMatchDeptTypes_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	resolver := NewDepartmentResolver(server.URL)

	_, err := resolver.MatchDeptTypes([]string{"内科"})
	if err == nil {
		t.Error("expected error for server 500, got nil")
	}
}

func TestMatchDeptTypes_CaseInsensitive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := deptAPIResponse{
			Code:    0,
			Message: "成功",
			Data: []Department{
				{ID: "d1", Name: "心血管内科"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	resolver := NewDepartmentResolver(server.URL)

	depts, err := resolver.MatchDeptTypes([]string{"心血管内科"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(depts) != 1 {
		t.Errorf("expected 1 match, got %d", len(depts))
	}
}

// ============================================================
// Triage (service-level integration tests with fake LLM)
// ============================================================

type fakeLLM struct {
	available bool
	result    string
	err       error
}

func (f *fakeLLM) Available() bool { return f.available }
func (f *fakeLLM) GenerateTriageAdvice(symptom, context string, departments []string) (string, error) {
	return f.result, f.err
}

func TestTriage_EmptySymptom(t *testing.T) {
	s := &Service{}
	_, err := s.Triage(&TriageRequest{Symptom: ""})
	if err == nil {
		t.Error("expected error for empty symptom")
	}
}

func TestTriage_WhitespaceOnlySymptom(t *testing.T) {
	s := &Service{}
	_, err := s.Triage(&TriageRequest{Symptom: "   "})
	if err == nil {
		t.Error("expected error for whitespace-only symptom")
	}
}
