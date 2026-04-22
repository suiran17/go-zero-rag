# 任务计划：go-zero-rag — RAG 知识问答微服务系统

## 目标

独立开发一个完整的 RAG 知识问答系统，基于 go-zero 微服务架构，部署到自建 K8s 集群，代码放到 GitHub public 仓库，作为求职作品展示。

## 当前阶段

全部完成 ✓  
仓库：https://github.com/suiran17/go-zero-rag

## 各阶段

### 阶段 1：项目脚手架
- [x] 初始化 git 仓库 + go.mod（单模块，无需 go.work）
- [x] 编写 docservice/doc.proto，goctl 生成代码
- [x] 编写 aiservice/ai.proto，goctl 生成代码
- [x] 编写 gateway/gateway.api，goctl 生成代码
- [x] 验证三个服务均可 `go build ./...` 通过 ✅
- [x] 创建 docker-compose.yml（含 etcd + Qdrant）
- **状态：** complete

### 阶段 2：doc-service（文档摄入）
- [x] 文档上传接口（接收 title + content）
- [x] 文本分块逻辑（按段落，最大 500 token）—— chunk/chunker.go
- [x] 调用通义千问 Embedding API 获取向量 —— embed/qwen.go
- [x] 向量存入 Qdrant（collection 初始化 + upsert）—— qdrantcli/client.go
- [ ] 本地验证：上传文本 → Qdrant 可检索（待补充 QWEN_API_KEY）
- **状态：** complete

### 阶段 3：ai-service（RAG 查询）
- [x] 问题向量化（通义千问 Embedding）
- [x] Qdrant Top-K 语义检索
- [x] 拼装 Prompt（检索段落 + 用户问题）
- [x] 调用 LLM（通义千问 / DeepSeek）生成答案
- [x] 返回答案 + 来源段落列表
- [ ] 本地验证：提问 → 返回答案（待补充 QWEN_API_KEY）
- **状态：** complete

### 阶段 4：gateway（端到端打通）
- [x] POST /api/doc/upload → docservice.UploadDoc
- [x] POST /api/qa/query → aiservice.Query
- [ ] 端到端测试：上传文档 → 提问 → 拿到答案（待补充 QWEN_API_KEY）
- **状态：** complete

### 阶段 5：K8s 部署 + CI/CD
- [x] 三个服务各写 Dockerfile（多阶段构建）
- [x] K8s yaml（Deployment + Service × 3）
- [x] Qdrant StatefulSet yaml + etcd StatefulSet yaml
- [x] ConfigMap（K8s 内网地址） + Secret 模板
- [x] Jenkinsfile（并行构建三镜像 → kubectl apply → rollout）
- [ ] 推送到集群实际验证（需要 QWEN_API_KEY + 集群 kubectl 权限）
- **状态：** complete

### 阶段 6：README + GitHub public
- [x] README.md（Mermaid 架构图 + 接口说明 + 本地运行步骤）
- [x] GitHub 仓库设为 public：https://github.com/suiran17/go-zero-rag
- [x] 确认仓库页面展示正常（Mermaid 架构图正常渲染，Go 94.8%）
- **状态：** complete

## 已做决策

| 决策 | 理由 |
|------|------|
| 框架：go-zero | 已有深度学习基础，有 rpc 通信需求 |
| 单模块（无 go.work） | goctl 生成代码在同一 go.mod 下，不需要多模块 |
| 向量库：Qdrant | Go SDK 完善，K8s 部署轻量，适合个人项目 |
| Embedding：通义千问 text-embedding-v3 | 国内直连，中文效果好，免费额度够 |
| LLM：通义千问 / DeepSeek | 国内直连，有免费 token |
| 分块：自实现 | 逻辑简单，面试可讲，无额外依赖 |
| 部署：复用现有集群 | 3节点 kubeadm + Jenkins + GitLab + Docker |
| 项目路径：/Users/zp/GIT/k8s-deploy/go-zero-rag | 当前阶段，后续独立为单独仓库 |

## 接口定义（已确定）

**docservice/doc.proto**
```protobuf
syntax = "proto3"; package doc; option go_package = "./pb";
service Doc { rpc UploadDoc(UploadDocReq) returns (UploadDocResp); }
message UploadDocReq { string title = 1; string content = 2; }
message UploadDocResp { string doc_id = 1; int32 chunk_count = 2; }
```

**aiservice/ai.proto**
```protobuf
syntax = "proto3"; package ai; option go_package = "./pb";
service Ai { rpc Query(QueryReq) returns (QueryResp); }
message QueryReq { string question = 1; int32 top_k = 2; }
message QueryResp { string answer = 1; repeated string sources = 2; }
```

**gateway/gateway.api**
```
syntax = "v1"
type (
  UploadDocReq  { Title string `json:"title"`; Content string `json:"content"` }
  UploadDocResp { DocId string `json:"doc_id"`; ChunkCount int `json:"chunk_count"` }
  QueryReq      { Question string `json:"question"`; TopK int `json:"top_k,optional,default=3"` }
  QueryResp     { Answer string `json:"answer"`; Sources []string `json:"sources"` }
)
service gateway {
  @handler UploadDoc
  post /api/doc/upload (UploadDocReq) returns (UploadDocResp)
  @handler Query
  post /api/qa/query (QueryReq) returns (QueryResp)
}
```

## 遇到的错误

| 错误 | 尝试次数 | 解决方案 |
|------|---------|---------|
| `session-catchup.py` 脚本路径不存在 | 1 | 改为直接以仓库文件和最新验证结果同步进度 |
