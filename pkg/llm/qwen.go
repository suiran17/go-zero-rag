package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const chatURL = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

type QwenChat struct {
	apiKey string
	model  string
	hc     *http.Client
}

func NewQwenChat(apiKey, model string) *QwenChat {
	if model == "" {
		model = "qwen-plus"
	}
	return &QwenChat{apiKey: apiKey, model: model, hc: &http.Client{}}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *QwenChat) Complete(messages []Message) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"model":    c.model,
		"messages": messages,
	})

	req, _ := http.NewRequest(http.MethodPost, chatURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("decode chat response: %w", err)
	}
	if result.Error.Message != "" {
		return "", fmt.Errorf("qwen chat error: %s", result.Error.Message)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty choices in response")
	}
	return result.Choices[0].Message.Content, nil
}
