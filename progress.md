# 进度日志

## 会话：2026-04-21

### 阶段 1-5：代码与部署主体
- **状态：** complete
- **开始时间：** 2026-04-21
- 执行的操作：
  - 创建项目目录 /Users/zp/GIT/k8s-deploy/go-zero-rag
  - 创建规划文件：task_plan.md、findings.md、progress.md、AGENTS.md
  - 生成并完善 `docservice`、`aiservice`、`gateway` 三个服务骨架
  - 完成文档上传、分块、Embedding、Qdrant 写入、检索和问答链路代码
  - 完成 Dockerfile、docker-compose、K8s YAML 和 Jenkins Pipeline
- 创建/修改的文件：
  - AGENTS.md
  - task_plan.md
  - findings.md
  - progress.md
  - docservice/
  - aiservice/
  - gateway/
  - pkg/
  - deploy/
  - docker-compose.yml
- 2026-04-21 完成：创建 docker-compose.yml（etcd 3.6.4 + Qdrant 1.9.2）
- 2026-04-21 完成：本地重新执行 `go build ./...`，编译通过
- 未完成项：
  - 缺少 `QWEN_API_KEY`，尚未完成本地问答闭环验证
  - 缺少集群权限与线上密钥，尚未完成 K8s 实际发布验证
- **阶段 1-5 全部完成 ✅**

### 阶段 6：README + GitHub public
- **状态：** complete ✅
- **开始时间：** 2026-04-21
- **完成时间：** 2026-04-23
- 执行的操作：
  - 修复 README 中绝对路径死链（Jenkins 链接）
  - 脱敏内网 IP（172.16.124.180、192.168.1.81），替换为 placeholder
  - 创建 .gitignore，排除内部文件（AGENTS.md、docs/、task_plan.md 等）
  - 在 go-zero-rag/ 内初始化独立 git 仓库并推送至 GitHub
  - Playwright 验证：Mermaid 架构图正常渲染，README 排版正常
- 创建/修改的文件：
  - .gitignore（新建）
  - README.md（死链修复）
  - deploy/jenkins/Jenkinsfile（内网 IP 脱敏）
  - deploy/k8s/docservice.yaml / aiservice.yaml / gateway.yaml（镜像地址脱敏）
- **GitHub 仓库：** https://github.com/suiran17/go-zero-rag（public）

## 会话：2026-06-18

### 模型切换：本地 Embedding + DeepSeek LLM
- **状态：** complete ✅
- **背景：** 避免依赖通义千问云端 key，改用本地 Embedding 模型 + DeepSeek 生成
- 执行的操作：
  - `pkg/embed`：原 DashScope 专用客户端 → OpenAI 兼容 `/v1/embeddings` 客户端（支持无 key 的本地服务）
  - `pkg/llm`：原通义千问 chat → 可配置 baseURL 的 OpenAI 兼容 chat 客户端（DeepSeek）
  - `pkg/qdrantcli`：向量维度由硬编码 `1536` 改为可配置，默认 `1024`
  - 两个 `config.go`：新增 `EmbedBaseURL/EmbedModel/EmbedApiKey/VectorSize`，aiservice 加 `LLMBaseURL/LLMModel/LLMApiKey`
  - 两个 `servicecontext.go`：字段 `Qwen` → `Embedder`，接线新客户端
  - 两个 `etc/*.yaml`：指向本地 `bge-m3`（`http://localhost:1234/v1`，1024 维）+ DeepSeek `deepseek-chat`
  - `aiservice/ai.go`：加 `conf.UseEnv()`，使 `${DEEPSEEK_API_KEY}` 从环境变量展开（修复了原本 `${...}` 从未展开的隐藏 bug）
  - README、K8s（configmap/secret/deployment）全部对齐新配置
  - `.gitignore`：排除含明文 key 的 `一些配置.md`
- 验证：
  - `curl` 实测本地 `bge-m3`：返回 1024 维向量 ✓
  - `curl` 实测 DeepSeek `deepseek-chat`：key 有效、接口格式正确 ✓
  - `go build ./...` 通过 ✓
- 未完成项：本地三服务端到端闭环、K8s 线上发布（需可达的 Embedding 地址）

## 测试结果

| 测试 | 输入 | 预期结果 | 实际结果 | 状态 |
|------|------|---------|---------|------|
| `go build ./...` | 当前仓库全部 Go 代码 | 全部编译通过 | 编译通过，无输出 | pass |
| README 文件检查 | README 需存在且覆盖运行说明 | 存在 README.md | 已创建并写入关键内容 | pass |
| 本地 Embedding 接口 | bge-m3 @ localhost:1234 | 返回向量 | 返回 1024 维向量 | pass |
| DeepSeek Chat 接口 | deepseek-chat + 真实 key | 返回答案 | key 有效，接口格式正确 | pass |
| 本地端到端问答 | 上传文档并发起提问 | 返回答案与来源 | 2026-06-18 三服务整链路实测打通，5 用例全 pass（见 端到端测试报告.md） | pass |
| K8s 实际发布验证 | 应用集群资源并检查 rollout | Pod 正常运行、Gateway 可访问 | 受集群权限 + Embedding 可达性限制，未执行 | blocked |
| GitHub public 发布 | https://github.com/suiran17/go-zero-rag | public 可访问，Mermaid 渲染正常 | Playwright 验证通过 | pass |

## 错误日志

| 时间戳 | 错误 | 尝试次数 | 解决方案 |
|--------|------|---------|---------|
| 2026-04-21 | `session-catchup.py` 路径不存在 | 1 | 直接读取仓库状态并补齐规划文件 |

## 五问重启检查

| 问题 | 答案 |
|------|------|
| 我在哪里？ | 全部阶段完成 ✅ |
| 我要去哪里？ | 补充真实 QWEN_API_KEY 后可做本地端到端验证（可选） |
| 目标是什么？ | go-zero RAG 微服务，K8s 部署，GitHub public 求职作品 |
| 我学到了什么？ | 见 findings.md |
| 我做了什么？ | 完成全部 6 个阶段，仓库已 public：https://github.com/suiran17/go-zero-rag |

---
*每个阶段完成后更新此文件*
