package embed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	qwenEmbedURL = "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
	batchSize    = 25 // Qwen API max texts per request
)

type QwenClient struct {
	apiKey string
	hc     *http.Client
}

func NewQwenClient(apiKey string) *QwenClient {
	return &QwenClient{apiKey: apiKey, hc: &http.Client{}}
}

// Embed returns one embedding vector per input text, in the same order.
func (c *QwenClient) Embed(texts []string) ([][]float32, error) {
	result := make([][]float32, len(texts))
	for start := 0; start < len(texts); start += batchSize {
		end := start + batchSize
		if end > len(texts) {
			end = len(texts)
		}
		vecs, err := c.embedBatch(texts[start:end])
		if err != nil {
			return nil, fmt.Errorf("embed batch [%d,%d): %w", start, end, err)
		}
		copy(result[start:end], vecs)
	}
	return result, nil
}

type embedRequest struct {
	Model  string      `json:"model"`
	Input  embedInput  `json:"input"`
	Params embedParams `json:"parameters"`
}

type embedInput struct {
	Texts []string `json:"texts"`
}

type embedParams struct {
	Dimension int `json:"dimension"`
}

type embedResponse struct {
	Output struct {
		Embeddings []struct {
			TextIndex int       `json:"text_index"`
			Embedding []float32 `json:"embedding"`
		} `json:"embeddings"`
	} `json:"output"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (c *QwenClient) embedBatch(texts []string) ([][]float32, error) {
	body, _ := json.Marshal(embedRequest{
		Model:  "text-embedding-v3",
		Input:  embedInput{Texts: texts},
		Params: embedParams{Dimension: 1536},
	})

	req, _ := http.NewRequest(http.MethodPost, qwenEmbedURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var er embedResponse
	if err := json.Unmarshal(raw, &er); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if er.Code != "" {
		return nil, fmt.Errorf("qwen error %s: %s", er.Code, er.Message)
	}

	vecs := make([][]float32, len(texts))
	for _, e := range er.Output.Embeddings {
		vecs[e.TextIndex] = e.Embedding
	}
	return vecs, nil
}
