# 多模型支持指南

github.com/chuanbosi666/agent_go 通过 OpenRouter 统一接口，支持 **100+ 种 AI 模型**，包括所有主流厂商。

## 🚀 快速开始

### 1. 获取 OpenRouter API Key

1. 访问 [OpenRouter](https://openrouter.ai/)
2. 注册并登录
3. 前往 [Keys 页面](https://openrouter.ai/keys) 创建 API Key
4. 充值一定金额（按使用量计费，支持信用卡/加密货币）

### 2. 配置环境变量

```bash
# .env 文件
OPENAI_API_KEY=sk-or-v1-your-openrouter-api-key
OPENAI_BASE_URL=https://openrouter.ai/api/v1
```

### 3. 使用任意模型

```go
import (
    github.com/chuanbosi666/agent_go "github.com/chuanbosi666/agent_go"
    "github.com/openai/openai-go/v3"
)

client := openai.NewClient()  // 自动读取环境变量

// 只需要改模型名，其他完全一样！
agent := github.com/chuanbosi666/agent_go.New("助手").
    WithModel("anthropic/claude-3.5-sonnet").  // 改这一行切换模型
    WithClient(client)
```

## 📋 支持的模型列表

### Claude 系列（Anthropic）- 最强推理

| 模型名 | OpenRouter ID | 特点 | 价格 |
|--------|--------------|------|------|
| Claude 3.5 Sonnet | `anthropic/claude-3.5-sonnet` | 最强推理、长上下文 | $3/$15 |
| Claude 3.5 Haiku | `anthropic/claude-3.5-haiku` | 快速、便宜 | $0.8/$4 |
| Claude 3 Opus | `anthropic/claude-3-opus` | 旧版最强 | $15/$75 |

**推荐场景**：复杂推理、代码分析、长文本理解

---

### Gemini 系列（Google）- 多模态

| 模型名 | OpenRouter ID | 特点 | 价格 |
|--------|--------------|------|------|
| Gemini 2.0 Flash Exp | `google/gemini-2.0-flash-exp` | 最新实验版 | 免费 |
| Gemini 1.5 Pro | `google/gemini-pro-1.5` | 超长上下文 200 万 token | $1.25/$5 |
| Gemini 1.5 Flash | `google/gemini-flash-1.5` | 快速、便宜 | $0.075/$0.3 |

**推荐场景**：长文档、视频分析、快速对话

---

### Grok 系列（xAI/Twitter）- 实时联网

| 模型名 | OpenRouter ID | 特点 | 价格 |
|--------|--------------|------|------|
| Grok Beta | `x-ai/grok-beta` | 实时 X 数据 | $5/$10 |
| Grok 2 | `x-ai/grok-2-1212` | 最新稳定版 | $2/$10 |

**推荐场景**：需要实时信息、社交媒体分析

---

### DeepSeek 系列（国产）- 超高性价比

| 模型名 | OpenRouter ID | 特点 | 价格 |
|--------|--------------|------|------|
| DeepSeek V3 | `deepseek/deepseek-chat` | 对话模型，极致性价比 | $0.14/$0.28 |
| DeepSeek Coder | `deepseek/deepseek-coder` | 代码专用 | $0.14/$0.28 |

**推荐场景**：预算有限、大量调用、代码生成

---

### Kimi 系列（Moonshot AI，国产）- 超长上下文

| 模型名 | OpenRouter ID | 特点 | 价格 |
|--------|--------------|------|------|
| Kimi | `moonshot/moonshot-v1-8k` | 中文优化、长上下文 | $0.12/$0.12 |

**推荐场景**：中文处理、长文档

---

### GLM 系列（智谱 AI，国产）- 中文优化

| 模型名 | OpenRouter ID | 特点 | 价格 |
|--------|--------------|------|------|
| GLM-4 Plus | `zhipuai/glm-4-plus` | 最新版 | $7.15/$7.15 |
| GLM-4 | `zhipuai/glm-4` | 标准版 | $1/$1 |

**推荐场景**：中文任务、企业应用

---

### GPT 系列（OpenAI）- 标杆模型

| 模型名 | OpenRouter ID | 特点 | 价格 |
|--------|--------------|------|------|
| GPT-4o | `openai/gpt-4o` | 多模态、最新 | $2.5/$10 |
| GPT-4o Mini | `openai/gpt-4o-mini` | 快速、便宜 | $0.15/$0.6 |
| o1 | `openai/o1` | 强化推理 | $15/$60 |

**推荐场景**：综合能力、作为对比基准

---

### 其他推荐模型

| 厂商 | 模型 | OpenRouter ID | 特点 |
|------|------|--------------|------|
| Meta | Llama 3.3 70B | `meta-llama/llama-3.3-70b-instruct` | 开源、免费 |
| Anthropic | Claude Instant | `anthropic/claude-instant-1.2` | 便宜、快速 |
| Cohere | Command R+ | `cohere/command-r-plus` | 企业级 |
| Mistral | Mistral Large | `mistral/mistral-large` | 欧洲模型 |

## 🎯 选择建议

### 按任务类型选择

| 任务类型 | 推荐模型 | 原因 |
|---------|---------|------|
| **复杂推理** | Claude 3.5 Sonnet / GPT-4o | 逻辑能力强 |
| **快速对话** | Gemini Flash / DeepSeek Chat | 快速、便宜 |
| **代码生成** | DeepSeek Coder / Claude Sonnet | 代码能力强 |
| **中文任务** | Kimi / GLM-4 / DeepSeek | 中文优化 |
| **长文本** | Gemini Pro (200 万 token) | 超长上下文 |
| **预算有限** | DeepSeek / Gemini Flash | 性价比高 |

### 按预算选择

| 预算 | 推荐组合 |
|------|---------|
| **土豪** | Claude Opus + GPT-4o |
| **平衡** | Claude Sonnet + Gemini Flash |
| **省钱** | DeepSeek + Gemini Flash |
| **免费** | Gemini 2.0 Flash Exp / Llama 3.3 |

## 📝 完整示例

```go
package main

import (
    "context"
    "log"

    github.com/chuanbosi666/agent_go "github.com/chuanbosi666/agent_go"
    "github.com/openai/openai-go/v3"
)

func main() {
    client := openai.NewClient()  // 从环境变量读取配置

    // 示例：使用 DeepSeek（性价比最高）
    agent := github.com/chuanbosi666/agent_go.New("助手").
        WithInstructions("你是一个编程助手").
        WithModel("deepseek/deepseek-chat").
        WithClient(client)

    result, err := github.com/chuanbosi666/agent_go.Run(context.Background(), agent, "写一个快速排序")
    if err != nil {
        log.Fatal(err)
    }

    log.Println(result.FinalOutput)
}
```

## 🔄 动态切换模型

```go
// 根据配置文件切换
modelConfig := map[string]string{
    "default":  "deepseek/deepseek-chat",      // 默认用便宜的
    "complex":  "anthropic/claude-3.5-sonnet", // 复杂任务用强的
    "fast":     "google/gemini-flash-1.5",     // 快速响应用 Flash
}

func createAgent(taskType string) *github.com/chuanbosi666/agent_go.Agent {
    model := modelConfig[taskType]
    return github.com/chuanbosi666/agent_go.New("助手").
        WithModel(model).
        WithClient(client)
}
```

## 💰 价格说明

OpenRouter 的价格格式：`$输入/$输出` per 1M tokens

**示例计算**：
- DeepSeek ($0.14/$0.28)：1000 次对话（每次 1k input + 1k output）= $0.42
- Claude Sonnet ($3/$15)：同样的使用 = $18
- 差价约 **40 倍**

💡 **省钱技巧**：
1. 简单任务用 DeepSeek / Gemini Flash
2. 复杂任务才用 Claude / GPT-4o
3. 使用工具路由自动选择模型

## 🔗 有用链接

- [OpenRouter 模型列表](https://openrouter.ai/models)
- [价格对比](https://openrouter.ai/models?order=newest)
- [使用统计](https://openrouter.ai/activity)
- [API 文档](https://openrouter.ai/docs)

## ❓ 常见问题

### Q: OpenRouter 会不会比直接调用贵很多？
A: 只贵 5-10%，但省去了：
- 管理多个 API Key
- 维护多个 SDK
- 写适配器代码
- 处理不同的限流策略

### Q: 可以同时用官方 API 和 OpenRouter 吗？
A: 可以！只需要创建不同的 client：

```go
// OpenRouter client
orClient := openai.NewClient(
    option.WithAPIKey("sk-or-v1-..."),
    option.WithBaseURL("https://openrouter.ai/api/v1"),
)

// 官方 OpenAI client
openaiClient := openai.NewClient(
    option.WithAPIKey("sk-..."),
)
```

### Q: 支持流式输出吗？
A: 支持！但 github.com/chuanbosi666/agent_go 当前版本未实现，需要后续添加。

### Q: 有免费额度吗？
A: OpenRouter 本身无免费额度，但部分模型免费：
- Gemini 2.0 Flash Exp
- Llama 3.3 70B
- 等实验性模型

## ✅ 总结

使用 OpenRouter：
- ✅ **一次配置，支持 100+ 模型**
- ✅ **切换模型只需改一行代码**
- ✅ **统一的监控和计费**
- ✅ **无需维护多个 SDK**
- ✅ **支持你要的所有模型**：Claude、Gemini、Grok、DeepSeek、Kimi、GLM

**推荐起步配置**：
```bash
# .env
OPENAI_API_KEY=sk-or-v1-your-key
OPENAI_BASE_URL=https://openrouter.ai/api/v1

# 默认用 DeepSeek（便宜）
DEFAULT_MODEL=deepseek/deepseek-chat
```
