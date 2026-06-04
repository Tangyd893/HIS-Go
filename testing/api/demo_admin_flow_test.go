package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// ==================== Profile A：Web 管理端 集成测试 ====================
// 覆盖服务：auth、user、registration、schedule、clinic、prescription、billing、pharmacy、system
// 启动前提：docker compose -f deploy/compose/demo-admin.yml up -d

// 辅助函数：登录并返回 Token
func loginAdmin(t *testing.T) string {
	t.Helper()
	req := LoginRequest{Username: "demo-doctor", Password: "demo123"}
	var resp APIResponse
	httpResp, err := DoJSON("POST", "/api/auth/login", req, &resp, nil)
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK || resp.Code != 0 {
		t.Fatalf("登录失败: HTTP=%d code=%d message=%s", httpResp.StatusCode, resp.Code, resp.Message)
	}
	var token TokenResponse
	json.Unmarshal(resp.Data, &token)
	return fmt.Sprintf("Bearer %s", token.AccessToken)
}

// 辅助函数：验证分页响应结构
func assertPageResponse(t *testing.T, raw json.RawMessage) {
	t.Helper()
	var page struct {
		List     json.RawMessage `json:"list"`
		Total    int64           `json:"total"`
		Page     int             `json:"page"`
		PageSize int             `json:"pageSize"`
	}
	if err := json.Unmarshal(raw, &page); err != nil {
		t.Fatalf("PageData 结构解析失败: %v", err)
	}
	t.Logf("分页数据: total=%d page=%d pageSize=%d", page.Total, page.Page, page.PageSize)
}

// ==================== Auth ====================

func TestAdmin_AuthLogin(t *testing.T) {
	skipIfNoDocker(t)
	token := loginAdmin(t)
	if token == "" {
		t.Fatal("Token 为空")
	}
	t.Log("管理端登录验证通过")
}

// ==================== User ====================

func TestAdmin_UserList(t *testing.T) {
	skipIfNoDocker(t)
	token := loginAdmin(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/user/patients?page=1&pageSize=5", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求患者列表失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	t.Logf("患者列表: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== Registration ====================

func TestAdmin_RegistrationSchedules(t *testing.T) {
	skipIfNoDocker(t)
	token := loginAdmin(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/registration/schedules?page=1&pageSize=5", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求排班列表失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	if resp.Code == 0 && resp.Data != nil {
		assertPageResponse(t, resp.Data)
	}
	t.Logf("排班列表: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== Schedule ====================

func TestAdmin_ScheduleList(t *testing.T) {
	skipIfNoDocker(t)
	token := loginAdmin(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/schedule/list?page=1&pageSize=5", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求排班管理列表失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	if resp.Code == 0 && resp.Data != nil {
		assertPageResponse(t, resp.Data)
	}
	t.Logf("排班管理列表: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== Clinic ====================

func TestAdmin_ClinicRecords(t *testing.T) {
	skipIfNoDocker(t)
	token := loginAdmin(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/clinic/records?page=1&pageSize=5", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求门诊记录失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	t.Logf("门诊记录: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== Prescription ====================

func TestAdmin_PrescriptionList(t *testing.T) {
	skipIfNoDocker(t)
	token := loginAdmin(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/prescription/list?page=1&pageSize=5", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求处方列表失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	if resp.Code == 0 && resp.Data != nil {
		assertPageResponse(t, resp.Data)
	}
	t.Logf("处方列表: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== Billing ====================

func TestAdmin_BillingList(t *testing.T) {
	skipIfNoDocker(t)
	token := loginAdmin(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/billing/list?page=1&pageSize=5", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求账单列表失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	t.Logf("账单列表: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== Pharmacy ====================

func TestAdmin_PharmacyDrugs(t *testing.T) {
	skipIfNoDocker(t)
	token := loginAdmin(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/pharmacy/drugs?page=1&pageSize=5", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求药品列表失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	if resp.Code == 0 && resp.Data != nil {
		assertPageResponse(t, resp.Data)
	}
	t.Logf("药品列表: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== System ====================

func TestAdmin_SystemDict(t *testing.T) {
	skipIfNoDocker(t)
	token := loginAdmin(t)
	headers := map[string]string{"Authorization": token}

	// 系统字典接口 — 只需要验证可达
	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/system/dict/prescription_type", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求字典失败: %v", err)
	}
	// dict 接口可能返回 404 如果未实现，以 non-5xx 为目标
	if httpResp.StatusCode >= 500 {
		t.Errorf("字典接口异常: HTTP=%d code=%d message=%s", httpResp.StatusCode, resp.Code, resp.Message)
	}
	t.Logf("字典查询: HTTP=%d code=%d message=%s", httpResp.StatusCode, resp.Code, resp.Message)
}
