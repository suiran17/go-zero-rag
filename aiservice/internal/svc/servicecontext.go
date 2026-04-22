package svc

import (
	"github.com/suiran17/go-zero-rag/aiservice/internal/config"
	"github.com/suiran17/go-zero-rag/pkg/embed"
	"github.com/suiran17/go-zero-rag/pkg/llm"
	"github.com/suiran17/go-zero-rag/pkg/qdrantcli"
)

type ServiceContext struct {
	Config config.Config
	Qwen   *embed.QwenClient
	Qdrant *qdrantcli.Client
	LLM    *llm.QwenChat
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Qwen:   embed.NewQwenClient(c.QwenApiKey),
		Qdrant: qdrantcli.NewClient(c.QdrantAddr),
		LLM:    llm.NewQwenChat(c.QwenApiKey, c.LLMModel),
	}
}
