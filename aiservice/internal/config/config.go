package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	QwenApiKey string
	QdrantAddr string
	LLMModel   string // default: qwen-plus
}
