<h1 align="center">
  <img src="./internal/github.com/chuanbosi666/agent_go.png" alt="github.com/chuanbosi666/agent_go Logo" width="200">
</h1>
<h2 align="center">
  <b>github.com/chuanbosi666/agent_go: Elegant multi-agent workflows in Go</b>
</h2>

<p align="center">
  <a href="#特性">特性</a> •
  <a href="#安装">安装</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#核心概念">核心概念</a> •
  <a href="#示例">示例</a> •
  <a href="#api-参考">API</a>
</p>

---

**github.com/chuanbosi666/agent_go** 是一个 Go 语言的 AI Agent 框架，参考 [OpenAI Agents SDK](https://github.com/openai/openai-agents-python) 设计。它提供了构建智能 Agent 应用所需的全部功能。

## 特性

- **Agent 定义** - 使用 Builder 模式轻松配置 Agent
- **工具调用** - 定义 FunctionTool 让 Agent 执行具体操作
- **MCP 集成** - 支持 [Model Context Protocol](https://modelcontextprotocol.io) 工具
- **Guardrails** - 输入/输出护栏确保安全
- **Session** - SQLite 会话管理，支持多轮对话
- **多 Agent 协作** - Agent-as-Tool 模式实现分工协作
- **ReAct 模式** - 内置推理-行动循环支持
- **工具路由** - 动态选择相关工具，优化性能

## 安装

```bash
go get github.com/agent_go
```

要求：
- Go 1.25+
- OpenAI API Key

## 快速开始

### 1. 设置环境变量

```bash
export OPENAI_API_KEY=sk-xxx
```

或创建 `.env` 文件：

```env
OPENAI_API_KEY=sk-xxx
```

### 2. 创建简单的 Agent

```go
package main

import (
    "context"
    "fmt"
    "log"

    agentgo "github.com/agent_go"
    "github.com/openai/openai-go/v3"
)

func main() {
    // 创建 OpenAI 客户端
    client := openai.NewClient()

    // 创建 Agent
    agent := github.com/chuanbosi666/agent_go.New("助手").
        WithInstructions("你是一个友好的 AI 助手").
        WithModel("gpt-4o-mini").
        WithClient(client)

    // 运行
    ctx := context.Background()
    result, err := github.com/chuanbosi666/agent_go.Run(ctx, agent, "你好！")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result.FinalOutput)
}
```

### 3. 使用工具

```go
// 定义工具
calculator := github.com/chuanbosi666/agent_go.FunctionTool{
    Name:        "calculator",
    Description: "执行数学计算",
    ParamsJSONSchema: map[string]any{
        "type": "object",
        "properties": map[string]any{
            "expression": map[string]any{
                "type":        "string",
                "description": "要计算的表达式",
            },
        },
        "required": []string{"expression"},
    },
    OnInvokeTool: func(ctx context.Context, args string) (any, error) {
        // 实现计算逻辑
        return "42", nil
    },
}

// 添加到 Agent
agent := github.com/chuanbosi666/agent_go.New("计算助手").
    WithInstructions("你是一个计算助手").
    WithModel("gpt-4o-mini").
    WithClient(client).
    WithTools([]github.com/chuanbosi666/agent_go.FunctionTool{calculator})
```

### 4. 多轮对话

```go
// 创建 Session
session, _ := memory.NewSQLiteSession(ctx, memory.SQLiteSessionConfig{
    SessionID: "chat-001",
    DBPath:    ":memory:", // 或文件路径
})
defer session.Close()

// 配置 Runner
runner := github.com/chuanbosi666/agent_go.Runner{
    Config: github.com/chuanbosi666/agent_go.RunConfig{
        Session: session,
    },
}

// 多轮对话
runner.Run(ctx, agent, "我叫小明")
runner.Run(ctx, agent, "你还记得我的名字吗？") // Agent 会记住
```

## 核心概念

### Agent

Agent 是 AI 模型的配置，包含：

| 属性 | 说明 |
|------|------|
| Name | Agent 名称 |
| Instructions | 系统提示词 |
| Model | 使用的模型（如 gpt-4o） |
| Tools | 可用的工具列表 |
| InputGuardrails | 输入护栏 |
| OutputGuardrails | 输出护栏 |

### Runner

Runner 执行 Agent 的主循环：

1. 调用 LLM 获取响应
2. 如果有工具调用，执行工具
3. 将结果反馈给 LLM
4. 重复直到生成最终输出

```go
runner := github.com/chuanbosi666/agent_go.Runner{
    Config: github.com/chuanbosi666/agent_go.RunConfig{
        MaxTurns: 10,          // 最大执行轮次
        Session:  session,      // 会话管理
        ToolRouter: router,     // 工具路由
    },
}
```

### FunctionTool

FunctionTool 定义 Agent 可调用的工具：

```go
tool := github.com/chuanbosi666/agent_go.FunctionTool{
    Name:             "tool_name",
    Description:      "工具描述",
    ParamsJSONSchema: schema,
    OnInvokeTool:     func(ctx context.Context, args string) (any, error) {
        // 执行逻辑
    },
}
```

### Guardrails

护栏用于检查输入输出的安全性：

```go
// 输入护栏
inputGuardrail := github.com/chuanbosi666/agent_go.NewInputGuardrail("check_input",
    func(ctx context.Context, agent *github.com/chuanbosi666/agent_go.Agent, input github.com/chuanbosi666/agent_go.Input) (github.com/chuanbosi666/agent_go.GuardrailFunctionOutput, error) {
        // 检查输入
        return github.com/chuanbosi666/agent_go.GuardrailFunctionOutput{
            TripwireTriggered: false, // true 会停止执行
        }, nil
    })

// 输出护栏
outputGuardrail := github.com/chuanbosi666/agent_go.NewOutputGuardrail("check_output",
    func(ctx context.Context, agent *github.com/chuanbosi666/agent_go.Agent, output any) (github.com/chuanbosi666/agent_go.GuardrailFunctionOutput, error) {
        // 检查输出
        return github.com/chuanbosi666/agent_go.GuardrailFunctionOutput{
            TripwireTriggered: false,
        }, nil
    })
```

### Session

Session 管理对话历史：

```go
session, _ := memory.NewSQLiteSession(ctx, memory.SQLiteSessionConfig{
    SessionID:     "session-001",
    DBPath:        ":memory:",      // 内存存储
    // DBPath:     "./history.db",  // 文件存储
})
```

## 示例

查看 [examples](./examples/) 目录获取完整示例：

| 示例 | 说明 |
|------|------|
| [01-basic](./examples/01-basic/) | 基础用法 |
| [02-tools](./examples/02-tools/) | 工具调用 |
| [03-multi-agent](./examples/03-multi-agent/) | 多 Agent 协作 |
| [04-react](./examples/04-react/) | ReAct 模式 |
| [05-guardrails](./examples/05-guardrails/) | 护栏功能 |
| [06-session](./examples/06-session/) | 会话管理 |

运行示例：

```bash
cd examples/01-basic
go run main.go
```

## API 参考

### Agent Builder 方法

```go
github.com/chuanbosi666/agent_go.New("name")                        // 创建 Agent
    .WithInstructions("prompt")         // 设置系统提示词
    .WithModel("gpt-4o")                // 设置模型
    .WithClient(client)                 // 设置 OpenAI 客户端
    .WithTools(tools)                   // 设置工具
    .WithInputGuardrails(guards)        // 设置输入护栏
    .WithOutputGuardrails(guards)       // 设置输出护栏
    .WithMCPServers(servers)            // 设置 MCP 服务器
    .WithModelSettings(settings)        // 设置模型参数
```

### Runner 配置

```go
github.com/chuanbosi666/agent_go.RunConfig{
    MaxTurns:             10,       // 最大轮次
    Session:              session,  // 会话
    ToolRouter:           router,   // 工具路由
    ToolRoutingThreshold: 5,        // 路由阈值
    InputGuardrails:      guards,   // 全局输入护栏
    OutputGuardrails:     guards,   // 全局输出护栏
}
```

### 高级功能

**Agent-as-Tool（多 Agent 协作）**

```go
subAgent := github.com/chuanbosi666/agent_go.New("专家").WithInstructions("...")
tool := github.com/chuanbosi666/agent_go.WrapAgentAsTool(subAgent, 5)
mainAgent.WithTools([]github.com/chuanbosi666/agent_go.FunctionTool{tool})
```

**ReAct 模式**

```go
agent := github.com/chuanbosi666/agent_go.New("ReAct").
    WithInstructionsGetter(github.com/chuanbosi666/agent_go.DefaultReActInstruction)
```

**工具路由**

```go
router := &github.com/chuanbosi666/agent_go.KeywordRouter{
    ToolKeywords: map[string][]string{
        "calculator": {"计算", "加", "减", "乘", "除"},
        "weather":    {"天气", "温度", "下雨"},
    },
    TopN: 3,
}
runner.Config.ToolRouter = router
```

## 项目结构

```
github.com/chuanbosi666/agent_go/
├── agent.go          # Agent 定义
├── runner.go         # 执行引擎
├── tool.go           # 工具接口
├── tool_router.go    # 工具路由
├── agent_tool.go     # Agent-as-Tool
├── guardrail.go      # 护栏
├── react.go          # ReAct 模式
├── mcp.go            # MCP 集成
├── instruction.go    # 动态指令
├── memory/           # 会话管理
│   ├── session.go
│   └── sqlite.go
└── examples/         # 示例代码
```

## 开发

```bash
# 构建
go build ./...

# 测试
go test -v ./...

# 格式化
go fmt ./...
```

## Acknowledgements

- [OpenAI Agents Python](https://github.com/openai/openai-agents-python) - 设计参考
- [OpenAI Go SDK](https://github.com/openai/openai-go) - OpenAI API 客户端
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - MCP 协议支持
- [NeMo Agent Toolkit](https://github.com/NVIDIA/NeMo-Agent-Toolkit) - 灵感来源

## License

MIT License
