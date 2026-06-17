package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client calls an OpenAI-compatible /v1/chat/completions endpoint. The same
// client works for DeepSeek, Qwen DashScope compatible-mode, and local servers
// — only baseURL, model and apiKey differ.
type Client struct {
	baseURL string // e.g. "https://api.deepseek.com/v1"
	model   string
	apiKey  string
	hc      *http.Client
}

func New(baseURL, model, apiKey string) *Client {
	if model == "" {
		model = "deepseek-chat"
	}
	return &Client{baseURL: baseURL, model: model, apiKey: apiKey, hc: &http.Client{}}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *Client) Complete(messages []Message) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"model":    c.model,
		"messages": messages,
	})

	req, _ := http.NewRequest(http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

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
		return "", fmt.Errorf("chat error: %s", result.Error.Message)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty choices in response")
	}
	return result.Choices[0].Message.Content, nil
}
