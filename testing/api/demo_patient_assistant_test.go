package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// ==================== 就诊助手 (Outpatient Assistant) 冒烟测试 ====================
// 覆盖：POST /api/outpatient/assistant/chat（通过 gateway 代理 + JWT）

func loginAny(t *testing.T, username string) string {
	t.Helper()
	req := LoginRequest{Username: username, Password: "demo123"}
	var resp APIResponse
	httpResp, err := DoJSON("POST", "/api/auth/login", req, &resp, nil)
	if err != nil {
		t.Fatalf("%s 登录失败: %v", username, err)
	}
	if httpResp.StatusCode != http.StatusOK || resp.Code != 0 {
		t.Fatalf("%s 登录失败: HTTP=%d code=%d", username, httpResp.StatusCode, resp.Code)
	}
	var token TokenResponse
	json.Unmarshal(resp.Data, &token)
	return "Bearer " + token.AccessToken
}

type TriageDept struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TriageData struct {
	Symptom      string       `json:"symptom"`
	Advice       string       `json:"advice"`
	Depts        []TriageDept `json:"depts"`
	KnowledgeRef string       `json:"knowledgeRef"`
	Urgency      string       `json:"urgency"`
	Mode         string       `json:"mode"`
	Disclaimer   string       `json:"disclaimer"`
}

func TestAssistant_Chat_PatientAuth(t *testing.T) {
	skipIfNoDocker(t)

	token := loginAny(t, "demo-patient")
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("POST", "/api/outpatient/assistant/chat",
		map[string]string{"symptom": "咳嗽三天发烧"},
		&resp, headers)
	if err != nil {
		t.Fatalf("请求就诊助手失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	if resp.Code != 0 {
		t.Fatalf("期望业务码=0，实际=%d message=%s", resp.Code, resp.Message)
	}

	var data TriageData
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	// 基本字段检查
	if data.Symptom == "" {
		t.Error("symptom 不应为空")
	}
	if data.Advice == "" {
		t.Error("advice 不应为空（关键词或 LLM 模式均应有建议）")
	}
	if len(data.Depts) == 0 {
		t.Error("depts 不应为空（至少应有兜底科室）")
	}
	if data.KnowledgeRef == "" {
		t.Error("knowledgeRef 不应为空")
	}
	if data.Disclaimer == "" {
		t.Error("disclaimer 不应为空")
	}
	if !strings.Contains(data.Disclaimer, "仅供参考") {
		t.Error("disclaimer 应包含「仅供参考」")
	}

	// Mode 必须为 llm 或 keyword
	if data.Mode != "llm" && data.Mode != "keyword" {
		t.Errorf("mode 应为 llm 或 keyword，实际=%s", data.Mode)
	}

	// Urgency 必须为 high/medium/low 之一
	switch data.Urgency {
	case "high", "medium", "low":
	default:
		t.Errorf("urgency 应为 high/medium/low，实际=%s", data.Urgency)
	}

	t.Logf("就诊助手响应: mode=%s urgency=%s depts=%d advice=%.60s...",
		data.Mode, data.Urgency, len(data.Depts), data.Advice)
}

func TestAssistant_Chat_EmptySymptom(t *testing.T) {
	skipIfNoDocker(t)

	token := loginAny(t, "demo-patient")
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("POST", "/api/outpatient/assistant/chat",
		map[string]string{"symptom": ""},
		&resp, headers)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	// 空症状应返回业务错误（非 500）
	if resp.Code == 0 {
		t.Error("空症状应返回错误码")
	}
	if httpResp.StatusCode >= 500 {
		t.Errorf("空症状不应返回 5xx，实际=%d", httpResp.StatusCode)
	}
	t.Logf("空症状响应: HTTP=%d code=%d message=%s", httpResp.StatusCode, resp.Code, resp.Message)
}

func TestAssistant_Chat_Unauthorized(t *testing.T) {
	skipIfNoDocker(t)

	var resp APIResponse
	httpResp, err := DoJSON("POST", "/api/outpatient/assistant/chat",
		map[string]string{"symptom": "咳嗽"},
		&resp, nil)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	// 无 Token 应被网关拦截
	if resp.Code == 0 {
		t.Error("无 Token 应被拒绝")
	}
	if httpResp.StatusCode != http.StatusUnauthorized && resp.Code != 10002 {
		t.Logf("无 Token 响应: HTTP=%d code=%d message=%s", httpResp.StatusCode, resp.Code, resp.Message)
	}
}
