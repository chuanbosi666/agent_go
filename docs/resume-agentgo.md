# agentgo 项目简历内容

> 适用岗位：大模型应用开发 / AI Agent 开发工程师

---

## 一、推荐版本（详细描述）

### **agentgo - Go 语言 AI Agent 框架**

**项目类型**：个人开源项目
**技术栈**：Go、OpenAI API、MCP 协议、SQLite、JSON Schema
**项目规模**：核心代码 2.8K 行，13 个接口，10 个完整示例

---

#### 项目简介

参考 OpenAI Agents SDK 设计理念，独立实现的生产级 AI Agent 框架。采用 5 层架构（执行层 → 能力层 → 基础设施层），支持工具调用、多智能体协作、会话记忆、安全护栏等核心能力，兼容 OpenAI / Claude / Gemini / DeepSeek / Qwen / 本地 Ollama 等 100+ 模型。

**核心特性**：
- 完整的 Agent Loop 执行引擎
- FunctionTool + MCP 协议双模式工具系统
- Agent-as-Tool 多智能体协作
- SQLite 会话持久化
- 输入/输出双向安全护栏
- 动态指令与 ReAct 推理模式

---

#### 工作内容

**1. 架构设计与核心引擎开发**
- 设计 5 层分层架构（执行层/能力层/基础设施层），定义 13 个核心接口，确保模块解耦与可扩展性
- 实现 Runner 执行引擎（620 行），完成 Agent Loop 主循环：输入 → 护栏校验 → LLM 推理 → 工具调用 → 输出
- 支持 MaxTurns 控制与异常熔断，工具执行失败自动反馈给 LLM 进行重试或降级

**2. 工具系统设计与 MCP 协议集成**
- 设计 FunctionTool 统一接口，支持 JSON Schema 参数定义、严格模式校验、条件启用（IsEnabled）
- 实现 MCP 协议适配层（380 行），自动转换外部工具 Schema，支持工具过滤（MCPToolFilter）
- 开发 ToolRouter 工具路由器，通过关键词匹配动态筛选相关工具，减少 30-50% Token 消耗

**3. 多智能体协作机制**
- 开发 Agent-as-Tool 包装器，支持将任意 Agent 封装为可调用工具
- 实现 Manager-Specialist 分层协作架构，支持任务分解与专家调度
- 编写 dev-team 示例（5 个 Agent 协作完成软件开发全流程）

**4. 动态指令与推理模式**
- 设计 Instructions 接口三种实现：静态字符串、动态函数、基于状态的模板渲染
- 实现 ReAct 推理模式（Thought → Action → Observation 循环）
- 支持 StateProvider 注入运行时状态到提示词

**5. 会话管理与安全护栏**
- 基于 SQLite 实现会话持久化（338 行），支持多轮对话历史的存储、检索与自动序列化
- 实现输入/输出双向护栏机制，支持内容过滤、敏感词拦截、Tripwire 熔断
- 设计 GuardrailTripwireTriggeredError 错误类型，支持安全审计与告警

---

## 二、精简版（一页简历适用）

### **agentgo - AI Agent 框架 (Go)**

#### 项目简介

参考 OpenAI Agents SDK 设计，独立实现的生产级 AI Agent 框架，支持工具调用、多智能体协作、会话记忆、安全护栏，兼容 100+ 大模型。

#### 工作内容

- **架构设计**：5 层分层架构，13 个核心接口，接口驱动、Builder 模式
- **执行引擎**：实现 Agent Loop（输入 → 护栏 → LLM → 工具调用 → 输出），支持异常熔断
- **工具系统**：FunctionTool + MCP 协议双模式，ToolRouter 动态路由节省 30-50% Token
- **多智能体**：Agent-as-Tool 封装，Manager-Specialist 分层协作架构
- **会话管理**：SQLite 持久化，多轮对话历史存储与检索
- **安全护栏**：输入/输出双向护栏，Tripwire 熔断机制

---

## 三、一句话版本

基于 Go 语言独立实现的 AI Agent 框架，支持 Tool Calling、Multi-Agent 协作、Memory 管理、Guardrails 安全护栏，兼容 100+ 大模型。

---

## 四、岗位关键词匹配表

| 岗位常见要求 | 项目对应点 |
|-------------|-----------|
| LLM API 调用经验 | OpenAI Chat Completions / Responses API 集成 |
| Agent 开发经验 | 完整 Agent Loop、ReAct 模式、动态 Prompt |
| Tool Calling / Function Calling | FunctionTool 接口、JSON Schema、MCP 协议 |
| 多智能体 / Multi-Agent | Agent-as-Tool、Manager-Specialist 架构 |
| Prompt Engineering | 动态 InstructionsFunc、ReAct 模板 |
| RAG / 记忆管理 | SQLite Session、多轮对话持久化 |
| 安全与合规 | InputGuardrail / OutputGuardrail 护栏系统 |
| Token 优化 | ToolRouter 工具路由、动态筛选 |
| 框架 / SDK 开发 | Builder 模式、接口驱动、可扩展设计 |

---

## 五、ATS 关键词（简历系统扫描用）

```
Go, Golang, AI Agent, LLM, Large Language Model, OpenAI API,
Tool Calling, Function Calling, Multi-Agent, MCP Protocol,
Prompt Engineering, ReAct, Memory Management, SQLite,
Guardrails, JSON Schema, Builder Pattern, Interface Design
```

---

## 六、面试问答准备

### Q1: 介绍一下你的 AI Agent 项目

> 我独立开发了一个 Go 语言的 AI Agent 框架，叫 agentgo。设计上参考了 OpenAI 官方的 Agents SDK，但用 Go 重新实现，更适合高并发的后端场景。
>
> 核心是一个完整的 Agent 执行循环：用户输入先过输入护栏做安全校验，然后调用 LLM 推理，如果模型返回工具调用请求，就执行对应工具并把结果反馈给模型继续推理，直到输出最终答案或达到最大轮次。
>
> 工具系统我设计了统一的 FunctionTool 接口，支持 JSON Schema 定义参数，也集成了 MCP 协议，可以接入外部工具服务器。

---

### Q2: 多智能体协作是怎么实现的？

> 我实现了一个 Agent-as-Tool 的包装器。核心思路是把任何 Agent 封装成一个 FunctionTool，这样主 Agent 调用它就像调用普通工具一样。
>
> 实际应用中可以搭建 Manager-Specialist 架构。比如我写了一个 dev-team 示例，有一个项目经理 Agent 负责任务分解和协调，下面挂了需求分析、架构设计、编码实现、测试验证四个专家 Agent。用户提一个需求，经理 Agent 会依次调用各个专家完成整个开发流程。
>
> 这种设计的好处是每个 Agent 职责单一、Prompt 精简，比一个大而全的 Agent 效果更好。

---

### Q3: Token 优化做了哪些工作？

> 主要是工具路由优化。当 Agent 可用工具很多时，每次都把所有工具描述传给模型会消耗大量 Token。
>
> 我实现了一个 KeywordRouter，给每个工具关联一组关键词，根据用户输入动态匹配最相关的 Top-N 个工具传给模型。实测可以减少 30-50% 的工具描述 Token。
>
> 另外动态指令也有帮助，不是每次都传完整的系统提示词，而是根据当前对话阶段生成精简版本。

---

### Q4: Guardrails 护栏系统的设计？

> 分输入护栏和输出护栏两层。
>
> 输入护栏在用户消息进入 LLM 之前执行，可以做敏感词检测、注入攻击拦截、格式校验这些。如果触发 Tripwire 条件就直接熔断，不再调用模型。
>
> 输出护栏在模型返回后执行，可以过滤不当内容、脱敏处理、格式规范化。
>
> 两层护栏都是接口设计，用户可以自定义实现，比如接入内容审核 API 或者规则引擎。

---

### Q5: 为什么选择 Go 而不是 Python？

> 几个原因：
>
> 1. **性能**：Go 的并发模型（goroutine）非常适合处理多个 Agent 并行执行、多工具并发调用的场景
> 2. **部署**：Go 编译成单一二进制文件，部署简单，不需要依赖管理
> 3. **类型安全**：静态类型在框架开发中能提前发现很多问题，接口设计也更清晰
> 4. **生态空白**：Python 有 LangChain、LlamaIndex，但 Go 生态缺少成熟的 Agent 框架

---

### Q6: MCP 协议是什么？为什么要集成？

> MCP（Model Context Protocol）是 Anthropic 提出的一个开放协议，用于标准化 LLM 与外部工具/数据源的交互。
>
> 集成 MCP 的好处是可以接入大量现成的工具服务器，比如文件系统、数据库、Web 服务等，不需要为每个工具单独写适配代码。
>
> 我在框架里实现了一个适配层，自动把 MCP 的工具 Schema 转换成 OpenAI Function Calling 的格式，对上层 Agent 透明。

---

### Q7: 项目中遇到的最大挑战是什么？

> Schema 转换是比较麻烦的一块。MCP 和 OpenAI 的 JSON Schema 规范有些细节差异，比如类型表示、必填字段处理、嵌套对象的处理方式。
>
> 我写了一个 strictschema 模块专门处理这些转换，确保各种边界情况都能正确处理。过程中读了很多 OpenAI 和 MCP 的文档，也参考了 Python SDK 的实现。

---

### Q8: 这个项目有实际应用场景吗？

> 有的。比如 dev-team 示例就是一个完整的应用：输入一个需求，系统自动完成需求分析、架构设计、代码生成、测试验证。
>
> 实际业务中可以用来做：
> - 智能客服（多轮对话 + 工具查询）
> - 自动化运维（Agent 调用监控、部署工具）
> - 代码助手（Agent 调用代码搜索、执行测试）
> - 数据分析（Agent 调用 SQL 查询、图表生成）

---

## 七、技术架构图（面试白板用）

```
┌─────────────────────────────────────────────────────────────┐
│                      User Input                              │
└─────────────────────────┬───────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                  Input Guardrails                            │
│            (敏感词检测 / 注入拦截 / 格式校验)                  │
└─────────────────────────┬───────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                     Agent Loop                               │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  1. Build Messages (System + History + User)        │    │
│  │  2. Select Tools (ToolRouter)                       │    │
│  │  3. Call LLM (OpenAI/Claude/Gemini/...)            │    │
│  │  4. Parse Response                                  │    │
│  │     ├─ Text → Output Guardrails → Return           │    │
│  │     └─ Tool Call → Execute → Feedback → Loop       │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────┬───────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        ▼                 ▼                 ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ FunctionTool │  │   MCP Tool   │  │ Agent Tool   │
│  (本地函数)   │  │  (外部协议)   │  │ (子 Agent)   │
└──────────────┘  └──────────────┘  └──────────────┘
```

---

## 八、项目亮点总结（30秒版）

1. **完整性**：从 Agent Loop 到工具调用到记忆管理，覆盖 AI Agent 全链路
2. **工程化**：接口驱动、Builder 模式、零值可用，符合 Go 最佳实践
3. **创新点**：Agent-as-Tool 多智能体协作、ToolRouter Token 优化
4. **实用性**：10 个示例覆盖真实场景，可直接用于生产环境
5. **前沿性**：集成 MCP 协议，支持 100+ 模型，紧跟行业发展

---

## 九、深度技术面试题（考察项目理解）

> 以下问题用于考察对项目架构、设计决策、实现细节的深度理解。分为 6 个维度，每个维度包含基础题和进阶题。

---

### 第一部分：架构设计（基础理解）

#### Q1.1 请描述 agentgo 项目的整体架构分层，各层之间的职责边界是什么？

**参考答案**：

项目采用 5 层架构，单向依赖：

```
应用层 (examples/)
    ↓
执行层 (pkg/runner/) —— 核心调度，编排整个 Agent 运行流程
    ↓
能力层 (pkg/agent/, pkg/tool/, pkg/pattern/) —— Agent 配置、工具管理、协作模式
    ↓
基础设施层 (pkg/memory/, pkg/types/, pkg/config/) —— 会话存储、共享类型、配置管理
    ↓
外部依赖层 (openai-go, mcp-sdk, sqlite3)
```

**各层职责边界**：
- `runner`：编排执行流程，不包含具体业务逻辑
- `agent`：定义配置（指令、模型、护栏），不包含执行逻辑
- `tool`：工具接口和实现，不依赖 Agent 具体实现
- `memory`：存储抽象，不感知 Agent/Runner 存在
- `types`：最底层共享类型，无业务逻辑

**表达技巧**：先说"分 5 层"给总览，再逐层点明核心职责，最后强调"单向依赖"原则。

---

#### Q1.2 项目中定义了 13 个接口。为什么采用接口驱动设计？这种设计的好处和潜在问题？

**参考答案**：

**好处**：
1. **可扩展性**：Session 可以是 SQLite、Redis、内存；ToolRouter 可以是关键词、语义、规则
2. **可测试性**：可以 mock 任何接口进行单元测试
3. **解耦**：各模块通过接口通信，实现可替换
4. **符合 Go 惯例**：Go 推崇小接口、隐式实现

**潜在问题**：
1. **学习成本**：新用户需要理解多个接口的关系
2. **过度抽象风险**：如果接口设计不当，可能导致不必要的复杂性
3. **性能开销**：接口调用有轻微的运行时开销（实际可忽略）

---

#### Q1.3（进阶）`pkg/types` 包中定义了 `AgentLike` 这个"最小接口"。为什么要单独抽出这个包？

**参考答案**：

**解决循环依赖问题**。

具体场景：
- `pkg/tool` 需要知道当前 Agent 的名称/模型（用于 MCP 工具过滤）
- 如果 `tool` 直接导入 `agent` 包，而 `agent` 又需要导入 `tool`（定义工具列表），就形成循环依赖

解决方案：
- 抽取最小接口 `AgentLike`（只有 `GetName()` 和 `GetModel()`）到独立的 `types` 包
- `tool` 依赖 `types.AgentLike`，而不是完整的 `agent.Agent`
- 打破了 agent ↔ tool 的循环引用

```go
// pkg/types/agent.go
type AgentLike interface {
    GetName() string
    GetModel() string
}
```

**这是 Go 语言处理循环依赖的标准模式**。

---

### 第二部分：核心执行流程

#### Q2.1 请完整描述 `Runner.run()` 方法的执行流程

**参考答案**：

```
1. 【输入处理】接收 Input（string 或 InputItems）
       ↓
2. 【输入护栏】执行所有 InputGuardrails
   - 任一返回 TripwireTriggered=true → 立即终止，返回错误
       ↓
3. 【主循环】最多执行 MaxTurns 轮（默认 10）
   │
   ├─ 3.1 获取动态指令（Instructions.GetInstructions）
   ├─ 3.2 收集工具（FunctionTools + MCP Tools）
   ├─ 3.3 工具路由（当工具数 > threshold 时过滤）
   ├─ 3.4 选择 API：
   │      - 有 Prompt 配置 → Responses API
   │      - 否则 → Chat Completions API
   ├─ 3.5 调用 LLM
   ├─ 3.6 解析响应：
   │      - 纯文本 → 跳出循环
   │      - 工具调用 → 执行工具 → 结果反馈 → 继续循环
   │      - Handoff → 切换 Agent → 继续循环
   └─ 3.7 检查终止条件
       ↓
4. 【输出护栏】执行所有 OutputGuardrails
       ↓
5. 【返回结果】RunResult {FinalOutput, NewItems, RawResponses, ...}
```

**关键决策点**：
- 工具执行失败 → 错误信息反馈给 LLM，LLM 可以重试或换方案
- MaxTurns 超限 → 返回 `MaxTurnsExceededError`

---

#### Q2.2 Responses API 和 Chat Completions API 的区别？Runner 如何选择？

**参考答案**：

| 维度 | Responses API | Chat Completions API |
|------|---------------|---------------------|
| 提供方 | OpenAI 独有 | OpenAI 兼容标准（Claude/Gemini/本地模型都支持） |
| 特性 | 支持 Prompt 模板、版本管理、变量替换 | 通用的消息列表格式 |
| 工具格式 | 更丰富的元数据 | 标准 Function Calling |
| 适用场景 | OpenAI 高级功能 | 跨模型兼容 |

**选择逻辑**：

```go
if agent.Prompt != nil {
    // 使用 Responses API（OpenAI only）
    return callResponsesAPI(...)
} else {
    // 使用 Chat Completions API（兼容所有模型）
    return callChatCompletionsAPI(...)
}
```

**简单记忆**：有 Prompt 配置用 Responses，否则用 Chat Completions。

---

#### Q2.3（进阶）KeywordRouter 的实现原理和性能考量

**参考答案**：

**实现原理**：

```go
type KeywordRouter struct {
    ToolKeywords map[string][]string  // 工具名 -> 关键词列表
    TopN         int                  // 返回前 N 个工具
}

// 路由算法：
func RouteTools(input, tools):
    1. 从用户输入提取关键词（分词、转小写）
    2. 遍历每个工具，计算匹配分数：
       score = len(输入关键词 ∩ 工具关键词)
    3. 按分数排序
    4. 返回前 TopN 个工具
```

**性能考量**：

| 考量点 | 设计决策 |
|--------|---------|
| 时间复杂度 | O(n×m)，n=工具数，m=关键词数；对于常规场景（<100 工具）完全够用 |
| 内存 | 关键词预加载到内存，查询时零分配 |
| Token 节省 | 从 20 个工具筛选到 5 个，节省约 75% 工具描述 Token |
| 精度 vs 召回 | 关键词匹配是粗筛，宁可多选不漏选（TopN 可调） |

**局限性**：关键词匹配无法理解语义。如果需要更精准的路由，可实现 `SemanticRouter`（基于 embedding 相似度）。

---

### 第三部分：指令与工具系统

#### Q3.1 Instructions 接口的三种实现及使用场景

**参考答案**：

| 实现 | 定义 | 使用场景 | 示例 |
|------|------|---------|------|
| `InstructionsStr` | `type InstructionsStr string` | 静态提示词，不变化 | `"你是一个翻译助手"` |
| `InstructionsFunc` | `func(ctx, agent) (string, error)` | 动态生成，依赖运行时状态 | 根据用户偏好语言生成不同提示 |
| `DynamicInstruction` | 结构体，含 BasePrompt + StateProvider + Template | 基于状态的模板渲染 | ReAct 模式中注入当前步骤和观察 |

**代码示例**：

```go
// 静态
agent.WithInstructions("你是客服助手")

// 动态函数
agent.WithInstructionsFunc(func(ctx context.Context, a *Agent) (string, error) {
    user := getUserFromContext(ctx)
    return fmt.Sprintf("你是 %s 的专属助手", user.Name), nil
})

// 动态模板
agent.WithDynamicInstruction(&DynamicInstruction{
    BasePrompt:    "你是 ReAct 推理助手",
    StateProvider: reactState,
    Template:      "当前步骤：{{step}}\n观察结果：{{observation}}",
})
```

---

#### Q3.2 FunctionTool 中 FailureErrorFunction 和 IsEnabled 的作用

**参考答案**：

**`FailureErrorFunction`**：工具执行失败时的自定义错误处理

```go
type ToolErrorFunction func(ctx context.Context, err error) string

// 使用场景：
// - 将技术错误转换为用户友好的提示
// - 添加重试建议
// - 记录错误日志

tool := FunctionTool{
    Name: "query_database",
    FailureErrorFunction: func(ctx context.Context, err error) string {
        log.Error("DB error", err)
        return "数据库暂时不可用，请稍后重试"
    },
}
```

**`IsEnabled`**：动态控制工具是否可用

```go
type FunctionToolEnabler func(ctx context.Context, agent AgentLike) bool

// 使用场景：
// - 根据用户权限启用/禁用工具
// - 根据时间段控制（如交易时间内才启用下单工具）
// - 根据上下文状态控制

tool := FunctionTool{
    Name: "admin_operation",
    IsEnabled: func(ctx context.Context, agent AgentLike) bool {
        user := getUserFromContext(ctx)
        return user.IsAdmin
    },
}
```

---

#### Q3.3（进阶）MCP 工具和本地 FunctionTool 如何统一管理？MCPToolFilter 的设计目的？

**参考答案**：

**统一管理机制**：

```go
// Runner 收集工具时的流程：
func collectTools(agent *Agent) []Tool {
    var allTools []Tool

    // 1. 收集本地 FunctionTools
    for _, ft := range agent.Tools {
        if ft.IsEnabled == nil || ft.IsEnabled(ctx, agent) {
            allTools = append(allTools, ft)
        }
    }

    // 2. 收集 MCP 工具（转换为统一接口）
    for _, server := range agent.MCPServers {
        mcpTools, _ := server.ListTools(ctx, agent)
        for _, mt := range mcpTools {
            // MCP Tool 包装成 Tool 接口
            allTools = append(allTools, wrapMCPTool(server, mt))
        }
    }

    return allTools
}
```

**MCPToolFilter 设计目的**：

```go
type MCPToolFilter interface {
    FilterMCPTool(ctx context.Context, filterCtx MCPToolFilterContext, tool *mcp.Tool) (bool, error)
}
```

使用场景：
1. **权限控制**：某些 MCP 工具只对特定 Agent 开放
2. **能力匹配**：根据 Agent 模型能力筛选工具（小模型不给复杂工具）
3. **成本控制**：过滤掉高成本的 MCP 工具
4. **动态启用**：根据运行时状态决定是否暴露某工具

```go
// 示例：只有 gpt-4 模型才能使用代码执行工具
type ModelBasedFilter struct{}

func (f *ModelBasedFilter) FilterMCPTool(ctx context.Context, fCtx MCPToolFilterContext, tool *mcp.Tool) (bool, error) {
    if tool.Name == "execute_code" {
        return strings.Contains(fCtx.Agent.GetModel(), "gpt-4"), nil
    }
    return true, nil
}
```

---

### 第四部分：设计模式与多智能体

#### Q4.1 Agent-as-Tool 模式的实现原理

**参考答案**：

**核心思路**：将任意 Agent 封装成 FunctionTool，主 Agent 调用它就像调用普通工具。

```go
func WrapAgentAsTool(subAgent *Agent, maxTurns int) FunctionTool {
    return FunctionTool{
        Name:        subAgent.Name,
        Description: fmt.Sprintf("调用 %s 专家处理任务", subAgent.Name),
        ParamsJSONSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "task": map[string]any{
                    "type":        "string",
                    "description": "要处理的任务描述",
                },
            },
            "required": []string{"task"},
        },
        OnInvokeTool: func(ctx context.Context, args string) (any, error) {
            var params struct{ Task string }
            json.Unmarshal([]byte(args), &params)

            // 创建新的 Runner 执行子 Agent
            runner := &Runner{Config: RunConfig{MaxTurns: maxTurns}}
            result, err := runner.Run(ctx, subAgent, params.Task)
            if err != nil {
                return nil, err
            }
            return result.FinalOutput, nil
        },
    }
}
```

**使用示例**：

```go
// 定义专家 Agent
mathExpert := agent.New("数学专家").WithInstructions("你擅长数学计算...")
codeExpert := agent.New("代码专家").WithInstructions("你擅长编程...")

// 包装成工具
mathTool := WrapAgentAsTool(mathExpert, 5)
codeTool := WrapAgentAsTool(codeExpert, 5)

// 主 Agent 使用这些工具
manager := agent.New("项目经理").
    WithInstructions("你负责协调任务，根据需要调用专家").
    WithTools([]FunctionTool{mathTool, codeTool})
```

---

#### Q4.2 ReAct 模式的实现

**参考答案**：

**ReAct = Reasoning + Acting**，让 LLM 显式输出思考过程。

**组件分工**：

| 组件 | 职责 |
|------|------|
| `DefaultReActInstruction` | 预定义的 ReAct 提示词模板，规定输出格式 |
| `ReActStateProvider` | 追踪执行步骤和观察结果，注入到提示词中 |

**提示词模板**：

```
你是一个 ReAct 推理助手。按以下格式思考和行动：

Thought: [分析当前情况，决定下一步]
Action: [选择要使用的工具]
Action Input: [工具的输入参数，JSON 格式]
Observation: [工具返回的结果]
... (重复直到得出答案)
Thought: 我已经得到了最终答案
Final Answer: [最终答案]
```

**StateProvider 实现**：

```go
type ReActStateProvider struct {
    Steps       []ReActStep
    CurrentStep int
}

func (p *ReActStateProvider) GetState(ctx context.Context) (map[string]string, error) {
    return map[string]string{
        "step":        fmt.Sprintf("%d", p.CurrentStep),
        "observation": p.LastObservation(),
        "history":     p.FormatHistory(),
    }, nil
}
```

---

#### Q4.3（进阶）设计"开发团队"多 Agent 协作架构

**参考答案**：

```
                    ┌─────────────────┐
                    │   User Input    │
                    └────────┬────────┘
                             ▼
                    ┌─────────────────┐
                    │ Project Manager │ ← 任务分解、进度协调、结果汇总
                    │    (主 Agent)   │
                    └────────┬────────┘
                             │
         ┌───────────┬───────┴───────┬───────────┐
         ▼           ▼               ▼           ▼
    ┌─────────┐ ┌─────────┐   ┌─────────┐ ┌─────────┐
    │   PM    │ │  Arch   │   │  Dev    │ │  Test   │
    │ 需求分析 │ │ 架构设计 │   │ 代码实现 │ │ 测试验证 │
    └─────────┘ └─────────┘   └─────────┘ └─────────┘
```

**实现代码结构**：

```go
// 1. 定义专家 Agent
reqAgent := agent.New("需求分析师").
    WithInstructions("你负责分析用户需求，输出结构化的需求文档...")

archAgent := agent.New("架构师").
    WithInstructions("你负责技术方案设计，输出架构图和接口定义...")

devAgent := agent.New("开发者").
    WithInstructions("你负责代码实现，遵循架构设计...").
    WithTools([]FunctionTool{codeWriteTool, codeSearchTool})

testAgent := agent.New("测试工程师").
    WithInstructions("你负责编写测试用例和执行测试...").
    WithTools([]FunctionTool{testRunTool})

// 2. 包装成工具
reqTool := WrapAgentAsTool(reqAgent, 5)
archTool := WrapAgentAsTool(archAgent, 5)
devTool := WrapAgentAsTool(devAgent, 10)
testTool := WrapAgentAsTool(testAgent, 5)

// 3. 主 Agent 协调
manager := agent.New("项目经理").
    WithInstructions(`
        你是项目经理，负责协调软件开发流程。
        收到需求后，依次执行：
        1. 调用需求分析师分析需求
        2. 调用架构师设计方案
        3. 调用开发者实现代码
        4. 调用测试工程师验证
        最后汇总报告给用户。
    `).
    WithTools([]FunctionTool{reqTool, archTool, devTool, testTool})
```

**信息流转设计**：
- 每个专家的输出作为下一个专家的输入
- 项目经理维护上下文，传递关键信息
- 支持回溯（测试失败 → 通知开发者修复）

---

### 第五部分：会话管理与状态

#### Q5.1 SQLiteSession 的存储设计

**参考答案**：

**数据库 Schema**：

```sql
CREATE TABLE agent_sessions (
    id TEXT PRIMARY KEY,           -- 会话 ID（UUID）
    created_at TIMESTAMP,          -- 创建时间
    updated_at TIMESTAMP           -- 最后更新时间
);

CREATE TABLE agent_messages (
    id TEXT PRIMARY KEY,           -- 消息 ID（UUID）
    session_id TEXT,               -- 关联的会话 ID
    item_json TEXT,                -- JSON 序列化的 ResponseItem
    created_at TIMESTAMP,          -- 创建时间
    FOREIGN KEY (session_id) REFERENCES agent_sessions(id)
);
```

**存储流程**：

```go
// 添加消息
func (s *SQLiteSession) AddItems(ctx context.Context, items []ResponseItem) error {
    for _, item := range items {
        jsonBytes, _ := json.Marshal(item)  // 序列化为 JSON
        _, err := s.db.Exec(
            "INSERT INTO agent_messages (id, session_id, item_json, created_at) VALUES (?, ?, ?, ?)",
            uuid.New().String(), s.sessionID, string(jsonBytes), time.Now(),
        )
    }
    return nil
}

// 获取历史
func (s *SQLiteSession) GetItems(ctx context.Context, limit int) ([]ResponseItem, error) {
    rows, _ := s.db.Query(
        "SELECT item_json FROM agent_messages WHERE session_id = ? ORDER BY created_at DESC LIMIT ?",
        s.sessionID, limit,
    )
    // 反序列化并返回
}
```

---

#### Q5.2 为什么用 Mutex 而不是 RWMutex？

**参考答案**：

**当前实现用 `sync.Mutex`**：

```go
type SQLiteSession struct {
    db        *sql.DB
    sessionID string
    mu        sync.Mutex  // 不是 RWMutex
}
```

**原因**：

1. **SQLite 写锁限制**：SQLite 本身在写操作时会锁定整个数据库，读写分离收益有限
2. **操作模式**：Agent 对话通常是"读历史 → 生成回复 → 写入新消息"的串行模式，读写交替，RWMutex 的并发读优势不明显
3. **简单性**：Mutex 实现更简单，不易出错

**何时可以改用 RWMutex**：

- 使用其他数据库（如 PostgreSQL、Redis）时
- 读操作远多于写操作（如只读分析场景）
- 需要支持多个 goroutine 并发读取会话历史

```go
// 改用 RWMutex 的示例
func (s *Session) GetItems(...) {
    s.mu.RLock()         // 读锁
    defer s.mu.RUnlock()
    // ...
}

func (s *Session) AddItems(...) {
    s.mu.Lock()          // 写锁
    defer s.mu.Unlock()
    // ...
}
```

---

#### Q5.3（进阶）实现分布式会话存储需要修改什么？

**参考答案**：

**Session 接口是否足够**：

```go
type Session interface {
    GetItems(ctx context.Context, limit int) ([]ResponseItem, error)
    AddItems(ctx context.Context, items []ResponseItem) error
    PopItem(ctx context.Context) (*ResponseItem, error)
    ClearSession(ctx context.Context) error
}
```

✅ **基本足够**，接口抽象得当。

**实现 RedisSession 需要的修改**：

```go
type RedisSession struct {
    client    *redis.Client
    sessionID string
}

func (s *RedisSession) GetItems(ctx context.Context, limit int) ([]ResponseItem, error) {
    // 使用 Redis List：LRANGE session:{id}:messages -limit -1
    vals, err := s.client.LRange(ctx, s.key(), -int64(limit), -1).Result()
    // 反序列化返回
}

func (s *RedisSession) AddItems(ctx context.Context, items []ResponseItem) error {
    // 使用 Redis List：RPUSH session:{id}:messages item1 item2 ...
    for _, item := range items {
        jsonBytes, _ := json.Marshal(item)
        s.client.RPush(ctx, s.key(), jsonBytes)
    }
    // 可选：设置过期时间
    s.client.Expire(ctx, s.key(), 24*time.Hour)
    return nil
}
```

**可能需要扩展的点**：

| 需求 | 当前接口 | 是否需要扩展 |
|------|---------|-------------|
| TTL 过期 | 无 | 可在实现层处理，不必改接口 |
| 分布式锁 | 无 | 实现层使用 Redis SETNX |
| 会话列表 | 无 | 可新增 `ListSessions()` 方法 |
| 批量操作 | 有 AddItems | ✅ 已支持 |

---

### 第六部分：护栏与错误处理

#### Q6.1 输入/输出护栏的执行时机和 TripwireTriggered 的作用

**参考答案**：

**执行时机**：

```
用户输入
    ↓
┌─────────────────────────┐
│    Input Guardrails     │ ← 在调用 LLM 之前
│  - 敏感词检测           │
│  - 注入攻击拦截         │
│  - 格式校验             │
└───────────┬─────────────┘
            ↓
        LLM 推理
            ↓
┌─────────────────────────┐
│   Output Guardrails     │ ← 在返回用户之前
│  - 内容过滤             │
│  - 脱敏处理             │
│  - 格式规范化           │
└───────────┬─────────────┘
            ↓
        返回用户
```

**TripwireTriggered = true 时**：

```go
type GuardrailFunctionOutput struct {
    TripwireTriggered bool   // true = 触发熔断
    OutputInfo        any    // 检查结果信息
}
```

- **立即终止执行**：不再调用 LLM / 不再返回给用户
- **返回错误**：`GuardrailTripwireTriggeredError`
- **记录日志**：`OutputInfo` 中包含触发原因

```go
// Runner 中的处理逻辑
for _, guard := range inputGuardrails {
    result, _ := guard.GuardrailFunc(ctx, agent, input)
    if result.TripwireTriggered {
        return nil, &GuardrailTripwireTriggeredError{
            GuardrailName: guard.Name,
            OutputInfo:    result.OutputInfo,
            IsInput:       true,
        }
    }
}
```

---

#### Q6.2 工具执行失败时的处理机制

**参考答案**：

**处理流程**：

```go
func (r *Runner) executeTool(ctx context.Context, tool Tool, args string) (any, error) {
    result, err := tool.Invoke(ctx, args)

    if err != nil {
        // 1. 如果有自定义错误处理函数
        if ft, ok := tool.(FunctionTool); ok && ft.FailureErrorFunction != nil {
            errorMsg := ft.FailureErrorFunction(ctx, err)
            return errorMsg, nil  // 返回友好错误信息给 LLM
        }

        // 2. 默认：返回原始错误信息
        return fmt.Sprintf("工具执行失败: %v", err), nil
    }

    return result, nil
}
```

**LLM 能否恢复？能！**

工具失败不会终止整个执行，而是：
1. 将错误信息作为工具结果反馈给 LLM
2. LLM 看到错误后可以：
   - 重试（换参数）
   - 换工具
   - 直接告诉用户无法完成

**示例对话**：

```
User: 查询北京天气
LLM: [调用 weather_api，参数 city="北京"]
Tool: 返回错误 "API 超时"
LLM: 抱歉，天气服务暂时不可用。您可以稍后再试，或者我可以帮您查询其他信息。
```

---

#### Q6.3（进阶）生产环境中如何监控和处理这些错误？

**参考答案**：

**监控策略**：

```go
// 1. 结构化日志
type AgentMetrics struct {
    SessionID       string
    AgentName       string
    TotalTurns      int
    ToolCalls       int
    ToolFailures    int
    GuardrailHits   int
    MaxTurnsHit     bool
    Duration        time.Duration
    Error           error
}

func (r *Runner) Run(...) (*RunResult, error) {
    metrics := &AgentMetrics{...}
    defer func() {
        log.Info("agent_execution", metrics)  // 输出到日志系统
        prometheus.RecordMetrics(metrics)     // 输出到监控系统
    }()
    // ...
}
```

**错误处理策略**：

| 错误类型 | 处理方式 | 告警级别 |
|---------|---------|---------|
| `MaxTurnsExceededError` | 记录日志，返回友好提示"任务太复杂，请拆分" | WARN |
| `GuardrailTripwireTriggeredError` | 记录详情（用于安全审计），返回拒绝提示 | INFO/WARN |
| 工具执行失败 | 重试 1-2 次，仍失败则降级 | INFO |
| LLM API 错误 | 重试 + 指数退避，超限告警 | ERROR |

**告警规则示例**：

```yaml
# Prometheus AlertManager
groups:
  - name: agent_alerts
    rules:
      - alert: HighGuardrailTriggerRate
        expr: rate(guardrail_triggers_total[5m]) > 0.1
        labels:
          severity: warning
        annotations:
          summary: "护栏触发率过高，可能存在攻击"

      - alert: HighMaxTurnsRate
        expr: rate(max_turns_exceeded_total[5m]) > 0.05
        labels:
          severity: info
        annotations:
          summary: "MaxTurns 超限频繁，考虑调高阈值或优化 Agent"
```

---

### 第七部分：开放题（综合能力）

#### Q7.1 如何添加流式输出（Streaming）支持？

**设计思路**：

```go
// 1. 新增流式结果类型
type StreamChunk struct {
    Type    string  // "text", "tool_call", "done"
    Content string
    ToolCall *ToolCallChunk
}

// 2. 修改 Runner 接口
func (r *Runner) RunStream(ctx context.Context, agent *Agent, input Input) (<-chan StreamChunk, error)

// 3. 修改 LLM 调用
// 使用 OpenAI 的 stream=true 参数
stream, _ := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
    Stream: openai.Bool(true),
    // ...
})

for chunk := range stream {
    resultChan <- StreamChunk{
        Type:    "text",
        Content: chunk.Choices[0].Delta.Content,
    }
}

// 4. 护栏需要适配
// 输入护栏：正常执行（输入是完整的）
// 输出护栏：可选择累积完再检查，或流式检查（更复杂）
```

---

#### Q7.2 如何实现语义路由？

**设计思路**：

```go
type SemanticRouter struct {
    Embedder     EmbeddingClient       // Embedding 模型
    ToolEmbeddings map[string][]float64 // 预计算的工具描述向量
    TopN         int
}

func (r *SemanticRouter) RouteTools(ctx context.Context, input Input, tools []Tool) ([]Tool, error) {
    // 1. 计算输入的 embedding
    inputEmb, _ := r.Embedder.Embed(ctx, input.String())

    // 2. 计算与每个工具的相似度
    scores := make([]float64, len(tools))
    for i, tool := range tools {
        scores[i] = cosineSimilarity(inputEmb, r.ToolEmbeddings[tool.GetName()])
    }

    // 3. 返回 Top-N
    return selectTopN(tools, scores, r.TopN), nil
}

// 优化：工具描述的 embedding 可以预计算并缓存
```

---

#### Q7.3 对比 OpenAI Agents SDK（Python）和 agentgo

| 维度 | OpenAI Agents SDK (Python) | agentgo (Go) |
|------|---------------------------|--------------|
| **语言特性** | 动态类型、装饰器语法简洁 | 静态类型、接口驱动、编译时检查 |
| **并发模型** | asyncio | goroutine（更轻量） |
| **部署** | 需要 Python 环境 | 单二进制，无依赖 |
| **类型安全** | 运行时错误 | 编译时捕获 |
| **生态** | 丰富（LangChain 等） | 相对空白 |
| **性能** | 一般 | 更高（适合高并发） |
| **学习曲线** | 低 | 略高（需懂 Go） |

**各自优势**：
- Python：快速原型、丰富生态、社区大
- Go：生产部署、高并发、类型安全

---

## 十、面试表达技巧总结

### 通用框架：STAR-T

| 步骤 | 说明 | 示例 |
|------|------|------|
| **S**ituation | 背景/问题 | "当工具数量超过 20 个时..." |
| **T**ask | 目标 | "需要减少 Token 消耗" |
| **A**ction | 行动/设计 | "我实现了 KeywordRouter..." |
| **R**esult | 结果 | "节省了 30-50% Token" |
| **T**rade-off | 权衡 | "牺牲了语义理解能力，但满足了性能需求" |

### 表达原则

1. **先总后分**：先给结论/总览，再展开细节
2. **用数字说话**："13 个接口"、"5 层架构"、"节省 30% Token"
3. **举具体代码**：提到包名、接口名、方法名，证明你读过代码
4. **主动提亮点**：不等追问，主动说设计意图
5. **承认局限**：展示技术判断力，"这个方案的局限是..."

---

*文档更新时间：2025-01-02*
*适用岗位：大模型应用开发 / AI Agent 开发工程师*
