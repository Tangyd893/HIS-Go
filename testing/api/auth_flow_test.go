package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
)

func skipIfNoDocker(t *testing.T) {
	t.Helper()
	if os.Getenv("HIS_INTEGRATION_TEST") != "true" {
		t.Skip("跳过集成测试，设置 HIS_INTEGRATION_TEST=true 以运行")
	}
}

// ==================== 健康检查测试 ====================

func TestGatewayHealth(t *testing.T) {
	skipIfNoDocker(t)

	var body map[string]interface{}
	resp, err := DoJSON("GET", "/health", nil, &body, nil)
	if err != nil {
		t.Fatalf("健康检查失败: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期望 HTTP 200，实际=%d", resp.StatusCode)
	}

	service, ok := body["service"]
	if !ok || service != "his-gateway" {
		t.Errorf("期望 service=his-gateway，实际=%v", service)
	}
	t.Logf("Gateway 健康检查通过: %v", body)
}

// ==================== 认证流程验收 ====================

func TestAuthLogin_Success(t *testing.T) {
	skipIfNoDocker(t)

	req := LoginRequest{
		Username: "demo-doctor",
		Password: "demo123",
	}

	var resp APIResponse
	httpResp, err := DoJSON("POST", "/api/auth/login", req, &resp, nil)
	if err != nil {
		t.Fatalf("登录请求失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("期望 HTTP 200，实际=%d", httpResp.StatusCode)
	}

	if resp.Code != 0 {
		t.Fatalf("期望业务码=0，实际=%d message=%s", resp.Code, resp.Message)
	}

	var token TokenResponse
	if err := json.Unmarshal(resp.Data, &token); err != nil {
		t.Fatalf("解析 Token 失败: %v", err)
	}

	if token.AccessToken == "" {
		t.Error("AccessToken 不能为空")
	}
	if token.RefreshToken == "" {
		t.Error("RefreshToken 不能为空")
	}
	if token.ExpiresIn <= 0 {
		t.Errorf("ExpiresIn 期望>0，实际=%d", token.ExpiresIn)
	}

	t.Logf("登录成功: access_token=%s...", token.AccessToken[:20])
}

func TestAuthLogin_WrongPassword(t *testing.T) {
	skipIfNoDocker(t)

	req := LoginRequest{
		Username: "demo-doctor",
		Password: "wrong_password",
	}

	var resp APIResponse
	httpResp, err := DoJSON("POST", "/api/auth/login", req, &resp, nil)
	if err != nil {
		t.Fatalf("登录请求失败: %v", err)
	}

	if httpResp.StatusCode == http.StatusOK && resp.Code == 0 {
		t.Error("期望错误密码登录失败")
	}
	t.Logf("错误密码登录返回: code=%d message=%s", resp.Code, resp.Message)
}

func TestAuthLogin_UserNotFound(t *testing.T) {
	skipIfNoDocker(t)

	req := LoginRequest{
		Username: "nonexistent_user_12345",
		Password: "password",
	}

	var resp APIResponse
	httpResp, err := DoJSON("POST", "/api/auth/login", req, &resp, nil)
	if err != nil {
		t.Fatalf("登录请求失败: %v", err)
	}

	if httpResp.StatusCode == http.StatusOK && resp.Code == 0 {
		t.Error("期望不存在的用户登录失败")
	}
	t.Logf("不存在用户登录返回: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== Gateway JWT 鉴权验收 ====================

func TestGatewayRejectsNoToken(t *testing.T) {
	skipIfNoDocker(t)

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/user/patients", nil, &resp, nil)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if httpResp.StatusCode != http.StatusUnauthorized && resp.Code == 0 {
		t.Errorf("期望无 Token 时被拒绝，HTTP=%d code=%d", httpResp.StatusCode, resp.Code)
	}
	t.Logf("无 Token 访问返回: code=%d message=%s", resp.Code, resp.Message)
}

func TestGatewayRejectsInvalidToken(t *testing.T) {
	skipIfNoDocker(t)

	headers := map[string]string{
		"Authorization": "Bearer invalid.token.here",
	}

	var resp APIResponse
	httpResp, err := DoJSON("GET", "/api/user/patients", nil, &resp, headers)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if httpResp.StatusCode != http.StatusUnauthorized && resp.Code == 0 {
		t.Errorf("期望无效 Token 被拒绝，HTTP=%d code=%d", httpResp.StatusCode, resp.Code)
	}
	t.Logf("无效 Token 访问返回: code=%d message=%s", resp.Code, resp.Message)
}

// ==================== 完整业务流程验收 ====================

func TestFullAuthFlow(t *testing.T) {
	skipIfNoDocker(t)

	// 1. 登录获取 Token
	req := LoginRequest{
		Username: "demo-doctor",
		Password: "demo123",
	}
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
	tokenHeader := fmt.Sprintf("Bearer %s", token.AccessToken)

	// 2. 携带 Token 访问受保护接口
	authHeaders := map[string]string{
		"Authorization": tokenHeader,
	}

	var userResp APIResponse
	_, err = DoJSON("GET", "/api/user/patients", nil, &userResp, authHeaders)
	if err != nil {
		t.Fatalf("获取患者列表失败: %v", err)
	}
	t.Logf("携带 Token 获取患者列表: code=%d message=%s", userResp.Code, userResp.Message)

	// 3. 刷新 Token
	var refreshResp APIResponse
	refreshBody := map[string]string{
		"refresh_token": token.RefreshToken,
	}
	_, err = DoJSON("POST", "/api/auth/refresh", refreshBody, &refreshResp, nil)
	if err != nil {
		t.Fatalf("刷新 Token 失败: %v", err)
	}
	t.Logf("刷新 Token: code=%d message=%s", refreshResp.Code, refreshResp.Message)
}

// ==================== 用户上下文透传验收 ====================

func TestUserContextPropagation(t *testing.T) {
	skipIfNoDocker(t)

	req := LoginRequest{
		Username: "demo-doctor",
		Password: "demo123",
	}
	var loginResp APIResponse
	httpResp, err := DoJSON("POST", "/api/auth/login", req, &loginResp, nil)
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if httpResp.StatusCode != http.StatusOK || loginResp.Code != 0 {
		t.Fatalf("登录失败: code=%d", loginResp.Code)
	}

	var token TokenResponse
	json.Unmarshal(loginResp.Data, &token)
	tokenHeader := fmt.Sprintf("Bearer %s", token.AccessToken)
	headers := map[string]string{
		"Authorization": tokenHeader,
	}

	resp2, err := DoJSON("GET", "/api/user/patients", nil, nil, headers)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	proxyBy := resp2.Header.Get("X-Proxy-By")
	if proxyBy != "his-gateway" {
		t.Errorf("期望 X-Proxy-By=his-gateway，实际=%s", proxyBy)
	}

	t.Logf("用户上下文透传验证通过: X-Proxy-By=%s", proxyBy)
}
