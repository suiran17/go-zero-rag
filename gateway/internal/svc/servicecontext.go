// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/suiran17/go-zero-rag/aiservice/ai"
	"github.com/suiran17/go-zero-rag/docservice/doc"
	"github.com/suiran17/go-zero-rag/gateway/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	DocRpc doc.Doc
	AiRpc  ai.Ai
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DocRpc: doc.NewDoc(zrpc.MustNewClient(c.DocRpc)),
		AiRpc:  ai.NewAi(zrpc.MustNewClient(c.AiRpc)),
	}
}
