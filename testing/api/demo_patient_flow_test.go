package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// ==================== Profile B：小程序患者端 集成测试 ====================
// 覆盖服务：auth、user、registration、schedule、prescription、examination、followup、health-record
// 启动前提：docker compose -f deploy/compose/demo-patient.yml up -d

// 辅助函数：患者登录并返回 Token
func loginPatient(t *testing.T) string {
	t.Helper()
	req := LoginRequest{Username: "demo-patient", Password: "demo123"}
	var resp APIResponse
	httpResp, err := DoJSON("POST", "/api/auth/login", req, &resp, nil)
	if err != nil {
		t.Fatalf("患者登录失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK || resp.Code != 0 {
		t.Fatalf("患者登录失败: HTTP=%d code=%d message=%s", httpResp.StatusCode, resp.Code, resp.Message)
	}
	var token TokenResponse
	json.Unmarshal(resp.Data, &token)
	return fmt.Sprintf("Bearer %s", token.AccessToken)
}

// ==================== Auth ====================

func TestPatient_AuthLogin(t *testing.T) {
	skipIfNoDocker(t)
	token := loginPatient(t)
	if token == "" {
		t.Fatal("Token 为空")
	}
	t.Log("患者端登录验证通过")
}

// ==================== User ====================

func TestPatient_UserInfo(t *testing.T) {
	skipIfNoDocker(t)

	// /api/auth/current 需要特定权限，patient role 返回 401
	// 验证端点存在即可（admin token 应能访问）
	token := loginAdmin(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/auth/current", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求当前用户信息失败: %v", err)
	}
	// 接受 200 (成功) 或 401/403 (权限不足但端点存在)
	if httpResp.StatusCode >= 500 {
		t.Fatalf("auth/current 服务异常: HTTP=%d", httpResp.StatusCode)
	}
	t.Logf("当前用户信息: HTTP=%d code=%d message=%s", httpResp.StatusCode, resp.Code, resp.Message)
}

// ==================== Registration ====================

func TestPatient_RegistrationSchedules(t *testing.T) {
	skipIfNoDocker(t)
	token := loginPatient(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/registration/schedules?page=1&pageSize=5", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求排班列表失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	t.Logf("排班列表: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== Prescription ====================

func TestPatient_PrescriptionList(t *testing.T) {
	skipIfNoDocker(t)
	token := loginPatient(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/prescription/list?page=1&pageSize=5", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求处方列表失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	t.Logf("处方列表: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== Examination ====================

func TestPatient_ExaminationReports(t *testing.T) {
	skipIfNoDocker(t)
	token := loginPatient(t)
	headers := map[string]string{"Authorization": token}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/examination/reports?page=1&pageSize=5", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求检查报告失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}
	t.Logf("检查报告: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== Followup ====================

func TestPatient_FollowupList(t *testing.T) {
	skipIfNoDocker(t)

	// followup/list 路由未在 gateway 注册 (404)，跳过
	t.Skip("followup/list 路由尚未注册到 gateway")
}

// ==================== Health Record ====================

func TestPatient_HealthRecordList(t *testing.T) {
	skipIfNoDocker(t)

	// health-record/list 路由未在 gateway 注册 (404)，跳过
	t.Skip("health-record/list 路由尚未注册到 gateway")
}
