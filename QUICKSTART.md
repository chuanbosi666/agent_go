# nvgo 快速启动指南

## 📋 项目状态

✅ **项目已就绪，可以运行！**

所有核心功能已实现并通过测试：
- ✅ 核心 Agent 框架
- ✅ 工具调用（Function Tools + MCP）
- ✅ 输入/输出护栏
- ✅ 会话管理（SQLite）
- ✅ 多 Agent 协作
- ✅ ReAct 模式
- ✅ 工具路由
- ✅ 所有示例可编译运行

## 🚀 快速开始（3 步）

### 步骤 1：配置 API

**选项 A：使用 OpenAI API**

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，填入你的 API Key
# OPENAI_API_KEY=sk-xxx
```

**选项 B：使用自定义 API（推荐给你）**

```bash
# 编辑 .env 文件，配置自定义 endpoint
# OPENAI_API_KEY=your-api-key
# OPENAI_BASE_URL=https://your-api-endpoint.com/v1
```

### 步骤 2：运行示例

```bash
# 基础示例
cd examples/01-basic
go run main.go

# 自定义 endpoint 示例（适合你）
cd examples/07-custom-endpoint
# 修改 main.go 中的 baseURL 和 apiKey
go run main.go
```

### 步骤 3：开始开发

```go
package main

import (
    "context"
    "log"

    nvgo "nvgo"
    "github.com/openai/openai-go/v3"
    "github.com/openai/openai-go/v3/option"
)

func main() {
    // 创建自定义客户端
    client := openai.NewClient(
        option.WithAPIKey("your-key"),
        option.WithBaseURL("https://your-endpoint/v1"),
    )

    // 创建 Agent
    agent := nvgo.New("助手").
        WithInstructions("你是一个友好的 AI 助手").
        WithModel("your-model-name").
        WithClient(client)

    // 运行
    result, _ := nvgo.Run(context.Background(), agent, "你好！")
    log.Println(result.FinalOutput)
}
```

## 📦 依赖说明

所有依赖已在 go.mod 中配置：

```
github.com/openai/openai-go/v3        # OpenAI SDK
github.com/modelcontextprotocol/go-sdk # MCP 支持
github.com/mattn/go-sqlite3            # SQLite 会话存储
```

运行 `go mod download` 自动下载。

## 🔧 支持的 API

nvgo 支持任何 **OpenAI 兼容的 API**，包括：

| 提供商 | 说明 | 配置方式 |
|--------|------|---------|
| **OpenRouter** ⭐ | **支持 100+ 模型**（推荐）<br>Claude, Gemini, Grok, DeepSeek, Kimi, GLM 等 | `OPENAI_BASE_URL=https://openrouter.ai/api/v1`<br>详见 [多模型指南](docs/MULTI_MODEL_GUIDE.md) |
| OpenAI | 官方 GPT 系列 | `OPENAI_BASE_URL=https://api.openai.com/v1` |
| Azure OpenAI | Azure 托管的 OpenAI | `OPENAI_BASE_URL=https://<resource>.openai.azure.com` |
| 本地 Ollama | 本地运行开源模型 | `OPENAI_BASE_URL=http://localhost:11434/v1` |
| LM Studio | 本地 GUI 工具 | `OPENAI_BASE_URL=http://localhost:1234/v1` |
| LiteLLM Proxy | 本地代理各种 API | `OPENAI_BASE_URL=http://localhost:4000` |

### 🌟 推荐：使用 OpenRouter 访问所有模型

**一次配置，支持 100+ 模型，无需额外代码！**

```bash
# .env 配置
OPENAI_API_KEY=sk-or-v1-your-openrouter-key
OPENAI_BASE_URL=https://openrouter.ai/api/v1
```

```go
// 切换模型只需改一行！
agent := nvgo.New("助手").
    WithModel("anthropic/claude-3.5-sonnet"). // Claude
    // WithModel("google/gemini-flash-1.5").  // Gemini
    // WithModel("deepseek/deepseek-chat").   // DeepSeek
    // WithModel("moonshot/moonshot-v1-8k").  // Kimi
    // WithModel("zhipuai/glm-4").            // GLM
    WithClient(client)
```

📖 **完整模型列表和选择指南**: [docs/MULTI_MODEL_GUIDE.md](docs/MULTI_MODEL_GUIDE.md)

## 📚 示例列表

| 示例 | 说明 | 路径 |
|------|------|------|
| 01-basic | 最简单的 Agent 使用 | `examples/01-basic/` |
| 02-tools | 工具调用功能 | `examples/02-tools/` |
| 03-multi-agent | 多 Agent 协作 | `examples/03-multi-agent/` |
| 04-react | ReAct 推理模式 | `examples/04-react/` |
| 05-guardrails | 输入输出护栏 | `examples/05-guardrails/` |
| 06-session | 会话管理 | `examples/06-session/` |
| **07-custom-endpoint** | **自定义 API endpoint** | `examples/07-custom-endpoint/` |

## 🎯 下一步

1. **阅读文档**：查看 [README.md](README.md) 了解详细 API
2. **运行测试**：`go test -v ./...` 验证环境
3. **查看示例**：浏览 `examples/` 目录学习用法
4. **开始开发**：基于示例创建你自己的 Agent

## ⚙️ 配置说明

### 环境变量

| 变量 | 说明 | 必需 |
|------|------|------|
| `OPENAI_API_KEY` | API 密钥 | ✅ 是 |
| `OPENAI_BASE_URL` | 自定义 endpoint | ❌ 否 |
| `OPENAI_ORG_ID` | 组织 ID | ❌ 否 |

### 代码配置

```go
// 在代码中配置（覆盖环境变量）
client := openai.NewClient(
    option.WithAPIKey("key"),
    option.WithBaseURL("url"),
    option.WithOrganization("org-id"),
)
```

## 🐛 常见问题

### Q: 示例编译失败？
A: 确保在项目根目录运行 `go mod download`

### Q: 运行时报错 "no API key"？
A: 检查 `.env` 文件或环境变量配置

### Q: 如何使用本地模型？
A: 参考 `examples/07-custom-endpoint/`，设置 `OPENAI_BASE_URL` 为本地服务地址

### Q: 支持哪些模型？
A: 支持任何 OpenAI Chat Completions API 兼容的模型

## 📊 项目结构

```
nvgo-main/
├── nvgo.go              # 主入口（导出所有 API）
├── pkg/                 # 核心包
│   ├── agent/          # Agent 定义
│   ├── runner/         # 执行引擎
│   ├── tool/           # 工具接口
│   ├── types/          # 共享类型
│   ├── pattern/        # 设计模式
│   └── memory/         # 会话管理
├── examples/            # 示例代码
├── .env.example        # 环境变量模板
└── README.md           # 完整文档
```

## 🎉 开始使用

```bash
# 1. 配置环境
cp .env.example .env
# 编辑 .env，填入你的配置

# 2. 运行示例
cd examples/07-custom-endpoint
go run main.go

# 3. 开始开发！
```

需要帮助？查看 [README.md](README.md) 或提交 Issue。
