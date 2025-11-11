# NVGo 项目全面分析与开发计划

## 📋 项目概述

**NVGo** 是一个用 Go 语言编写的**优雅的多智能体工作流框架**，灵感来自 OpenAI 的 Agents Python SDK 和 NVIDIA 的 NeMo Agent Toolkit。

### 核心定位
- 提供多智能体编排能力
- 支持 Model Context Protocol (MCP) 集成
- 与 OpenAI API 深度集成
- 提供 guardrails（护栏）机制用于输入/输出验证

### 项目信息
- **仓库**: `github.com/demo/nvgo`
- **Go 版本**: 1.25
- **许可证**: MIT License (Copyright 2025 qntx.sol)
- **状态**: 🚧 早期开发阶段，核心功能未完成

---

## 🏗️ 架构设计

### 设计模式
1. **Builder 模式**: Agent 配置使用链式调用
2. **接口抽象**: Tool, Prompter, InstructionsGetter 等都是接口
3. **依赖注入**: 通过接口和配置对象传递依赖
4. **函数式编程**: 支持函数式的动态配置（InstructionsFunc, DynamicPromptFunction）

### 设计亮点
1. **灵活的工具系统**: 支持动态启用/禁用、错误处理、严格模式
2. **多层次的配置覆盖**: Agent 级别和 Run 级别的 ModelSettings 合并
3. **MCP 工具过滤**: 白名单/黑名单静态过滤 + 可扩展的动态过滤接口
4. **会话管理**: 支持内存和持久化存储，限制条数等高级功能
5. **类型安全**: 大量使用 Go 的类型系统确保编译时安全

---

## 📁 文件结构与功能详解

### 🎯 核心文件

#### 1. `agent.go` - 智能体核心定义
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\agent.go`
**行数**: 167 行
**完成度**: ✅ 95%

**核心结构体 - Agent**:
```go
type Agent struct {
    Name                string                    // 智能体名称
    Instructions        InstructionsGetter        // 系统提示词（可动态生成）
    Prompt              Prompter                  // OpenAI Responses API 的提示配置
    Model               string                    // 使用的模型名称
    Client              openai.Client             // OpenAI 客户端
    ModelSettings       ModelSettings             // 模型参数配置
    MCPServers          []MCPServer               // MCP 服务器列表
    MCPConfig           MCPConfig                 // MCP 配置
    InputGuardrails     []InputGuardrail          // 输入护栏
    OutputGuardrails    []OutputGuardrail         // 输出护栏
    OutputType          OutputTypeInterface       // 输出类型定义
}
```

**Builder 方法**:
- `New(name string) *Agent` - 创建新智能体
- `WithInstructions(instr string) *Agent` - 设置指令
- `WithInstructionsFunc(fn InstructionsFunc) *Agent` - 动态指令
- `WithPrompt(prompt Prompter) *Agent` - 设置提示
- `WithModel(model string) *Agent` - 设置模型
- `WithClient(client openai.Client) *Agent` - 设置客户端
- `WithModelSettings(settings ModelSettings) *Agent` - 设置模型参数
- `WithMCPServers(mcpServers []MCPServer) *Agent` - 批量设置 MCP 服务器
- `AddMCPServer(mcpServer MCPServer) *Agent` - 添加单个 MCP 服务器
- `WithMCPConfig(mcpConfig MCPConfig) *Agent` - 设置 MCP 配置
- `WithInputGuardrails(gr []InputGuardrail) *Agent` - 批量设置输入护栏
- `AddInputGuardrail(gr InputGuardrail) *Agent` - 添加单个输入护栏
- `WithOutputGuardrails(gr []OutputGuardrail) *Agent` - 批量设置输出护栏
- `AddOutputGuardrail(gr OutputGuardrail) *Agent` - 添加单个输出护栏
- `WithOutputType(outputType OutputTypeInterface) *Agent` - 设置输出类型

**注意事项**:
- 第 128-130 行有注释掉的 `AddMCPStdioServer` 方法，需要决定是否实现

---

#### 2. `runner.go` - 运行器（🚨 核心未完成）
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\runner.go`
**行数**: 152 行
**完成度**: ❌ 10%

**核心结构体**:

```go
// Runner 执行器
type Runner struct {
    Config RunConfig
}

// RunConfig 运行配置
type RunConfig struct {
    Model               string                // 全局模型覆盖
    ModelSettings       ModelSettings         // 全局模型设置覆盖
    InputGuardrails     []InputGuardrail      // 全局输入护栏
    OutputGuardrails    []OutputGuardrail     // 全局输出护栏
    WorkflowName        string                // 工作流名称（用于追踪）
    MaxTurns            uint64                // 最大循环次数（默认 10）
    PreviousResponseID  string                // 上一次响应 ID（OpenAI Responses API）
    Session             memory.Session        // 会话对象
}

// RunResult 运行结果
type RunResult struct {
    Input                   Input                      // 原始输入
    NewItems                []RunItem                  // 新生成的项
    RawResponses            []ModelResponse            // 原始 LLM 响应
    FinalOutput             any                        // 最终输出
    InputGuardrailResults   []InputGuardrailResult     // 输入护栏结果
    OutputGuardrailResults  []OutputGuardrailResult    // 输出护栏结果
    LastAgent               *Agent                     // 最后运行的智能体
}

// Usage 使用统计
type Usage struct {
    Requests            uint64                                         // 请求次数
    InputTokens         uint64                                         // 输入 tokens
    InputTokensDetails  responses.ResponseUsageInputTokensDetails      // 输入详情
    OutputTokens        uint64                                         // 输出 tokens
    OutputTokensDetails responses.ResponseUsageOutputTokensDetails     // 输出详情
    TotalTokens         uint64                                         // 总 tokens
}
```

**执行流程设计**:
```go
// Run 执行工作流（外部 API）
func (r Runner) Run(ctx context.Context, startingAgent *Agent, input string) (*RunResult, error)

// run 内部执行逻辑（🚨 未实现）
func (r Runner) run(ctx context.Context, startingAgent *Agent, input Input) (*RunResult, error) {
    return nil, nil  // ❌ 当前只返回 nil
}
```

**预期执行流程**:
1. 调用 Agent 的 LLM，传入输入和上下文
2. 如果产生最终输出（匹配 Agent.OutputType），结束循环
3. 如果有 handoff（切换到其他 Agent），使用新 Agent 重新循环
4. 如果有工具调用，执行工具并将结果反馈给 LLM，继续循环
5. 检查是否超过 MaxTurns
6. 运行 guardrails 检查

**🚨 关键缺失**:
- 完整的循环逻辑
- LLM 调用实现
- 工具调用机制
- Handoff 机制
- MaxTurnsExceededError 实现
- GuardrailTripwireTriggeredError 实现

---

#### 3. `tool.go` - 工具定义
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\tool.go`
**行数**: 79 行
**完成度**: ✅ 100%

**核心接口**:
```go
type Tool interface {
    ToolName() string
    isTool()
}
```

**FunctionTool 结构体**:
```go
type FunctionTool struct {
    Name                  string                                                // 工具名称
    Description           string                                                // 工具描述
    ParamsJSONSchema      map[string]any                                       // 参数 JSON Schema
    OnInvokeTool          func(ctx context.Context, arguments string) (any, error)  // 执行函数
    FailureErrorFunction  *ToolErrorFunction                                   // 错误处理函数
    StrictJSONSchema      param.Opt[bool]                                      // 是否使用严格 JSON Schema
    IsEnabled             FunctionToolEnabler                                  // 动态启用接口
}
```

**错误处理**:
```go
type ToolErrorFunction func(ctx context.Context, err error) (any, error)

// 默认错误处理：返回错误信息给 LLM
func DefaultToolErrorFunction(_ context.Context, err error) (any, error) {
    return fmt.Sprintf("An error occurred while running the tool. Please try again. Error: %s", err), nil
}
```

**设计亮点**:
- 支持动态启用/禁用工具
- 可自定义错误处理策略
- 支持严格模式 JSON Schema（提高 LLM 输入准确性）

---

#### 4. `mcp.go` - MCP 集成（最复杂）
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\mcp.go`
**行数**: 399 行
**完成度**: ✅ 90%

**核心接口**:
```go
type MCPServer interface {
    Connect(context.Context) error                                                      // 连接服务器
    Cleanup(context.Context) error                                                      // 清理资源
    Name() string                                                                       // 服务器名称
    UseStructuredContent() bool                                                         // 是否使用结构化内容
    ListTools(context.Context, *Agent) ([]*mcp.Tool, error)                            // 列出工具
    CallTool(context.Context, string, map[string]any) (*mcp.CallToolResult, error)    // 调用工具
    ListPrompts(context.Context) (*mcp.ListPromptsResult, error)                       // 列出提示
    GetPrompt(context.Context, string, map[string]string) (*mcp.GetPromptResult, error) // 获取提示
}
```

**工具过滤机制**:

```go
// 过滤接口
type MCPToolFilter interface {
    FilterMCPTool(ctx context.Context, filterCtx MCPToolFilterContext, tool *mcp.Tool) (bool, error)
}

// 静态过滤器（白名单/黑名单）
type MCPToolFilterStatic struct {
    AllowedToolNames []string  // 白名单
    BlockedToolNames []string  // 黑名单
}
```

**三种传输方式**:

1. **MCPServerStdio** - 标准输入输出传输
   ```go
   func NewMCPServerStdio(p MCPServerStdioParams) *MCPServerStdio
   ```

2. **MCPServerSSE** - Server-Sent Events（已废弃）
   ```go
   func NewMCPServerSSE(p MCPServerSSEParams) *MCPServerSSE
   ```

3. **MCPServerStreamableHTTP** - 可流式 HTTP
   ```go
   func NewMCPServerStreamableHTTP(p MCPServerStreamableHTTPParams) *MCPServerStreamableHTTP
   ```

**关键功能**:

1. **工具列表缓存**:
   - `cacheToolsList` - 是否缓存
   - `cacheDirty` - 缓存是否失效
   - `InvalidateToolsCache()` - 手动失效缓存

2. **MCP 工具转 FunctionTool**:
   ```go
   func ToFunctionTool(tool *mcp.Tool, server MCPServer, strict bool) (FunctionTool, error)
   ```

3. **批量获取工具**:
   ```go
   func GetAllFunctionTools(ctx context.Context, servers []MCPServer, strict bool, agent *Agent) ([]Tool, error)
   ```
   - 检查工具名称去重
   - 应用过滤器

4. **工具调用**:
   ```go
   func InvokeMCPTool(ctx context.Context, server MCPServer, tool *mcp.Tool, input string) (string, error)
   ```
   - 支持结构化内容和普通内容
   - 自动序列化/反序列化 JSON

**线程安全**:
- 使用 `sync.Mutex` 保护 `Cleanup` 操作

---

### 🛡️ 安全与验证

#### 5. `guardrail.go` - 护栏机制
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\guardrail.go`
**行数**: 102 行
**完成度**: ✅ 80%

**核心结构**:

```go
// 护栏输出
type GuardrailFunctionOutput struct {
    TripwireTriggered bool  // 是否触发紧急停止
    OutputInfo        any   // 检查详情信息
}

// 输入护栏
type InputGuardrail struct {
    Name          string
    GuardrailFunc func(ctx context.Context, agent *Agent, input Input) (GuardrailFunctionOutput, error)
}

// 输出护栏
type OutputGuardrail struct {
    Name          string
    GuardrailFunc func(ctx context.Context, agent *Agent, output any) (GuardrailFunctionOutput, error)
}
```

**执行方法**:
```go
func (g InputGuardrail) Run(ctx context.Context, agent *Agent, input Input) (InputGuardrailResult, error)
func (g OutputGuardrail) Run(ctx context.Context, agent *Agent, output any) (OutputGuardrailResult, error)
```

**设计理念**:
- 输入护栏在智能体执行前并行运行（只对第一个 Agent）
- 输出护栏在最终输出后运行
- Tripwire 触发时会中止整个工作流
- 可携带详细检查信息供调试

**使用场景**:
- 输入护栏: 敏感词检测、越权检查、格式验证
- 输出护栏: 内容审核、合规性检查、质量评估

---

### 🧠 记忆系统

#### 6. `memory/session.go` - 会话接口
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\memory\session.go`
**行数**: 28 行
**完成度**: ✅ 100%

**接口定义**:
```go
type Session interface {
    // 获取历史记录
    // limit <= 0: 返回全部（升序）
    // limit > 0: 返回最新 N 条（升序）
    GetItems(ctx context.Context, limit int) ([]responses.ResponseInputItemUnionParam, error)

    // 添加新项
    AddItems(ctx context.Context, items []responses.ResponseInputItemUnionParam) error

    // 弹出最新项（用于撤销等场景）
    PopItem(context.Context) (*responses.ResponseInputItemUnionParam, error)

    // 清空会话
    ClearSession(context.Context) error
}
```

---

#### 7. `memory/sqlite.go` - SQLite 实现
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\memory\sqlite.go`
**行数**: 339 行
**完成度**: ✅ 100%

**核心结构**:
```go
type SQLiteSession struct {
    sessionID     string       // 会话 ID
    db            *sql.DB      // 数据库连接
    sessionsTable string       // 会话表名
    messagesTable string       // 消息表名
    isMemoryDB    bool         // 是否为内存数据库
    mu            sync.Mutex   // 互斥锁
}
```

**数据库 Schema**:

```sql
-- 会话元数据表
CREATE TABLE IF NOT EXISTS agent_sessions (
    session_id TEXT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 消息数据表
CREATE TABLE IF NOT EXISTS agent_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    message_data TEXT NOT NULL,  -- JSON 序列化
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES agent_sessions (session_id) ON DELETE CASCADE
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_agent_messages_session_id
ON agent_messages (session_id, created_at);
```

**关键特性**:

1. **配置灵活**:
   ```go
   type SQLiteSessionConfig struct {
       SessionID     string  // 必需
       DBPath        string  // 默认 ":memory:"
       SessionsTable string  // 默认 "agent_sessions"
       MessagesTable string  // 默认 "agent_messages"
   }
   ```

2. **线程安全**: 所有操作都使用 `mu.Lock()` 保护

3. **事务支持**: `AddItems` 使用事务确保原子性

4. **高级查询**:
   - 限制返回条数
   - 自动翻转顺序（获取最新 N 条时）
   - 使用 RETURNING 子句优化 PopItem

5. **优雅关闭**:
   ```go
   func (s *SQLiteSession) Close() error
   ```

6. **错误处理**: 所有数据库错误都包装为语义化错误

---

#### 8. `memory/errors.go` - 错误定义
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\memory\errors.go`
**行数**: 23 行
**完成度**: ✅ 100%

**定义的错误**:
```go
var (
    ErrInvalidSessionID   = errors.New("session ID is required")
    ErrDatabaseOpen       = errors.New("failed to open database")
    ErrDatabaseInit       = errors.New("failed to initialize database schema")
    ErrSessionNotFound    = errors.New("session not found")
    ErrInvalidItemData    = errors.New("invalid item data")
    ErrTransactionFailed  = errors.New("database transaction failed")
    ErrOperationFailed    = errors.New("database operation failed")
    ErrDatabaseClose      = errors.New("failed to close database")
)
```

---

### 🔧 配置与辅助

#### 9. `setting.go` - 模型设置
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\setting.go`
**行数**: 162 行
**完成度**: ✅ 100%

**完整的 ModelSettings 结构**:
```go
type ModelSettings struct {
    // 基础参数
    Temperature      param.Opt[float64]  `json:"temperature"`
    TopP             param.Opt[float64]  `json:"top_p"`
    FrequencyPenalty param.Opt[float64]  `json:"frequency_penalty"`
    PresencePenalty  param.Opt[float64]  `json:"presence_penalty"`
    MaxTokens        param.Opt[int64]    `json:"max_tokens"`

    // 工具相关
    ToolChoice        ToolChoice         `json:"tool_choice"`
    ParallelToolCalls param.Opt[bool]    `json:"parallel_tool_calls"`

    // 高级参数
    Truncation       param.Opt[Truncation]            `json:"truncation"`
    Reasoning        openai.ReasoningParam            `json:"reasoning"`

    // 元数据
    Metadata         map[string]string                `json:"metadata"`
    Store            param.Opt[bool]                  `json:"store"`
    IncludeUsage     param.Opt[bool]                  `json:"include_usage"`
    ResponseInclude  []responses.ResponseIncludable   `json:"response_include"`

    // HTTP 参数
    ExtraQuery       map[string]string                `json:"extra_query"`
    ExtraHeaders     map[string]string                `json:"extra_headers"`

    // 自定义钩子
    CustomizeResponsesRequest      func(context.Context, *responses.ResponseNewParams, []option.RequestOption) (*responses.ResponseNewParams, []option.RequestOption, error) `json:"-"`
    CustomizeChatCompletionsRequest func(context.Context, *openai.ChatCompletionNewParams, []option.RequestOption) (*openai.ChatCompletionNewParams, []option.RequestOption, error) `json:"-"`
}
```

**ToolChoice 枚举**:
```go
type ToolChoice interface {
    isToolChoice()
}

type ToolChoiceString string
const (
    ToolChoiceAuto     ToolChoiceString = "auto"
    ToolChoiceRequired ToolChoiceString = "required"
    ToolChoiceNone     ToolChoiceString = "none"
)

type ToolChoiceMCP struct {
    ServerLabel string `json:"server_label"`
    Name        string `json:"name"`
}
```

**配置合并**:
```go
// 将 override 的非零值覆盖到当前设置
func (ms ModelSettings) Resolve(override ModelSettings) ModelSettings
```

---

#### 10. `prompt.go` - 提示配置
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\prompt.go`
**行数**: 64 行
**完成度**: ✅ 100%

**核心结构**:
```go
type Prompt struct {
    ID        string                                            // 提示 ID
    Version   param.Opt[string]                                // 提示版本
    Variables map[string]responses.ResponsePromptVariableUnionParam  // 变量替换
}

type Prompter interface {
    Prompt(context.Context, *Agent) (Prompt, error)
}

type DynamicPromptFunction func(context.Context, *Agent) (Prompt, error)
```

**工具函数**:
```go
// 转换为 OpenAI API 参数
func (promptUtil) ToModelInput(
    ctx context.Context,
    prompter Prompter,
    agent *Agent,
) (responses.ResponsePromptParam, bool, error)
```

**使用场景**:
- 使用 OpenAI Responses API 的静态提示
- 根据 Agent 状态动态生成提示

---

#### 11. `instruction.go` - 指令系统
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\instruction.go`
**行数**: 31 行
**完成度**: ✅ 100%

**接口定义**:
```go
type InstructionsGetter interface {
    GetInstructions(context.Context, *Agent) (string, error)
}
```

**实现类型**:

1. **静态指令**:
   ```go
   type InstructionsStr string
   func (s InstructionsStr) GetInstructions(context.Context, *Agent) (string, error)
   ```

2. **动态指令**:
   ```go
   type InstructionsFunc func(context.Context, *Agent) (string, error)
   func (fn InstructionsFunc) GetInstructions(ctx context.Context, a *Agent) (string, error)
   ```

**使用示例**:
```go
// 静态
agent.WithInstructions("You are a helpful assistant")

// 动态
agent.WithInstructionsFunc(func(ctx context.Context, a *Agent) (string, error) {
    return fmt.Sprintf("Current time: %s", time.Now()), nil
})
```

---

#### 12. `input.go` - 输入类型
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\input.go`
**行数**: 39 行
**完成度**: ✅ 100%

**接口定义**:
```go
type Input interface {
    isInput()
}
```

**实现类型**:

1. **简单字符串**:
   ```go
   type InputString string
   ```

2. **结构化输入项**:
   ```go
   type InputItems []responses.ResponseInputItemUnionParam
   ```

**辅助函数**:
```go
func CopyInput(input Input) Input  // 深拷贝输入
```

---

#### 13. `output.go` - 输出类型
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\output.go`
**行数**: 30 行
**完成度**: ✅ 100%

**接口定义**:
```go
type OutputTypeInterface interface {
    IsPlainText() bool                                    // 是否为纯文本
    Name() string                                         // 输出类型名称
    JSONSchema() (map[string]any, error)                 // JSON Schema
    IsStrictJSONSchema() bool                            // 是否为严格模式
    ValidateJSON(ctx context.Context, jsonStr string) (any, error)  // 验证 JSON
}
```

**使用场景**:
- 定义智能体的结构化输出格式
- 利用 OpenAI 的 Structured Outputs 功能
- 验证 LLM 输出的 JSON 格式

---

#### 14. `error.go` - 错误定义
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\error.go`
**行数**: 12 行
**完成度**: ✅ 100%

**定义的错误**:
```go
var (
    ErrMCPServerNotInitialized = errors.New("server not initialized: make sure you call `Connect()` first")
    ErrMCPAgentRequired        = errors.New("agent is required for dynamic tool filtering")
)
```

---

### 🔍 内部工具

#### 15. `internal/strictschema/` - JSON Schema 严格化
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\internal\strictschema\`
**完成度**: ✅ 100%

**主文件**: `strictschema.go` (156 行)

**核心函数**:
```go
func EnsureStrictJSONSchema(schema map[string]any) (map[string]any, error)
```

**转换规则**:

1. **空 Schema 补全**:
   ```json
   {}
   →
   {
     "type": "object",
     "additionalProperties": false,
     "properties": {},
     "required": []
   }
   ```

2. **自动添加 additionalProperties**:
   - 对所有 `type: "object"` 添加 `additionalProperties: false`

3. **从属性生成 required**:
   ```json
   {
     "type": "object",
     "properties": {"a": {...}, "b": {...}}
   }
   →
   {
     "type": "object",
     "properties": {"a": {...}, "b": {...}},
     "required": ["a", "b"]  // 自动添加
   }
   ```

4. **解析 $ref 引用**:
   - 当 `$ref` 与其他字段共存时，展开引用并合并

5. **移除 null 默认值**:
   ```json
   {"default": null}  →  {}
   ```

6. **递归处理**:
   - `properties` 内的嵌套对象
   - `items` 数组项
   - `anyOf`, `allOf` 联合/交集
   - `$defs`, `definitions` 定义区

**文档**: `README.md` 包含详细示例

**测试**: `strictschema_test.go`

---

#### 16. `internal/transform/` - 命名转换
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\internal\transform\`
**完成度**: ✅ 100%

**主文件**: `transform.go` (140 行)

**功能**:
```go
// 命名约定
type NamingConvention string
const (
    SnakeCase NamingConvention = "snake_case"
    CamelCase NamingConvention = "camelCase"
)

// 转换函数
func ToCase(name string) string  // 使用当前约定
func ApplyCase(name string, convention NamingConvention) string  // 指定约定
func ToCamelCase(name string) string
func ToSnakeCase(name string) string
func TransformStringFunctionStyle(name string) string  // 函数名规范化
```

**环境变量配置**:
- `OPENAI_AGENTS_NAMING_CONVENTION`: 设置为 "snake_case" 或 "camelCase"
- 默认: "snake_case"

**示例**:
```go
ToSnakeCase("getUserInfo")    // → "get_user_info"
ToCamelCase("get_user_info")  // → "getUserInfo"
ToCamelCase("GetUserInfo")    // → "getUserInfo"
```

**测试**: `transform_test.go`

---

### 🧪 测试文件

#### 17. `mcp_test.go` - MCP 测试
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\mcp_test.go`
**行数**: 168 行
**完成度**: ✅ 80%

**测试覆盖**:

1. **TestNewMCPToolFilterStatic** (第 13-65 行)
   - 测试静态过滤器的创建
   - 验证白名单/黑名单配置

2. **TestMCPToolFilterStatic_FilterMCPTool** (第 67-125 行)
   - 测试过滤逻辑
   - 验证白名单/黑名单优先级

3. **TestApplyMCPToolFilter** (第 127-167 行)
   - 测试批量过滤
   - 验证过滤结果

**缺失测试**:
- MCP 服务器连接测试
- 工具调用集成测试
- 工具缓存测试

---

#### 18. `runner_test.go` - 运行器测试
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\runner_test.go`
**行数**: 2 行
**完成度**: ❌ 0%

**现状**: 完全为空，只有 `package nvgo`

**需要添加的测试**:
- Runner 执行流程测试
- MaxTurns 限制测试
- Guardrails 测试
- Session 集成测试
- 工具调用测试
- Handoff 测试

---

### 📦 配置文件

#### 19. `go.mod` / `go.sum` - 依赖管理
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\go.mod`

**依赖列表**:
```go
require (
    github.com/mattn/go-sqlite3 v1.14.32                // SQLite 驱动
    github.com/modelcontextprotocol/go-sdk v0.3.0       // MCP SDK
    github.com/openai/openai-go/v2 v2.1.1               // OpenAI Go SDK
    github.com/stretchr/testify v1.11.0                 // 测试框架
)

require (
    github.com/davecgh/go-spew v1.1.1                   // 结构体打印
    github.com/google/jsonschema-go v0.2.0              // JSON Schema
    github.com/pmezard/go-difflib v1.0.0                // Diff 工具
    github.com/tidwall/gjson v1.14.4                    // JSON 查询
    github.com/tidwall/match v1.1.1                     // 字符串匹配
    github.com/tidwall/pretty v1.2.1                    // JSON 美化
    github.com/tidwall/sjson v1.2.5                     // JSON 设置
    github.com/yosida95/uritemplate/v3 v3.0.2           // URI 模板
    gopkg.in/yaml.v3 v3.0.1                            // YAML 解析
)
```

---

#### 20. `Makefile` - 构建工具
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\Makefile`
**行数**: 58 行

**可用命令**:
```bash
make          # 运行 pre-commit（默认）
make deps     # 安装和整理依赖
make fmt      # 格式化代码
make lint     # 代码检查（需要 golangci-lint）
make test     # 运行测试
make docs     # 启动 godoc 服务器（localhost:6060）
make pre-commit  # 运行所有检查
make help     # 显示帮助信息
```

---

#### 21. `.golangci.yml` - Linter 配置
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\.golangci.yml`
**行数**: 130 行

**启用的 Linter (80+)**:
- 代码质量: staticcheck, govet, errcheck, unused, ineffassign
- 安全性: gosec
- 命名: revive, misspell, godot
- 性能: prealloc, perfsprint
- 测试: testifylint, testableexamples
- 更多...

**禁用的 Linter**:
- exhaustruct (不要求所有字段初始化)
- funlen (不限制函数长度)
- gochecknoglobals (允许全局变量)
- varnamelen (不限制变量名长度)
- 等...

**格式化工具**:
- gci: import 排序
- gofmt: 标准格式化
- gofumpt: 更严格的格式化
- goimports: import 管理
- golines: 长行分割

**格式化规则**:
```yaml
gofmt:
  rewrite-rules:
    - pattern: "interface{}"
      replacement: "any"
```

---

#### 22. `LICENSE` - MIT 许可证
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\LICENSE`
**版权所有者**: qntx.sol (2025)

---

### 🤖 CI/CD

#### 23. `.github/workflows/go.yml` - Go CI
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\.github\workflows\go.yml`
**行数**: 12 行
**完成度**: ❌ 30%

**触发条件**:
```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

**🚨 缺失**: 具体的 CI 步骤（构建、测试、linting）

---

#### 24. `.github/workflows/stale.yml` - 过期 Issue 清理
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\.github\workflows\stale.yml`
**行数**: 8 行
**完成度**: ❌ 30%

**触发条件**:
```yaml
on:
  schedule:
    - cron: "30 1 * * *"  # 每天 01:30 UTC
  workflow_dispatch:
```

**🚨 缺失**: 具体的清理步骤

---

#### 25. `.github/dependabot.yml` - Dependabot 配置
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\.github\dependabot.yml`
**行数**: 12 行
**完成度**: ✅ 100%

**配置**:
```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
```

---

#### 26. `README.md` - 项目文档
**路径**: `e:\Lab\work\develop\write_agent\nvgo-main\README.md`
**行数**: 20 行
**完成度**: ❌ 15%

**现有内容**:
- Logo
- 项目标题
- 安装命令
- 致谢链接

**🚨 缺失**:
- 快速开始指南
- 使用示例
- API 文档
- 架构说明
- 贡献指南

---

## 📊 项目成熟度评估

| 模块 | 文件 | 完成度 | 说明 |
|------|------|--------|------|
| **核心类型定义** | agent.go, tool.go, input.go, output.go, instruction.go, prompt.go | ✅ 95% | 架构完整，缺少少量实现 |
| **MCP 集成** | mcp.go | ✅ 90% | 功能完整，测试覆盖不足 |
| **运行器** | runner.go | ❌ 10% | **核心逻辑完全缺失** |
| **护栏机制** | guardrail.go | ✅ 80% | 框架完整，缺实际使用示例 |
| **记忆系统** | memory/*.go | ✅ 100% | SQLite 实现完整且健壮 |
| **模型设置** | setting.go | ✅ 100% | 支持所有 OpenAI 参数 |
| **内部工具** | internal/* | ✅ 100% | StrictSchema 和 Transform 都完整 |
| **测试** | *_test.go | ⚠️ 20% | 只有部分单元测试 |
| **文档** | README.md | ❌ 15% | 严重不足 |
| **CI/CD** | .github/workflows/* | ⚠️ 30% | 配置不完整 |

**总体完成度**: 约 **50%**

---

## 🚨 关键缺失功能详解

### 1. Runner 核心逻辑 (最高优先级)

**位置**: `runner.go:149-151`

**当前代码**:
```go
func (r Runner) run(ctx context.Context, startingAgent *Agent, input Input) (*RunResult, error) {
    return nil, nil  // ❌ 未实现
}
```

**需要实现的逻辑**:

```go
func (r Runner) run(ctx context.Context, startingAgent *Agent, input Input) (*RunResult, error) {
    // 1. 初始化运行结果
    result := &RunResult{
        Input:      CopyInput(input),
        NewItems:   []RunItem{},
        RawResponses: []ModelResponse{},
    }

    // 2. 运行输入护栏（只对第一个 Agent）
    if len(r.Config.InputGuardrails) > 0 || len(startingAgent.InputGuardrails) > 0 {
        // 合并 Runner 和 Agent 的输入护栏
        // 并行运行所有输入护栏
        // 检查是否有 TripwireTriggered
    }

    // 3. 初始化循环变量
    currentAgent := startingAgent
    currentInput := input
    turnCount := uint64(0)
    maxTurns := r.Config.MaxTurns
    if maxTurns == 0 {
        maxTurns = DefaultMaxTurns
    }

    // 4. 主循环
    for turnCount < maxTurns {
        turnCount++

        // 4.1 构建 LLM 请求参数
        // - 合并 Agent 和 Runner 的 ModelSettings
        // - 获取 Instructions
        // - 获取 Prompt（如果有）
        // - 获取工具列表（包括 MCP 工具）
        // - 构建历史消息（从 Session 加载）

        // 4.2 调用 LLM
        // - 根据是否有 Prompt 决定使用 Responses API 还是 Chat Completions API
        // - 使用 CustomizeResponsesRequest / CustomizeChatCompletionsRequest 钩子

        // 4.3 处理响应
        response := // LLM 返回的响应
        result.RawResponses = append(result.RawResponses, response)

        // 4.4 解析输出
        for _, outputItem := range response.Output {
            switch item := outputItem.(type) {
            case *responses.ResponseMessage:
                // 普通消息
                result.NewItems = append(result.NewItems, ...)

            case *responses.ResponseFunctionCall:
                // 工具调用
                // 执行工具
                // 将结果添加到 NewItems

            case *responses.ResponseHandoff:
                // 切换 Agent
                // 查找目标 Agent
                // 更新 currentAgent

            // 其他类型...
            }
        }

        // 4.5 检查是否有最终输出
        if currentAgent.OutputType != nil {
            // 尝试解析最终输出
            // 如果成功，设置 result.FinalOutput 并退出循环
        }

        // 4.6 保存到 Session
        if r.Config.Session != nil {
            r.Config.Session.AddItems(ctx, result.NewItems)
        }

        // 4.7 如果没有工具调用也没有 handoff，退出循环
        if /* 没有待处理的工具调用和 handoff */ {
            break
        }
    }

    // 5. 检查是否超过最大循环次数
    if turnCount >= maxTurns {
        return nil, &MaxTurnsExceededError{MaxTurns: maxTurns}
    }

    // 6. 运行输出护栏
    if result.FinalOutput != nil {
        // 合并 Runner 和 Agent 的输出护栏
        // 运行所有输出护栏
        // 检查是否有 TripwireTriggered
    }

    // 7. 设置最后的 Agent
    result.LastAgent = currentAgent

    return result, nil
}
```

**需要实现的辅助类型**:
```go
type MaxTurnsExceededError struct {
    MaxTurns uint64
}

func (e *MaxTurnsExceededError) Error() string {
    return fmt.Sprintf("max turns exceeded: %d", e.MaxTurns)
}

type GuardrailTripwireTriggeredError struct {
    GuardrailName string
    OutputInfo    any
}

func (e *GuardrailTripwireTriggeredError) Error() string {
    return fmt.Sprintf("guardrail tripwire triggered: %s", e.GuardrailName)
}
```

---

### 2. 工具调用机制

**需要实现**:
- 从 LLM 响应中提取工具调用
- 调用对应的 `FunctionTool.OnInvokeTool`
- 处理工具错误（使用 `FailureErrorFunction`）
- 将工具结果格式化为 LLM 可理解的格式
- 支持并行工具调用（根据 `ParallelToolCalls` 设置）

**伪代码**:
```go
func executeTool(ctx context.Context, tool FunctionTool, arguments string) (any, error) {
    // 1. 检查工具是否启用
    if tool.IsEnabled != nil {
        enabled, err := tool.IsEnabled.IsEnabled(ctx, agent)
        if err != nil || !enabled {
            return nil, err
        }
    }

    // 2. 执行工具
    result, err := tool.OnInvokeTool(ctx, arguments)

    // 3. 处理错误
    if err != nil {
        if tool.FailureErrorFunction != nil {
            return (*tool.FailureErrorFunction)(ctx, err)
        }
        return DefaultToolErrorFunction(ctx, err)
    }

    return result, nil
}
```

---

### 3. Handoff (智能体切换) 机制

**需要实现**:
- 定义 Handoff 类型
- 从 LLM 响应中识别 handoff 请求
- 查找目标 Agent（需要 Agent 注册机制）
- 切换到新 Agent 并传递上下文

**可能的设计**:
```go
// 在 Agent 中添加
type Agent struct {
    // ...
    HandoffDescription string  // 用于 LLM 选择 handoff 目标
}

// 在 Runner 中添加
type RunConfig struct {
    // ...
    AvailableAgents map[string]*Agent  // Agent 名称到 Agent 的映射
}

// Handoff 逻辑
func handleHandoff(handoff *responses.ResponseHandoff, availableAgents map[string]*Agent) (*Agent, error) {
    targetAgent, ok := availableAgents[handoff.TargetAgent]
    if !ok {
        return nil, fmt.Errorf("handoff target agent not found: %s", handoff.TargetAgent)
    }
    return targetAgent, nil
}
```

---

### 4. 测试覆盖

**需要添加的测试**:

1. **Runner 基础测试**:
   ```go
   func TestRunner_Run_SimpleMessage(t *testing.T)
   func TestRunner_Run_WithTools(t *testing.T)
   func TestRunner_Run_WithHandoff(t *testing.T)
   func TestRunner_Run_MaxTurns(t *testing.T)
   func TestRunner_Run_WithSession(t *testing.T)
   ```

2. **Guardrails 测试**:
   ```go
   func TestInputGuardrail_TripwireTriggered(t *testing.T)
   func TestOutputGuardrail_TripwireTriggered(t *testing.T)
   ```

3. **MCP 集成测试**:
   ```go
   func TestMCPServer_Connect(t *testing.T)
   func TestMCPServer_CallTool(t *testing.T)
   func TestMCPServer_ToolCaching(t *testing.T)
   ```

4. **Memory 集成测试**:
   ```go
   func TestSQLiteSession_Integration(t *testing.T)
   func TestSession_WithRunner(t *testing.T)
   ```

---

### 5. CI/CD 完善

**Go CI Workflow**:
```yaml
name: Go CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.25']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Install dependencies
        run: make deps

      - name: Format check
        run: |
          make fmt
          git diff --exit-code

      - name: Lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest

      - name: Test
        run: make test

      - name: Build
        run: go build -v ./...
```

**Stale Workflow**:
```yaml
name: Close Stale Issues and PRs

on:
  schedule:
    - cron: "30 1 * * *"
  workflow_dispatch:

jobs:
  stale:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/stale@v9
        with:
          stale-issue-message: 'This issue is stale because it has been open 60 days with no activity.'
          stale-pr-message: 'This PR is stale because it has been open 60 days with no activity.'
          days-before-stale: 60
          days-before-close: 7
```

---

### 6. 文档完善

**README.md 应包含**:

1. **快速开始**:
   ```markdown
   ## Quick Start
   
   ```go
   package main
   
   import (
       "context"
       "fmt"
       "github.com/demo/nvgo"
       "github.com/openai/openai-go/v2"
   )
   
   func main() {
       client := openai.NewClient()
   
       agent := nvgo.New("assistant").
           WithInstructions("You are a helpful assistant.").
           WithModel("gpt-4").
           WithClient(client)
   
       result, err := nvgo.Run(context.Background(), agent, "Hello!")
       if err != nil {
           panic(err)
       }
   
       fmt.Println(result.FinalOutput)
   }
   ```
   ```

2. **功能示例**:
   - 添加工具
   - 使用 MCP 服务器
   - 配置 Guardrails
   - 使用 Session
   - 多智能体 Handoff

3. **架构说明**:
   - 组件关系图
   - 执行流程图
   - 数据流图

---

## 🎯 开发路线图

### Phase 1: 核心功能实现 (P0 - 关键阻塞)

**目标**: 让项目可以运行最基本的 Agent 工作流

#### 1.1 实现 Runner 核心逻辑 ⏱️ 预计 3-5 天
- [ ] 实现 `Runner.run()` 方法
- [ ] 实现 LLM 调用逻辑
- [ ] 实现消息历史管理
- [ ] 实现基础循环控制

**验收标准**:
```go
// 可以运行这个最简单的例子
agent := nvgo.New("test").
    WithInstructions("You are a helpful assistant.").
    WithModel("gpt-4").
    WithClient(client)

result, err := nvgo.Run(ctx, agent, "Hello!")
// 应该返回有效的响应
```

#### 1.2 实现工具调用机制 ⏱️ 预计 2-3 天
- [ ] 从 LLM 响应解析工具调用
- [ ] 执行 `FunctionTool.OnInvokeTool`
- [ ] 处理工具错误
- [ ] 格式化工具结果
- [ ] 支持串行工具调用

**验收标准**:
```go
tool := nvgo.FunctionTool{
    Name: "get_weather",
    Description: "Get weather information",
    OnInvokeTool: func(ctx context.Context, args string) (any, error) {
        return "Sunny, 25°C", nil
    },
}

agent := nvgo.New("test").
    WithInstructions("Use tools to answer questions.").
    WithModel("gpt-4").
    AddTool(tool)  // 需要添加这个方法

result, err := nvgo.Run(ctx, agent, "What's the weather?")
// 应该调用工具并返回结果
```

#### 1.3 实现错误类型 ⏱️ 预计 1 天
- [ ] 实现 `MaxTurnsExceededError`
- [ ] 实现 `GuardrailTripwireTriggeredError`
- [ ] 实现其他必要的错误类型

#### 1.4 基础测试 ⏱️ 预计 2-3 天
- [ ] 添加 `TestRunner_Run_SimpleMessage`
- [ ] 添加 `TestRunner_Run_WithTools`
- [ ] 添加 `TestRunner_Run_MaxTurns`
- [ ] 修复所有发现的 bug

**Phase 1 总时间**: 约 8-12 天

---

### Phase 2: 高级功能 (P1 - 重要功能)

**目标**: 实现多智能体协作和完整的护栏机制

#### 2.1 实现 Handoff 机制 ⏱️ 预计 2-3 天
- [ ] 设计 Handoff 类型
- [ ] 实现 Agent 注册机制
- [ ] 实现 Agent 切换逻辑
- [ ] 添加 Handoff 测试

**验收标准**:
```go
salesAgent := nvgo.New("sales").
    WithInstructions("Handle sales questions")

supportAgent := nvgo.New("support").
    WithInstructions("Handle support questions")

router := nvgo.New("router").
    WithInstructions("Route to appropriate agent").
    WithHandoffTargets(salesAgent, supportAgent)

result, err := nvgo.Run(ctx, router, "I need support")
// 应该自动切换到 supportAgent
```

#### 2.2 完善 Guardrails ⏱️ 预计 2 天
- [ ] 实现 Guardrail 并行执行
- [ ] 实现 Tripwire 检查
- [ ] 添加内置 Guardrail 示例
- [ ] 添加 Guardrail 测试

#### 2.3 Session 集成 ⏱️ 预计 1-2 天
- [ ] 在 Runner 中集成 Session
- [ ] 自动保存历史
- [ ] 从 Session 加载历史
- [ ] 添加 Session 集成测试

#### 2.4 并行工具调用 ⏱️ 预计 1-2 天
- [ ] 实现并行执行逻辑
- [ ] 根据 `ParallelToolCalls` 控制
- [ ] 添加并行工具测试

#### 2.5 文档编写 ⏱️ 预计 2-3 天
- [ ] 完善 README.md
- [ ] 添加快速开始指南
- [ ] 添加 API 文档
- [ ] 添加架构说明
- [ ] 添加使用示例

**Phase 2 总时间**: 约 8-12 天

---

### Phase 3: 增强与优化 (P2 - 体验增强)

**目标**: 提升项目的可用性和稳定性

#### 3.1 测试覆盖完善 ⏱️ 预计 3-4 天
- [ ] MCP 集成测试
- [ ] StrictSchema 边界测试
- [ ] Transform 边界测试
- [ ] 端到端测试
- [ ] 压力测试

#### 3.2 CI/CD 完善 ⏱️ 预计 1-2 天
- [ ] 完善 Go CI workflow
- [ ] 完善 Stale workflow
- [ ] 添加测试覆盖率报告
- [ ] 添加 Release workflow

#### 3.3 日志与追踪 ⏱️ 预计 2-3 天
- [ ] 集成 structured logging
- [ ] 添加 OpenTelemetry 支持
- [ ] 实现 WorkflowName 追踪
- [ ] 添加性能指标

#### 3.4 示例项目 ⏱️ 预计 2-3 天
- [ ] 创建 `examples/` 目录
- [ ] 简单对话示例
- [ ] 工具调用示例
- [ ] MCP 集成示例
- [ ] 多智能体协作示例

#### 3.5 性能优化 ⏱️ 预计 2-3 天
- [ ] 优化内存分配
- [ ] 优化并发性能
- [ ] 添加性能基准测试
- [ ] 优化 Session 查询

**Phase 3 总时间**: 约 10-15 天

---

### Phase 4: 高级特性 (P3 - 可选增强)

**目标**: 添加高级功能和生态集成

#### 4.1 高级 Output Types ⏱️ 预计 2-3 天
- [ ] 实现内置的 OutputType
- [ ] 支持流式输出
- [ ] 支持增量解析

#### 4.2 更多 Guardrail 类型 ⏱️ 预计 2 天
- [ ] 内容审核 Guardrail
- [ ] 成本控制 Guardrail
- [ ] 延迟控制 Guardrail

#### 4.3 插件系统 ⏱️ 预计 3-4 天
- [ ] 设计插件接口
- [ ] 实现插件加载机制
- [ ] 添加插件示例

#### 4.4 可视化工具 ⏱️ 预计 3-5 天
- [ ] 工作流可视化
- [ ] 执行追踪 UI
- [ ] 调试工具

**Phase 4 总时间**: 约 10-14 天

---

## 📋 开发检查清单

### 立即可做 (Phase 1)

#### 核心实现
- [ ] 实现 `Runner.run()` 方法
  - [ ] 初始化 RunResult
  - [ ] 实现主循环
  - [ ] 调用 LLM
  - [ ] 处理响应
  - [ ] 检查最终输出
  - [ ] 保存到 Session
- [ ] 实现工具调用
  - [ ] 解析工具调用
  - [ ] 执行工具
  - [ ] 处理错误
  - [ ] 格式化结果
- [ ] 实现 MaxTurnsExceededError
- [ ] 实现 GuardrailTripwireTriggeredError

#### 测试
- [ ] TestRunner_Run_SimpleMessage
- [ ] TestRunner_Run_WithTools
- [ ] TestRunner_Run_MaxTurns

#### 代码清理
- [ ] 决定是否实现 `agent.go:128-130` 的 AddMCPStdioServer
- [ ] 添加缺失的包级文档注释

---

### 近期任务 (Phase 2)

#### 功能实现
- [ ] 实现 Handoff 机制
- [ ] 完善 Guardrails 执行
- [ ] 集成 Session 到 Runner
- [ ] 支持并行工具调用

#### 测试
- [ ] TestRunner_Run_WithHandoff
- [ ] TestRunner_Run_WithSession
- [ ] TestInputGuardrail_TripwireTriggered
- [ ] TestOutputGuardrail_TripwireTriggered
- [ ] TestRunner_Run_ParallelTools

#### 文档
- [ ] 完善 README.md
- [ ] 添加快速开始指南
- [ ] 添加 API 使用示例
- [ ] 编写架构文档

---

### 中期任务 (Phase 3)

#### 测试
- [ ] 添加 MCP 集成测试
- [ ] 添加端到端测试
- [ ] 提高测试覆盖率到 80%+

#### CI/CD
- [ ] 完善 Go CI workflow
- [ ] 完善 Stale workflow
- [ ] 添加覆盖率报告
- [ ] 添加 Release workflow

#### 工具
- [ ] 集成 structured logging
- [ ] 添加 OpenTelemetry
- [ ] 实现追踪功能

#### 示例
- [ ] 创建 examples/ 目录
- [ ] 添加各种示例代码

---

### 长期任务 (Phase 4)

#### 高级功能
- [ ] 实现高级 Output Types
- [ ] 添加更多内置 Guardrails
- [ ] 设计插件系统
- [ ] 开发可视化工具

#### 生态
- [ ] 发布到 pkg.go.dev
- [ ] 编写博客文章
- [ ] 创建视频教程
- [ ] 建立社区

---

## 🔍 代码问题与建议

### 现有问题

1. **agent.go:128-130** - 注释掉的代码
   ```go
   // Add MCPStdioServer appends an MCP server to the agent's MCP server list.
   // func (a *Agent) AddMCPStdioServer() *Agent {
   // 	return
   // }
   ```
   **建议**: 决定是否需要这个方法，如果不需要就删除

2. **runner.go:149-151** - 核心逻辑缺失
   ```go
   func (r Runner) run(ctx context.Context, startingAgent *Agent, input Input) (*RunResult, error) {
       return nil, nil
   }
   ```
   **建议**: 这是最高优先级，必须实现

3. **runner_test.go** - 空文件
   **建议**: 添加基础测试用例

4. **缺少 Agent.AddTool 方法**
   **建议**: 添加便捷方法直接添加工具
   ```go
   func (a *Agent) AddTool(tool Tool) *Agent {
       // 需要在 Agent 中添加 Tools 字段
       return a
   }
   ```

5. **缺少 Agent 注册机制**
   **建议**: 为 Handoff 添加 Agent 注册功能
   ```go
   func (a *Agent) WithHandoffTargets(agents ...*Agent) *Agent {
       // 需要在 Agent 中添加 HandoffTargets 字段
       return a
   }
   ```

---

### 设计建议

1. **Agent 结构扩展**
   ```go
   type Agent struct {
       // ... 现有字段
   
       // 新增字段
       Tools            []Tool              // 直接配置的工具
       HandoffTargets   map[string]*Agent   // Handoff 目标 Agent
       HandoffDescription string            // Handoff 描述
   }
   ```

2. **RunConfig 扩展**
   ```go
   type RunConfig struct {
       // ... 现有字段
   
       // 新增字段
       EnableTracing    bool                 // 启用追踪
       Logger           Logger               // 日志记录器
       MaxConcurrentTools int                // 最大并行工具数
   }
   ```

3. **错误处理增强**
   ```go
   type RunError struct {
       Agent     *Agent
       TurnCount uint64
       Cause     error
   }
   
   func (e *RunError) Error() string
   func (e *RunError) Unwrap() error
   ```

4. **事件回调**
   ```go
   type RunCallbacks struct {
       OnAgentStart   func(ctx context.Context, agent *Agent)
       OnAgentEnd     func(ctx context.Context, agent *Agent, output any)
       OnToolStart    func(ctx context.Context, tool Tool, arguments string)
       OnToolEnd      func(ctx context.Context, tool Tool, result any, err error)
       OnHandoff      func(ctx context.Context, from *Agent, to *Agent)
   }
   ```

---

## 💡 使用示例（预期）

### 示例 1: 简单对话

```go
package main

import (
    "context"
    "fmt"
    "github.com/demo/nvgo"
    "github.com/openai/openai-go/v2"
)

func main() {
    client := openai.NewClient()

    agent := nvgo.New("assistant").
        WithInstructions("You are a helpful assistant.").
        WithModel("gpt-4").
        WithClient(client)

    result, err := nvgo.Run(context.Background(), agent, "Hello!")
    if err != nil {
        panic(err)
    }

    fmt.Println(result.FinalOutput)
}
```

---

### 示例 2: 使用工具

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/demo/nvgo"
    "github.com/openai/openai-go/v2"
)

func main() {
    client := openai.NewClient()

    // 定义天气工具
    weatherTool := nvgo.FunctionTool{
        Name:        "get_weather",
        Description: "Get the current weather for a location",
        ParamsJSONSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "location": map[string]any{
                    "type":        "string",
                    "description": "The city name",
                },
            },
            "required": []string{"location"},
        },
        OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
            var args struct {
                Location string `json:"location"`
            }
            if err := json.Unmarshal([]byte(arguments), &args); err != nil {
                return nil, err
            }

            // 模拟天气查询
            return fmt.Sprintf("The weather in %s is sunny, 25°C", args.Location), nil
        },
    }

    agent := nvgo.New("assistant").
        WithInstructions("You are a helpful assistant with access to weather information.").
        WithModel("gpt-4").
        WithClient(client).
        AddTool(weatherTool)  // 需要实现

    result, err := nvgo.Run(context.Background(), agent, "What's the weather in New York?")
    if err != nil {
        panic(err)
    }

    fmt.Println(result.FinalOutput)
}
```

---

### 示例 3: 使用 MCP 服务器

```go
package main

import (
    "context"
    "fmt"
    "github.com/demo/nvgo"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    "github.com/openai/openai-go/v2"
)

func main() {
    client := openai.NewClient()

    // 创建 MCP 服务器
    mcpServer := nvgo.NewMCPServerStdio(nvgo.MCPServerStdioParams{
        Transport: &mcp.CommandTransport{
            Command: mcp.Command{
                Path: "npx",
                Args: []string{"-y", "@modelcontextprotocol/server-filesystem", "/path/to/files"},
            },
        },
        CommonMCPServerParams: nvgo.CommonMCPServerParams{
            Name: "filesystem",
            CacheToolsList: true,
        },
    })

    // 连接服务器
    if err := mcpServer.Connect(context.Background()); err != nil {
        panic(err)
    }
    defer mcpServer.Cleanup(context.Background())

    agent := nvgo.New("file-assistant").
        WithInstructions("You can read and write files.").
        WithModel("gpt-4").
        WithClient(client).
        AddMCPServer(mcpServer)

    result, err := nvgo.Run(context.Background(), agent, "List all files in the directory")
    if err != nil {
        panic(err)
    }

    fmt.Println(result.FinalOutput)
}
```

---

### 示例 4: 多智能体协作

```go
package main

import (
    "context"
    "fmt"
    "github.com/demo/nvgo"
    "github.com/openai/openai-go/v2"
)

func main() {
    client := openai.NewClient()

    // 销售 Agent
    salesAgent := nvgo.New("sales").
        WithInstructions("You are a sales representative. Handle pricing and product questions.").
        WithModel("gpt-4").
        WithClient(client)

    // 技术支持 Agent
    supportAgent := nvgo.New("support").
        WithInstructions("You are a technical support specialist. Handle technical issues.").
        WithModel("gpt-4").
        WithClient(client)

    // 路由 Agent
    routerAgent := nvgo.New("router").
        WithInstructions("You are a customer service router. Direct customers to the right department.").
        WithModel("gpt-4").
        WithClient(client).
        WithHandoffTargets(salesAgent, supportAgent)  // 需要实现

    result, err := nvgo.Run(context.Background(), routerAgent, "My product is not working")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Handled by: %s\n", result.LastAgent.Name)
    fmt.Println(result.FinalOutput)
}
```

---

### 示例 5: 使用 Guardrails

```go
package main

import (
    "context"
    "fmt"
    "github.com/demo/nvgo"
    "github.com/openai/openai-go/v2"
    "strings"
)

func main() {
    client := openai.NewClient()

    // 内容审核 Guardrail
    contentGuardrail := nvgo.NewInputGuardrail("content_filter",
        func(ctx context.Context, agent *nvgo.Agent, input nvgo.Input) (nvgo.GuardrailFunctionOutput, error) {
            inputStr, ok := input.(nvgo.InputString)
            if !ok {
                return nvgo.GuardrailFunctionOutput{}, nil
            }

            // 检查敏感词
            if strings.Contains(strings.ToLower(string(inputStr)), "badword") {
                return nvgo.GuardrailFunctionOutput{
                    TripwireTriggered: true,
                    OutputInfo:        "Content contains prohibited words",
                }, nil
            }

            return nvgo.GuardrailFunctionOutput{}, nil
        },
    )

    agent := nvgo.New("assistant").
        WithInstructions("You are a helpful assistant.").
        WithModel("gpt-4").
        WithClient(client).
        AddInputGuardrail(contentGuardrail)

    result, err := nvgo.Run(context.Background(), agent, "Hello!")
    if err != nil {
        if _, ok := err.(*nvgo.GuardrailTripwireTriggeredError); ok {
            fmt.Println("Input blocked by guardrail")
            return
        }
        panic(err)
    }

    fmt.Println(result.FinalOutput)
}
```

---

### 示例 6: 使用 Session (会话记忆)

```go
package main

import (
    "context"
    "fmt"
    "github.com/demo/nvgo"
    "github.com/demo/nvgo/memory"
    "github.com/openai/openai-go/v2"
)

func main() {
    client := openai.NewClient()

    // 创建 Session
    session, err := memory.NewSQLiteSession(context.Background(), memory.SQLiteSessionConfig{
        SessionID: "user-123",
        DBPath:    "./conversations.db",
    })
    if err != nil {
        panic(err)
    }
    defer session.Close()

    agent := nvgo.New("assistant").
        WithInstructions("You are a helpful assistant with memory of past conversations.").
        WithModel("gpt-4").
        WithClient(client)

    runner := nvgo.Runner{
        Config: nvgo.RunConfig{
            Session: session,
        },
    }

    // 第一轮对话
    result1, _ := runner.Run(context.Background(), agent, "My name is Alice")
    fmt.Println(result1.FinalOutput)

    // 第二轮对话（记住名字）
    result2, _ := runner.Run(context.Background(), agent, "What's my name?")
    fmt.Println(result2.FinalOutput)  // 应该回答 "Alice"
}
```

---

## 🛠️ 开发环境配置

### 必需工具

1. **Go 1.25+**
   ```bash
   go version
   ```

2. **golangci-lint**
   ```bash
   # macOS
   brew install golangci-lint
   
   # Linux
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
   
   # Windows
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

3. **SQLite** (内置在 mattn/go-sqlite3)

### 可选工具

1. **godoc** (查看文档)
   ```bash
   go install golang.org/x/tools/cmd/godoc@latest
   ```

2. **pre-commit** (Git hooks)
   ```bash
   # macOS
   brew install pre-commit
   
   # Linux/Windows
   pip install pre-commit
   ```

### IDE 配置

**VS Code 推荐插件**:
- Go (golang.go)
- golangci-lint (golangci.golangci-lint)
- Go Test Explorer (premparihar.go-test-explorer)

**GoLand / IntelliJ IDEA**:
- 原生支持，无需额外配置

---

## 📚 学习资源

### 官方文档
- [OpenAI Agents Python](https://github.com/openai/openai-agents-python) - Python 版本参考
- [OpenAI API 文档](https://platform.openai.com/docs)
- [MCP 协议](https://modelcontextprotocol.io)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)

### 相关项目
- [NeMo Agent Toolkit](https://github.com/NVIDIA/NeMo-Agent-Toolkit) - NVIDIA 的智能体框架
- [LangChain Go](https://github.com/tmc/langchaingo) - LangChain 的 Go 实现

---

## 🤝 贡献指南

### 分支策略
- `main`: 稳定分支
- `develop`: 开发分支
- `feature/*`: 功能分支
- `fix/*`: 修复分支

### 提交规范
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档
- `style`: 格式
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具

**示例**:
```
feat(runner): implement core run loop

- Add LLM invocation logic
- Implement tool calling mechanism
- Add max turns check

Closes #123
```

### Pull Request 流程
1. Fork 项目
2. 创建 feature 分支
3. 编写代码和测试
4. 运行 `make pre-commit`
5. 提交 PR
6. 等待 Review

---

## 📈 项目指标

### 当前状态
- **代码行数**: ~2,500 行
- **测试覆盖率**: ~15%
- **文档完成度**: ~15%
- **功能完成度**: ~50%

### 目标状态
- **代码行数**: ~5,000 行
- **测试覆盖率**: >80%
- **文档完成度**: >90%
- **功能完成度**: 100%

---

## 🔗 相关链接

- **GitHub**: (待发布)
- **文档**: (待完善)
- **示例**: (待创建)
- **讨论**: (待建立)

---

## 📝 总结

NVGo 是一个设计良好但实现未完成的多智能体框架。架构清晰，类型系统完善，但缺少最关键的执行引擎实现。

**优势**:
- ✅ 清晰的架构设计
- ✅ 完整的类型系统
- ✅ 良好的 MCP 集成
- ✅ 健壮的 Memory 系统
- ✅ 灵活的配置机制

**挑战**:
- ❌ Runner 核心逻辑缺失
- ❌ 测试覆盖不足
- ❌ 文档严重不足
- ❌ 缺少使用示例

**下一步**:
1. 实现 Runner.run() 方法 (最高优先级)
2. 实现工具调用机制
3. 添加基础测试
4. 完善文档

按照本计划，预计 **4-6 周**可以完成 Phase 1 和 Phase 2，使项目达到可用状态。

---

**文档版本**: v1.0
**最后更新**: 2025-10-28
**维护者**: (待定)w
