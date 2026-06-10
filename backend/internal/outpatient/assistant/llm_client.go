package assistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"his-go/pkg/errors"
)

// LLMClient DeepSeek Chat 客户端
type LLMClient struct {
	cfg     *Config
	httpCli *http.Client
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest DeepSeek API 请求
type chatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// chatResponse DeepSeek API 响应
type chatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// NewLLMClient 创建 DeepSeek 客户端
func NewLLMClient(cfg *Config) *LLMClient {
	return &LLMClient{
		cfg: cfg,
		httpCli: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Available 检查客户端是否可用
func (c *LLMClient) Available() bool {
	return c.cfg.IsDeepSeekAvailable()
}

// GenerateTriageAdvice 生成分诊建议
func (c *LLMClient) GenerateTriageAdvice(symptom string, context string, departments []string) (string, error) {
	if !c.Available() {
		return "", fmt.Errorf("DeepSeek API 未配置")
	}

	deptList := strings.Join(departments, "、")
	if deptList == "" {
		deptList = "内科、外科、儿科、妇产科、急诊科"
	}

	systemPrompt := `你是一个医院就诊助手，帮助患者根据症状推荐挂号科室。
请根据患者描述的症状和以下参考信息，给出简洁的挂号建议。

要求：
1. 推荐1-2个最适合的科室
2. 简要解释推荐理由（1-2句话）
3. 如果症状较严重或紧急，提醒患者尽快就医
4. 回复末尾必须包含免责声明："⚠️ 本建议仅供参考，不能替代专业医疗诊断。如症状严重请及时就医。"

回复格式：
**推荐科室**：科室A / 科室B
**推荐理由**：简要原因
**注意事项**：（如适用）
**免责声明**：⚠️ 本建议仅供参考，不能替代专业医疗诊断。如症状严重请及时就医。`

	userPrompt := fmt.Sprintf(`患者症状描述：%s

参考信息（知识库匹配）：
%s

本院可用科室：%s

请给出挂号建议。`, symptom, context, deptList)

	req := chatRequest{
		Model: c.cfg.DeepSeekModel,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: false,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	url := strings.TrimRight(c.cfg.DeepSeekBaseURL, "/") + "/v1/chat/completions"
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", errors.WrapCreateError("请求", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.cfg.DeepSeekAPIKey)

	resp, err := c.httpCli.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("请求 DeepSeek API 失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepSeek API 返回错误 (%d): %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("DeepSeek 未返回有效回复")
	}

	return chatResp.Choices[0].Message.Content, nil
}
