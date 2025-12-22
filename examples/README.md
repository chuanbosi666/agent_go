# nvgo 示例代码

本目录包含 nvgo 框架的各种使用示例，帮助你快速上手。

## 示例列表

| 目录 | 说明 | 难度 |
|------|------|------|
| [01-basic](./01-basic/) | 基础用法 - 创建 Agent 并对话 | ⭐ 入门 |
| [02-tools](./02-tools/) | 工具调用 - 定义和使用 FunctionTool | ⭐⭐ 基础 |
| [03-multi-agent](./03-multi-agent/) | 多 Agent 协作 - Agent-as-Tool 模式 | ⭐⭐⭐ 进阶 |
| [04-react](./04-react/) | ReAct 模式 - 思考-行动-观察循环 | ⭐⭐⭐ 进阶 |
| [05-guardrails](./05-guardrails/) | 护栏功能 - 输入输出安全检查 | ⭐⭐ 基础 |
| [06-session](./06-session/) | 会话管理 - 多轮对话和历史存储 | ⭐⭐ 基础 |

## 快速开始

### 1. 环境准备

```bash
# 设置 OpenAI API Key
export OPENAI_API_KEY=sk-xxx
```

### 2. 运行示例

```bash
# 运行基础示例
cd examples/01-basic
go run main.go

# 运行工具调用示例
cd examples/02-tools
go run main.go

# 运行多 Agent 示例
cd examples/03-multi-agent
go run main.go

# 运行 ReAct 模式示例
cd examples/04-react
go run main.go

# 运行护栏示例
cd examples/05-guardrails
go run main.go

# 运行会话示例（演示模式）
cd examples/06-session
go run main.go --demo

# 运行会话示例（交互模式）
cd examples/06-session
go run main.go
```

## 示例详解

### 01-basic: 基础用法

最简单的 nvgo 使用示例，展示如何：
- 创建 OpenAI 客户端
- 使用 Builder 模式配置 Agent
- 运行 Agent 并获取结果

```go
agent := nvgo.New("助手").
    WithInstructions("你是一个友好的 AI 助手").
    WithModel("gpt-4o-mini").
    WithClient(client)

result, err := nvgo.Run(ctx, agent, "你好！")
```

### 02-tools: 工具调用

展示如何定义 FunctionTool 让 Agent 执行具体操作：
- 定义工具的 JSON Schema
- 实现工具的执行逻辑
- Agent 自动选择和调用工具

### 03-multi-agent: 多 Agent 协作

展示 Agent-as-Tool 模式：
- 使用 `WrapAgentAsTool` 将 Agent 包装成工具
- 主 Agent 可以调用其他专业 Agent
- 构建分工协作的多 Agent 系统

### 04-react: ReAct 模式

展示 ReAct（Reasoning + Acting）推理模式：
- 使用内置的 `DefaultReActInstruction`
- Agent 按照思考->行动->观察循环解决问题
- 适合需要多步推理的复杂任务

### 05-guardrails: 护栏功能

展示安全护栏的使用：
- InputGuardrail：检查用户输入
- OutputGuardrail：检查 Agent 输出
- Tripwire 机制：触发时停止执行

### 06-session: 会话管理

展示 Session 功能：
- SQLiteSession：内存或文件存储
- 多轮对话上下文保持
- 对话历史管理

## 学习路径

建议按以下顺序学习：

1. **入门**: 01-basic → 理解基本概念
2. **工具**: 02-tools → 学习工具调用
3. **会话**: 06-session → 掌握多轮对话
4. **安全**: 05-guardrails → 了解安全机制
5. **进阶**: 03-multi-agent, 04-react → 高级模式

## 注意事项

1. 所有示例需要有效的 `OPENAI_API_KEY`
2. 默认使用 `gpt-4o-mini` 模型，可按需修改
3. 示例代码仅供学习，生产环境请添加错误处理
4. MCP 集成示例需要额外的 MCP 服务器配置
