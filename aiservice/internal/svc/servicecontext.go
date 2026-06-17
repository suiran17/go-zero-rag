package svc

import (
	"github.com/suiran17/go-zero-rag/aiservice/internal/config"
	"github.com/suiran17/go-zero-rag/pkg/embed"
	"github.com/suiran17/go-zero-rag/pkg/llm"
	"github.com/suiran17/go-zero-rag/pkg/qdrantcli"
)

type ServiceContext struct {
	Config   config.Config
	Embedder *embed.Client
	Qdrant   *qdrantcli.Client
	LLM      *llm.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Embedder: embed.New(c.EmbedBaseURL, c.EmbedModel, c.EmbedApiKey),
		Qdrant:   qdrantcli.NewClient(c.QdrantAddr, c.VectorSize),
		LLM:      llm.New(c.LLMBaseURL, c.LLMModel, c.LLMApiKey),
	}
}
