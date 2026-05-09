// Package api HIS-Go API 集成测试客户端
// 用法：当 Docker 服务启动后，运行 go test -v ./...
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	baseURL   string
	client    *http.Client
)

func init() {
	baseURL = os.Getenv("HIS_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	client = &http.Client{Timeout: 10 * time.Second}
}

// HTTPClient 返回配置好的 HTTP 客户端
func HTTPClient() *http.Client {
	return client
}

// BaseURL 返回 API 基地址
func BaseURL() string {
	return baseURL
}

// APIPath 构建完整 API 路径
func APIPath(path string) string {
	return baseURL + path
}

// DoJSON 发送 JSON 请求并解析响应
func DoJSON(method, path string, body, result interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, APIPath(path), reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Request-ID", fmt.Sprintf("test-%d", time.Now().UnixNano()))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return resp, fmt.Errorf("解析响应失败: %w", err)
		}
	}

	return resp, nil
}

// APIResponse 统一 API 响应结构
type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// TokenResponse 登录 token 响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
