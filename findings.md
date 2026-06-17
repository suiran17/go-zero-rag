# 发现与决策

## 需求

- 做一个 100% 自己写的 Go RAG 微服务项目，作为求职作品
- 部署到已有 3 节点 K8s 集群，接入现有 Jenkins CI/CD
- 放到 GitHub public 仓库，附完整 README + 架构图

## 环境信息

- Go：已安装
- goctl：已安装
- K8s：3 节点虚拟机（1 master + 2 worker），kubeadm 手动搭建
- CI/CD：Jenkins + GitLab + Docker
- 监控：Prometheus + Grafana（已有）
- 服务注册：etcd（go-zero 默认）
- 项目当前路径：/Users/zp/GIT/k8s-deploy/go-zero-rag（后续独立）

## RAG 流程

```
上传阶段：
文档文本 → 按段落分块（≤500 token）
        → 本地 bge-m3 Embedding（1024 维）
        → 向量 + 原文存入 Qdrant

查询阶段：
用户问题 → 本地 bge-m3 Embedding
         → Qdrant Top-K 检索（默认 k=3）
         → 拼装 Prompt（检索段落 + 问题）
         → DeepSeek deepseek-chat
         → 返回答案 + 来源段落
```

> 2026-06-18 更新：Embedding/LLM 由通义千问切换为「本地 bge-m3 + DeepSeek」，均走 OpenAI 兼容接口；切 Embedding 模型时 `VectorSize` 必须与模型维度一致（bge-m3 = 1024）。

## 技术决策

| 决策 | 理由 |
|------|------|
| go-zero 而非 gin | 有 rpc 通信需求，已深度学习 |
| Qdrant 而非 Milvus | Go SDK 好，K8s StatefulSet 简单，轻量 |
| 通义千问 Embedding | 国内直连，中文效果好，免费额度够用 |
| 单模块 `go.mod`，不使用 `go.work` | 当前仓库规模可控，goctl 生成代码也能直接协同 |
| 分块自实现 | 无额外依赖，逻辑清晰，面试可讲 |

## 当前仓库状态

- `README.md` 已补齐，包含架构图、接口说明、本地运行步骤和 K8s/Jenkins 部署入口
- 本地重新执行 `go build ./...` 通过，可确认当前代码至少处于可编译状态
- `docs/superpowers/specs/` 与 `docs/superpowers/plans/` 目前只有 `.gitkeep`
- 仓库尚无首个 commit，也还没完成 GitHub public 发布
- 本地与集群的端到端验证仍受 `QWEN_API_KEY` 和实际 `kubectl` 权限限制

## 遇到的问题

| 问题 | 解决方案 |
|------|---------|
| 规划恢复脚本路径在当前环境不存在 | 改为直接读取仓库文件并用最新构建结果校准进度 |

## 资源

- 通义千问 API：https://help.aliyun.com/zh/dashscope/
- Qdrant Go client：github.com/qdrant/go-client
- go-zero 文档：https://go-zero.dev/docs/
- go-zero-looklook（参考项目，本人本地实操部署过）：https://github.com/Mikaelemmmm/go-zero-looklook

---
*每执行2次查看/搜索操作后更新此文件*
