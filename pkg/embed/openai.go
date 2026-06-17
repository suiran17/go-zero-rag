package embed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const batchSize = 25 // max texts per request

// Client calls an OpenAI-compatible /v1/embeddings endpoint. It works against
// local servers (LM Studio, Ollama, vLLM) as well as cloud providers; for
// keyless local servers the apiKey is simply left empty.
type Client struct {
	baseURL string // e.g. "http://localhost:1234/v1"
	model   string
	apiKey  string // optional; empty for keyless local servers
	hc      *http.Client
}

func New(baseURL, model, apiKey string) *Client {
	return &Client{baseURL: baseURL, model: model, apiKey: apiKey, hc: &http.Client{}}
}

// Embed returns one embedding vector per input text, in the same order.
func (c *Client) Embed(texts []string) ([][]float32, error) {
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

type embedResponse struct {
	Data []struct {
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Client) embedBatch(texts []string) ([][]float32, error) {
	body, _ := json.Marshal(map[string]any{
		"model": c.model,
		"input": texts,
	})

	req, _ := http.NewRequest(http.MethodPost, c.baseURL+"/embeddings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

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
	if er.Error.Message != "" {
		return nil, fmt.Errorf("embed error: %s", er.Error.Message)
	}
	if len(er.Data) != len(texts) {
		return nil, fmt.Errorf("expected %d embeddings, got %d", len(texts), len(er.Data))
	}

	vecs := make([][]float32, len(texts))
	for _, d := range er.Data {
		vecs[d.Index] = d.Embedding
	}
	return vecs, nil
}
