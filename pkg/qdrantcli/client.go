package qdrantcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	CollectionName = "documents"
	VectorSize     = 1536
)

type Client struct {
	addr string // e.g. "http://localhost:6333"
	hc   *http.Client
}

func NewClient(addr string) *Client {
	return &Client{addr: addr, hc: &http.Client{}}
}

// EnsureCollection creates the collection if it does not already exist.
func (c *Client) EnsureCollection() error {
	url := fmt.Sprintf("%s/collections/%s", c.addr, CollectionName)
	resp, err := c.hc.Get(url)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil // already exists
	}

	body, _ := json.Marshal(map[string]any{
		"vectors": map[string]any{
			"size":     VectorSize,
			"distance": "Cosine",
		},
	})
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp2, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp2.Body)
		return fmt.Errorf("create collection status %d: %s", resp2.StatusCode, raw)
	}
	return nil
}

// Point represents one chunk stored in Qdrant.
type Point struct {
	ID      string    // UUID string
	Vector  []float32
	DocID   string
	Text    string
	ChunkIdx int
}

// SearchResult is one hit returned by Qdrant.
type SearchResult struct {
	Score    float32
	DocID    string
	Text     string
	ChunkIdx int
}

// Search performs a cosine similarity search and returns the top-k results.
func (c *Client) Search(vector []float32, topK int) ([]SearchResult, error) {
	body, _ := json.Marshal(map[string]any{
		"vector":       vector,
		"limit":        topK,
		"with_payload": true,
	})
	url := fmt.Sprintf("%s/collections/%s/points/search", c.addr, CollectionName)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var result struct {
		Result []struct {
			Score   float32 `json:"score"`
			Payload struct {
				DocID    string `json:"doc_id"`
				Text     string `json:"text"`
				ChunkIdx int    `json:"chunk_index"`
			} `json:"payload"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	hits := make([]SearchResult, len(result.Result))
	for i, r := range result.Result {
		hits[i] = SearchResult{
			Score:    r.Score,
			DocID:    r.Payload.DocID,
			Text:     r.Payload.Text,
			ChunkIdx: r.Payload.ChunkIdx,
		}
	}
	return hits, nil
}

// Upsert stores points into the collection.
func (c *Client) Upsert(points []Point) error {
	type payload struct {
		DocID    string `json:"doc_id"`
		Text     string `json:"text"`
		ChunkIdx int    `json:"chunk_index"`
	}
	type point struct {
		ID      string    `json:"id"`
		Vector  []float32 `json:"vector"`
		Payload payload   `json:"payload"`
	}
	var pts []point
	for _, p := range points {
		pts = append(pts, point{
			ID:     p.ID,
			Vector: p.Vector,
			Payload: payload{
				DocID:    p.DocID,
				Text:     p.Text,
				ChunkIdx: p.ChunkIdx,
			},
		})
	}

	body, _ := json.Marshal(map[string]any{"points": pts})
	url := fmt.Sprintf("%s/collections/%s/points", c.addr, CollectionName)
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upsert status %d: %s", resp.StatusCode, raw)
	}
	return nil
}
