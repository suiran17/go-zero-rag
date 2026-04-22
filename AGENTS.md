# 工作方式说明：using-superpowers + planning-with-files-zh

本项目使用两个 skill 协同工作，分别负责不同层面。

---

## 两个 Skill 的分工

### `/using-superpowers`
**作用：工作纪律层**

规定"做任何事之前先检查有没有适用的 skill"。
它是元规则，保证不会跳过流程直接动手。

核心要求：
- 用户消息到来 → 先检查是否有 skill 适用 → 再响应
- 哪怕只有 1% 的可能性适用，也必须先调用 skill

### `/planning-with-files-zh`
**作用：任务执行层**

规定"如何有条不紊地执行多步骤任务"。
用持久化 Markdown 文件代替易失的上下文记忆。

核心文件：
| 文件 | 用途 |
|------|------|
| `task_plan.md` | 阶段划分、进度、决策记录 |
| `findings.md` | 研究发现、技术调研结果 |
| `progress.md` | 每次会话的操作日志、错误记录 |

### SDD（Spec-Driven Development）
**作用：开发质量层**

在动手编码之前，先明确"要做什么"再明确"怎么做"，最后才执行。
三步强制顺序，防止在模糊需求上浪费实现成本。

| 阶段 | Skill | 产出位置 |
|------|-------|---------|
| 1. 写 Spec | `superpowers:brainstorming` | `docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md` |
| 2. 写 Plan | `superpowers:writing-plans` | `docs/superpowers/plans/YYYY-MM-DD-<feature>.md` |
| 3. 执行 | `superpowers:executing-plans` 或 `superpowers:subagent-driven-development` | 代码实现 |

核心规则：
- **Spec 未经用户确认，不得进入 Plan 阶段**（brainstorming skill 有 hard gate）
- **Plan 未写完，不得开始执行**
- Spec 变更时，Plan 必须同步更新，才能继续执行

---

## 搭配使用的完整流程

```
1. 用户发消息
        ↓
2. [using-superpowers] 检查是否有 skill 适用
        ↓
3. [SDD 新任务] brainstorming skill → 对话澄清需求 → 写 design.md → 用户确认
        ↓
4. [SDD 新任务] writing-plans skill → 把 spec 拆成可执行任务 → 写 plan.md
        ↓
5. [planning-with-files-zh] 读取 task_plan.md 恢复上下文
        ↓
6. [SDD] 执行当前阶段任务
        ↓
7. [planning-with-files-zh] 更新 task_plan.md / progress.md
        ↓
8. 响应用户
```

---

## 本项目的规划文件

- [docs/superpowers/specs/](docs/superpowers/specs/) — 设计文档，由 brainstorming skill 生成（SDD 第一步）
- [docs/superpowers/plans/](docs/superpowers/plans/) — 实现计划，由 writing-plans skill 生成（SDD 第二步）
- [task_plan.md](task_plan.md) — 项目阶段、进度、接口定义
- [findings.md](findings.md) — 技术决策、环境信息、RAG 流程
- [progress.md](progress.md) — 每次会话的操作日志

---

## 每次新会话的开始方式

告诉 Claude：
> "继续 go-zero-rag 项目"

Claude 会自动：
1. 读取 `task_plan.md` 确认当前阶段
2. 读取 `progress.md` 了解上次做了什么
3. 从断点继续，不重复已完成的工作

---

## 注意事项

- `task_plan.md` 只写**内部决策和计划**，不写外部网页内容（安全规则）
- 外部内容（API 文档、搜索结果）只写入 `findings.md`
- 每个阶段完成后立即更新状态：`pending → in_progress → complete`
