package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	QdrantAddr   string
	VectorSize   int    `json:",default=1024"`
	EmbedBaseURL string // OpenAI-compatible base, e.g. http://localhost:1234/v1
	EmbedModel   string
	EmbedApiKey  string `json:",optional"` // empty for keyless local servers
}
