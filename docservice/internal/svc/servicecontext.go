package svc

import (
	"github.com/suiran17/go-zero-rag/docservice/internal/config"
	"github.com/suiran17/go-zero-rag/pkg/embed"
	"github.com/suiran17/go-zero-rag/pkg/qdrantcli"
)

type ServiceContext struct {
	Config config.Config
	Qwen   *embed.QwenClient
	Qdrant *qdrantcli.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Qwen:   embed.NewQwenClient(c.QwenApiKey),
		Qdrant: qdrantcli.NewClient(c.QdrantAddr),
	}
}
