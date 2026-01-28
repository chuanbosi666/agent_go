# AI Agent 框架面试问答整理

> 本文档整理了针对 agentgo 项目的模拟面试问答，供后续复习使用。
>
> 面试岗位：AI Agent 工程师
>
> 整理日期：2026-01-11

---

## 目录

1. [架构设计：五层架构](#问题-1架构设计)
2. [核心引擎：Runner 执行流程](#问题-2runner-执行流程)
3. [工具系统：ToolRouter 工具路由](#问题-3toolrouter-工具路由)
4. [工具系统：MCP 协议集成](#问题-4mcp-协议集成)
5. [多智能体：Agent-as-Tool 协作](#问题-5agent-as-tool-多智能体协作)
6. [多模型支持：API 兼容性](#问题-6多模型支持)
7. [安全护栏：Tripwire 熔断机制](#问题-7tripwire-熔断机制)
8. [综合问题：改进方向与技术挑战](#问题-8改进方向与技术挑战)
9. [面试总结与建议](#面试总结)

---

## 问题 1：架构设计

### 面试官问题

> 你提到项目采用"执行层 → 能力层 → 基础设施层"的五层架构，但描述中只提到了三层。请具体说明这五层分别是什么？每一层的职责边界是如何划分的？当时为什么选择这种分层方式而不是其他架构模式（比如管道模式或事件驱动）？

### 标准答案

#### 五层架构详解

| 层级 | 名称 | 核心模块 | 职责 |
|-----|------|---------|------|
| 0 | 应用层 | examples/ | 用户应用、示例代码 |
| 1 | **执行层** | Runner | Agent Loop 主循环、多轮对话控制 |
| 2 | **能力层** | Agent, Tool, Guardrails, Pattern | 能力定义：工具、护栏、推理模式 |
| 3 | **数据层** | Session, Types, Config | 会话存储、类型定义、配置管理 |
| 4 | **基础设施层** | MCP, OpenAI SDK, SQLite | 外部依赖适配 |

#### 为什么选择分层而非管道/事件驱动？

1. **分层架构**：职责清晰，便于测试和替换（如 Session 可换实现）
2. **管道模式**：适合线性数据流，但 Agent Loop 是**循环结构**，不适合
3. **事件驱动**：增加复杂度，Agent 场景更适合同步调用链

#### 职责边界划分原则

- **执行层**只负责"怎么跑"，不关心工具细节
- **能力层**定义"能做什么"，不关心存储
- **数据层**只负责"存什么"，与业务逻辑解耦

#### 加分点

- 提到 **13 个核心接口**实现模块解耦
- 提到 **Builder 模式**让配置更灵活
- 提到这种设计便于**单元测试**（可 mock 各层接口）

---

## 问题 2：Runner 执行流程

### 面试官问题

> 你提到执行层的 Runner 负责 Agent Loop 主循环。请描述一下这个循环的核心流程：从用户输入进来，到最终输出返回，中间经过哪些关键步骤？如果某一步失败了（比如工具调用失败），Runner 是如何处理的？

### 标准答案

#### 核心执行流程

```
用户输入
    ↓
1. 输入护栏检查 (InputGuardrails)
    ↓ (通过)
2. 进入多轮循环 (turnCount < MaxTurns)
    ├─ 获取 Instructions (系统提示词)
    ├─ 获取工具列表 (Tools + MCP)
    ├─ 工具路由 (ToolRouter，可选)
    ├─ 加载会话历史 (Session)
    ↓
3. 调用 LLM API
    ├─ Responses API (OpenAI 专用)
    └─ Chat Completions API (通用)
    ↓
4. 处理响应
    ├─ 工具调用 → 执行工具 → 结果反馈给 LLM → 继续循环
    └─ 最终消息 → 提取 FinalOutput → 退出循环
    ↓
5. 输出护栏检查 (OutputGuardrails)
    ↓
返回 RunResult
```

#### 错误处理机制（重点！）

| 错误类型 | 处理方式 | 代码位置 |
|---------|---------|---------|
| **护栏触发** | 直接返回 `GuardrailTripwireTriggeredError`，**中断执行** | runner.go:149-155 |
| **工具调用失败** | 错误信息作为工具输出**反馈给 LLM**，让 LLM 决定下一步 | runner.go:250-257 |
| **工具未找到** | 同上，反馈 "Tool not found" 给 LLM | runner.go:241-247 |
| **超过 MaxTurns** | 返回 `MaxTurnsExceededError` | runner.go:342-344 |

#### 关键代码示例

```go
// 工具调用失败 → 反馈给 LLM（不是直接报错）
toolResult, err := executeTool(ctx, currentAgent, t, item.Arguments)
if err != nil {
    errorOutput := responses.ResponseInputItemParamOfFunctionCallOutput(
        item.CallID,
        fmt.Sprintf("Tool execution failed: %v", err),  // 错误信息作为工具输出
    )
    result.NewItems = append(result.NewItems, WrapRunItem(errorOutput))
    continue  // 继续处理其他工具调用
}
```

#### 循环终止条件

1. `result.FinalOutput != nil` — LLM 返回了最终消息
2. `turnCount >= maxTurns` — 达到最大轮数限制

---

## 问题 3：ToolRouter 工具路由

### 面试官问题

> 你提到了 ToolRouter 工具路由器，说它能"减少 30-50% Token 消耗"。请解释：工具路由的具体实现原理是什么？它是如何根据用户输入筛选相关工具的？这个 30-50% 的数字是怎么得出的？

### 标准答案

#### 核心原理：关键词匹配（不是 LLM！）

**重要**：KeywordRouter **不使用 LLM**，它是纯规则匹配，在 LLM 调用**之前**执行。

```
用户输入: "帮我计算 12 乘以 34"
         ↓
    extractText(input)
         ↓
    inputLower = "帮我计算 12 乘以 34"
         ↓
┌─────────────────────────────────────────────────────┐
│ 遍历每个工具，计算关键词匹配得分                      │
│                                                      │
│ Tool: calculator                                     │
│   keywords: ["计算", "加", "减", "乘", "除"]         │
│   匹配: "计算"✓ "乘"✓ → score = 2                   │
│                                                      │
│ Tool: weather                                        │
│   keywords: ["天气", "温度", "下雨"]                 │
│   匹配: 无 → score = 0                              │
└─────────────────────────────────────────────────────┘
         ↓
    按 score 降序排序 → 取 TopN (默认 5)
```

#### 核心代码

```go
// router.go:39-44
for _, keyword := range keywords {
    if strings.Contains(inputLower, strings.ToLower(keyword)) {
        score++
    }
}
```

#### 30-50% 数字的计算

| 场景 | 工具数 | 路由后 | Token 减少 |
|-----|-------|--------|-----------|
| 10 个工具，TopN=5 | 10 | 5 | ~50% |
| 20 个工具，TopN=5 | 20 | 5 | ~75% |
| 6 个工具，TopN=5 | 6 | 5 | ~17% |

每个工具的 JSON Schema 描述大约 100-300 tokens，工具越多节省越明显。

#### 面试标准表述

> "工具路由是一个**轻量级的预过滤机制**，在 LLM 调用之前执行，不消耗额外 Token。实现原理是：预先为每个工具定义关键词列表，然后用简单的字符串匹配计算每个工具与用户输入的相关度得分，最后按得分排序返回 TopN 个工具。"

---

## 问题 4：MCP 协议集成

### 面试官问题

> 你的框架支持 MCP 协议集成。请解释：MCP (Model Context Protocol) 是什么？你的框架是如何将外部 MCP 工具转换成内部 FunctionTool 格式的？转换过程中有哪些需要处理的兼容性问题？

### 标准答案

#### MCP 定义

**MCP** = Model Context Protocol（模型上下文协议），由 Anthropic 提出的开放标准，用于统一 AI 应用与外部工具/数据源的交互接口。

#### 转换流程

```
┌─────────────────────────────────────────────────────────────┐
│ 外部 MCP Server                                             │
│ (通过 Stdio 或 StreamableHTTP 连接)                          │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
              ListTools(ctx) → []*mcp.Tool
                       ↓
┌─────────────────────────────────────────────────────────────┐
│ MCPToolFilter 过滤                                           │
│ • AllowedToolNames: 白名单                                   │
│ • BlockedToolNames: 黑名单                                   │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────────────┐
│ ToFunctionTool() 格式转换  (mcp.go:110-143)                  │
│                                                              │
│ 1. Schema 转换:                                              │
│    mcp.Tool.InputSchema → map[string]any                    │
│                                                              │
│ 2. 确保 properties 字段存在                                  │
│                                                              │
│ 3. 严格模式转换 (可选):                                       │
│    strictschema.EnsureStrictJSONSchema(schema)              │
│                                                              │
│ 4. 包装调用函数:                                             │
│    OnInvokeTool = func() { InvokeMCPTool(...) }             │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
              FunctionTool (统一格式)
```

#### 兼容性问题

| 问题 | 代码位置 | 解决方案 |
|-----|---------|---------|
| **Schema 格式差异** | mcp.go:114-126 | JSON 序列化/反序列化转换，确保 `properties` 字段存在 |
| **严格模式要求** | mcp.go:128-135 | 使用 `strictschema.EnsureStrictJSONSchema` |
| **工具名称冲突** | mcp.go:81-85 | 遍历检查，发现重复名称返回错误 |
| **结果格式差异** | mcp.go:157-176 | 判断 `StructuredContent` vs `Content`，统一转 JSON |
| **两种传输方式** | mcp.go:328-380 | 抽象 `MCPServer` 接口，支持 Stdio 和 StreamableHTTP |

#### 面试标准表述

> "MCP 是 Model Context Protocol，由 Anthropic 提出的开放标准。我的框架通过 `MCPServer` 接口抽象了 MCP 服务器，支持 Stdio 和 StreamableHTTP 两种传输方式。
>
> 转换过程中的主要兼容性问题包括：JSON Schema 格式差异（需确保 properties 字段存在）、严格模式处理、结果格式统一、工具名称去重。"

---

## 问题 5：Agent-as-Tool 多智能体协作

### 面试官问题

> 你简历中提到了 Agent-as-Tool 多智能体协作机制，并举了一个 "dev-team 示例（5 个 Agent 协作完成软件开发全流程）"。请具体描述这个 dev-team 的架构设计。

### 标准答案

#### 5 个 Agent 角色

| 角色 | Agent 名称 | 职责 | 配置的工具 |
|-----|-----------|------|-----------|
| 产品经理 | REQ-Agent | 需求分析 | 无 |
| 架构师 | ARCH-Agent | 系统设计 | 文件读取 + 搜索（只读） |
| 程序员 | CODE-Agent | 代码实现 | 完整文件操作 + 搜索 |
| 测试员 | TEST-Agent | 测试验证 | 文件操作 + 执行 + 搜索 |
| 协调者 | Manager | 调度协调 | 4 个专家 Agent（作为工具） |

#### 协作机制

```
┌────────────────────────────────────────┐
│ Manager Agent (协调者)                  │
│ Instructions 定义四阶段工作流程         │
└──────────────┬─────────────────────────┘
               │ WrapAgentAsTool()
               ├────────────────────┬──────────────┬──────────────┐
               ↓                    ↓              ↓              ↓
        ┌─────────────┐     ┌───────────┐  ┌──────────┐  ┌────────────┐
        │ REQ-Agent   │     │ARCH-Agent │  │CODE-Agent│  │TEST-Agent  │
        │ (as-tool)   │     │ (as-tool) │  │(as-tool) │  │ (as-tool)  │
        └─────────────┘     └───────────┘  └──────────┘  └────────────┘
```

#### 协调者如何决定调用哪个专家？

通过 **Instructions 定义的工作流程 + LLM 自主判断**：

```
## 工作流程
### 第一阶段：需求分析 → 调用 REQ-Agent
### 第二阶段：架构设计 → 调用 ARCH-Agent
### 第三阶段：代码实现 → 调用 CODE-Agent
### 第四阶段：测试验证 → 调用 TEST-Agent
```

#### 错误处理

**两层错误处理**：
1. **Runner 层**：工具执行失败会反馈给 LLM
2. **Manager 层**：Instructions 定义了重试策略（最多 2 次）和兜底方案

```
## 错误处理
如果某个阶段失败：
1. 分析失败原因
2. 尝试修正输入后重试（最多 2 次）
3. 如果仍然失败，向用户报告问题并请求指导
```

#### 加分点：最小权限原则

不同专家 Agent 配置了不同的工具集：
- 架构师**只能读不能写**
- 只有程序员和测试员才能修改文件

---

## 问题 6：多模型支持

### 面试官问题

> 你简历中提到项目支持 "OpenAI、Claude、Gemini、DeepSeek、Qwen、本地 Ollama 等 100+ 模型"。请解释：你的框架是如何做到兼容这么多模型的？Responses API 和 Chat Completions API 有什么区别？

### 标准答案

#### 兼容原理

```
┌─────────────────────────────────────────────────────────────┐
│  你的框架 (agentgo)                                          │
│  ↓                                                          │
│  OpenAI Go SDK                                              │
│  ↓                                                          │
│  option.WithBaseURL("https://openrouter.ai/api/v1")        │
│  ↓                                                          │
│  OpenRouter (API 聚合网关)                                  │
│  ↓                                                          │
│  ┌─────┬─────┬─────┬─────┬─────┬─────┐                     │
│  │Claude│Gemini│GPT │DeepSeek│Qwen│Ollama│ ...100+ 模型    │
│  └─────┴─────┴─────┴─────┴─────┴─────┘                     │
└─────────────────────────────────────────────────────────────┘
```

#### 代码实现

```go
// 只需要改这两个参数，就能切换到任何兼容 API
client := openai.NewClient(
    option.WithAPIKey(apiKey),
    option.WithBaseURL(baseURL),  // 改这里
)

agent := agentgo.New("助手").
    WithModel("anthropic/claude-3.5-sonnet").  // 改这里
    WithClient(client)
```

#### 两种 API 对比

| 特性 | Responses API | Chat Completions API |
|-----|--------------|---------------------|
| **兼容性** | OpenAI 专用 | 通用标准（所有模型都支持） |
| **使用场景** | 需要高级功能 | 兼容第三方模型 |
| **请求结构** | `prompt` + `input` | `messages` 数组 |
| **Token 详情** | ✅ 有详细分类 | ❌ 只有总数 |
| **代码位置** | runner.go:372-446 | runner.go:448-508 |

#### API 选择逻辑

```go
// runner.go:217-230
if currentAgent.Prompt != nil {
    // OpenAI 专用：Responses API
    modelResponse, err = r.callResponsesAPI(...)
} else {
    // 通用：Chat Completions API
    modelResponse, err = r.callChatCompletionsAPI(...)
}
```

#### 面试标准表述

> "我的框架通过两个层面实现多模型兼容：
>
> **第一层：API 适配**。框架支持 Responses API（OpenAI 专用）和 Chat Completions API（业界通用标准）。
>
> **第二层：网关聚合**。通过 OpenRouter 这样的 API 聚合网关，只需改变 BaseURL 和模型名称，就能访问 100+ 模型。代码层面完全不需要改动。"

---

## 问题 7：Tripwire 熔断机制

### 面试官问题

> 你简历中提到了 "输入/输出双向安全护栏" 和 "Tripwire 熔断"。请解释：什么是 Tripwire 熔断机制？它和普通的错误处理有什么区别？

### 标准答案

#### Tripwire 定义

**Tripwire**（绊线）是一个安全术语，原意是"触发警报的细线"。

```go
// guardrail.go:10-15
type GuardrailFunctionOutput struct {
    TripwireTriggered bool  // 关键：是否触发熔断
    OutputInfo        any   // 附加信息
}
```

当 `TripwireTriggered = true` 时，**立即中断整个 Agent 执行**。

#### Tripwire vs 普通错误处理

| 对比项 | 普通错误处理 | Tripwire 熔断 |
|-------|------------|--------------|
| **触发后行为** | 可能重试、降级、继续 | **立即中断，不继续** |
| **设计目的** | 处理技术异常 | 安全防护，阻止风险 |
| **返回类型** | `error` | `GuardrailTripwireTriggeredError` |

#### 代码对比

```go
// 工具调用失败 → 反馈给 LLM，让它决定下一步
if err != nil {
    errorOutput := "Tool execution failed: " + err.Error()
    continue  // 继续循环
}

// Tripwire 触发 → 直接返回错误，中断执行
if grResult.Output.TripwireTriggered {
    return nil, &GuardrailTripwireTriggeredError{...}  // 立即终止
}
```

#### 实际安全场景

| 风险类型 | 护栏位置 | 具体实现 |
|---------|---------|---------|
| 敏感信息泄露 | 输入护栏 | 检测"密码"、"信用卡"、"身份证" |
| 禁止话题 | 输入护栏 | 检测"非法"、"暴力"、"犯罪" |
| 输出长度攻击 | 输出护栏 | 限制输出长度 |
| PII 信息泄露 | 输出护栏 | 检测个人身份信息 |

#### 拦截流程示例

```
用户输入: "帮我查一下密码是多少"
    ↓
输入护栏检查 → 发现"密码"关键词
    ↓
TripwireTriggered = true
    ↓
返回 GuardrailTripwireTriggeredError
    ↓
Agent 执行中断，不调用 LLM
```

#### 面试标准表述

> "Tripwire 熔断是一种安全防护机制，和普通错误处理的核心区别是：**普通错误可能重试或降级，而 Tripwire 触发后立即中断执行**。
>
> 设计意图是：当检测到安全风险时，不能让请求继续处理，必须立刻阻止。框架支持输入/输出双向护栏：输入护栏在 LLM 调用前检查，防止恶意请求消耗资源；输出护栏在返回前检查，防止 LLM 生成不当内容。"

---

## 问题 8：改进方向与技术挑战

### 面试官问题

> 回顾你的整个项目，如果让你重新设计，有哪些地方你会改进？在开发过程中遇到的最大技术挑战是什么？

### 标准答案

#### 改进方向

**1. 推理模式升级：ReAct → Reflexion**

| 特性 | ReAct | Reflexion |
|-----|-------|-----------|
| 模式 | Thought → Action → Observation | ReAct + 自我反思 + 记忆 |
| 失败处理 | 没有显式反思 | 失败后反思原因，记录经验 |
| 学习能力 | 单次执行 | 跨任务积累经验 |

> "ReAct 只有 Thought-Action-Observation 循环，当任务失败时缺乏系统性的反思机制。Reflexion 增加了自我反思和经验记忆，Agent 可以从错误中学习。"

**2. 工具路由优化：关键词匹配 → Embedding 语义匹配**

目前的 KeywordRouter 是静态关键词匹配，可以改成基于 Embedding 的语义匹配，准确率会更高。

#### 技术挑战（架构层面）

| 挑战 | 具体问题 | 解决方案 |
|-----|---------|---------|
| **循环依赖** | Agent 和 Tool 互相引用 | 抽象出 `types.AgentLike` 最小接口 |
| **API 兼容性** | Responses API vs Chat Completions | 在 Runner 中用条件判断分流 |
| **工具调用结果统一** | MCP 返回格式多样 | 统一转成 JSON 字符串 |
| **并发安全** | Session 多协程访问 | 使用 `sync.Mutex` 保护 |

#### 面试标准表述

> "如果重新设计，我会做两个改进：
>
> **第一，推理模式升级**：把 ReAct 改成 Reflexion，增加自我反思和经验记忆机制。
>
> **第二，工具路由优化**：目前的 KeywordRouter 是静态关键词匹配，可以改成基于 Embedding 的语义匹配。
>
> **最大的技术挑战**是解决模块间的循环依赖。我的解决方案是抽象出一个最小化的 `AgentLike` 接口，只暴露必要的方法，让 Tool 依赖接口而不是具体实现。"

---

## 面试总结

### 考察点与表现评估

| 问题领域 | 考察点 | 关键知识点 |
|---------|-------|-----------|
| 架构设计 | 五层分层、职责划分 | 每层职责、为什么用分层 |
| 核心引擎 | Runner 循环、错误处理 | 工具失败反馈 LLM、护栏直接中断 |
| 工具系统 | ToolRouter、MCP 集成 | 关键词匹配（不是 LLM）、Schema 转换 |
| 多模型支持 | API 兼容性 | 两种 API 区别、OpenRouter 聚合 |
| 安全护栏 | Tripwire 熔断 | 和普通错误的区别、立即中断 |
| 多智能体 | Agent-as-Tool | 协作机制、最小权限原则 |

### 面试建议

1. **重新熟悉核心代码**：
   - `pkg/runner/runner.go` — 执行引擎
   - `pkg/agent/guardrail.go` — 护栏系统
   - `pkg/tool/mcp.go` — MCP 集成
   - `pkg/tool/router.go` — 工具路由

2. **准备具体数字**：比如"减少 30-50% Token"要能解释计算方式

3. **理解设计决策**：每个接口、每个分层都要能说出"为什么这样设计"

4. **准备失败案例**：面试官喜欢问"遇到什么困难，怎么解决的"

### 核心文件路径

```
nvgo-main/
├── pkg/
│   ├── runner/runner.go      ← 执行引擎核心（必读）
│   ├── agent/
│   │   ├── agent.go          ← Agent 配置
│   │   ├── guardrail.go      ← 护栏系统（必读）
│   │   └── instruction.go    ← 指令管理
│   ├── tool/
│   │   ├── tool.go           ← 工具接口
│   │   ├── router.go         ← 工具路由（必读）
│   │   └── mcp.go            ← MCP 集成（必读）
│   ├── memory/
│   │   └── sqlite.go         ← Session 实现
│   └── pattern/
│       ├── agent_tool.go     ← Multi-Agent
│       └── react.go          ← ReAct 模式
└── examples/                  ← 10 个完整示例
```

---

> 文档生成时间：2026-01-11
>
> 建议复习频率：面试前 1-2 天重点复习
