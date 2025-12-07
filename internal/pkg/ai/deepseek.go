package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type deepSeekClient struct {
	apiKey  string
	baseURL string
}

func NewDeepSeekClient(apiKey, baseURL string) Client {
	if baseURL == "" {
		baseURL = "https://api.deepseek.com" // 默认 BaseURL
	}
	// 移除末尾的斜杠，防止拼接出错
	baseURL = strings.TrimRight(baseURL, "/")
	return &deepSeekClient{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// DeepSeek 请求结构
type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeek 响应结构
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *deepSeekClient) GradeQuestion(q, std, stu string, fullScore int) (*GradingResult, error) {
	// 1. 构造 Prompt
	prompt := fmt.Sprintf(`你是一个智能助教。请根据以下信息对学生的回答进行评分。

题目内容: %s
标准答案: %s
学生回答: %s
题目分值: %d

请严格按照以下 JSON 格式输出结果（不要包含 Markdown 标记或其他多余文本）：
{
    "score": <int, 0到满分>,
    "is_correct": <bool>,
    "comment": "<string, 简短评语>"
}`, q, std, stu, fullScore)

	// 2. 构造请求
	reqBody := chatRequest{
		Model: "deepseek-chat",
		Messages: []message{
			{Role: "system", Content: "You are a helpful assistant that outputs only JSON."},
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("deepseek api error: %d %s", resp.StatusCode, string(bodyBytes))
	}

	// 4. 解析响应
	var chatResp chatResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return nil, fmt.Errorf("parse response failed: %v", err)
	}
	if chatResp.Error != nil {
		return nil, errors.New(chatResp.Error.Message)
	}
	if len(chatResp.Choices) == 0 {
		return nil, errors.New("empty choices from deepseek")
	}

	content := chatResp.Choices[0].Message.Content
	// 清理可能的 Markdown 标记 (```json ... ```)
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")

	var result GradingResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// 如果解析失败，可能是 AI 返回了非 JSON 格式，降级处理
		return &GradingResult{
			Score:     0,
			IsCorrect: false,
			Comment:   "AI 评分格式解析失败，请人工复核。原始回复: " + content,
		}, nil
	}

	return &result, nil
}
