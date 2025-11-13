# Runner.go 实现指导文档

## 📌 文档信息

- **创建时间**: 2025-11-11
- **学习目标**: 手写完成 `runner.go` 的核心实现
- **学习方式**: 分步指导 + 自主编写 + 代码审查
- **当前状态**: 阶段 4 完成，准备阶段 5

---

## 🎯 项目目标

实现 `Runner.run()` 方法，这是整个 nvgo 框架的**执行引擎**，负责：

1. ✅ 调用 LLM（OpenAI Responses API）
2. ✅ 执行工具调用（包括 MCP 工具）
3. ✅ 管理 Agent 循环（最多 MaxTurns 次）
4. ✅ 处理 Guardrails（输入/输出护栏）
5. ✅ 管理 Session（会话历史）
6. ✅ 处理 Handoff（智能体切换）

---

## 📊 当前状态分析

### 已有代码
```go
// runner.go:149-151
func (r Runner) run(ctx context.Context, startingAgent *Agent, input Input) (*RunResult, error) {
    return nil, nil  // ❌ 完全未实现
}
```

### 已有结构体
- ✅ `Runner` - 运行器
- ✅ `RunConfig` - 运行配置
- ✅ `RunResult` - 运行结果
- ✅ `Usage` - 使用统计
- ✅ `ModelResponse` - 模型响应

### 缺失内容
- ❌ 错误类型（`MaxTurnsExceededError`, `GuardrailTripwireTriggeredError`）
- ❌ 完整的导入包
- ❌ `run()` 方法的实现
- ❌ 辅助函数（如工具执行、guardrail 运行等）

---

## 🗺️ 整体实施计划

```
阶段 1: 准备工作 (30分钟)
  └─ 添加导入包 + 定义错误类型

阶段 2: 基础框架 (1小时)
  └─ 实现 run() 的主循环结构（不含具体逻辑）

阶段 3: LLM 调用 (2-3小时)
  └─ 集成 OpenAI Responses API

阶段 4: 工具调用 (2-3小时)
  └─ 处理 function calls 和 MCP 工具

阶段 5: Guardrails (1-2小时)
  └─ 实现输入/输出护栏执行

阶段 6: Session (1小时)
  └─ 集成会话历史管理

阶段 7: 测试验证 (2-3小时)
  └─ 编写测试用例验证功能

总计: 约 10-15 小时
```

---

## 📝 阶段 1: 准备工作

### 🎯 目标
- 补充必要的导入包
- 定义缺失的错误类型

### 📋 任务清单

#### 任务 1.1: 补充导入包

**当前导入**:
```go
import (
	"context"

	"github.com/agent_go/memory"
	"github.com/openai/openai-go/v3/responses"
)
```

**需要添加的包**:

| 包名 | 用途 | 何时使用 |
|------|------|----------|
| `fmt` | 格式化错误消息 | 创建错误、日志 |
| `errors` | 错误处理 | 包装错误、错误判断 |
| `github.com/openai/openai-go/v3` | OpenAI 客户端 | 调用 Chat Completions API |
| `github.com/openai/openai-go/v3/option` | 请求选项 | 自定义 HTTP 请求 |

**导入格式规范**:
```go
import (
	// 标准库（按字母排序）
	"context"
	"errors"
	"fmt"

	// 第三方库（按字母排序）
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"

	// 本地包
	"github.com/agent_go/memory"
)
```

**🔨 动手任务**:
打开 `runner.go`，修改导入部分，添加缺失的包。

**✅ 验收标准**:
- [ ] 导入按标准库、第三方库、本地包分组
- [ ] 每组内按字母排序
- [ ] 没有未使用的导入（编译器会警告）

---

#### 任务 1.2: 定义错误类型

**需要定义的错误**:

##### 1. `MaxTurnsExceededError` - 超过最大循环次数

**接口要求**:
- 实现 `error` 接口（需要 `Error() string` 方法）
- 包含 `MaxTurns uint64` 字段（记录限制值）

**思路**:
```go
type MaxTurnsExceededError struct {
	MaxTurns uint64
}

func (e *MaxTurnsExceededError) Error() string {
	// TODO: 返回格式化的错误消息
	// 提示: 使用 fmt.Sprintf
}
```

**错误消息格式建议**:
```
"max turns exceeded: reached limit of 10 turns"
```

---

##### 2. `GuardrailTripwireTriggeredError` - 护栏触发

**接口要求**:
- 实现 `error` 接口
- 包含 `GuardrailName string` 字段（哪个护栏触发）
- 包含 `OutputInfo any` 字段（护栏的详细信息）
- 包含 `IsInput bool` 字段（是输入护栏还是输出护栏）

**思路**:
```go
type GuardrailTripwireTriggeredError struct {
	GuardrailName string
	OutputInfo    any
	IsInput       bool  // true = 输入护栏, false = 输出护栏
}

func (e *GuardrailTripwireTriggeredError) Error() string {
	// TODO: 返回格式化的错误消息
	// 提示: 区分输入/输出护栏
}
```

**错误消息格式建议**:
```
"input guardrail 'content_filter' triggered"
"output guardrail 'safety_check' triggered"
```

---

**🔨 动手任务**:
在 `runner.go` 文件末尾（在最后一个函数之后）添加这两个错误类型定义。

**提示**:
1. 放在文件末尾，保持代码组织清晰
2. 添加注释说明每个字段的用途
3. `Error()` 方法使用 `fmt.Sprintf` 格式化消息

**✅ 验收标准**:
- [ ] 两个错误类型都定义了
- [ ] 每个类型都有 `Error()` 方法
- [ ] 错误消息清晰易懂
- [ ] 代码能编译通过（`go build`）

---

### 🎓 知识点

#### Go 的错误接口
```go
type error interface {
    Error() string
}
```

任何实现了 `Error() string` 方法的类型都是 `error`。

#### 指针接收器 vs 值接收器
```go
// ✅ 推荐：使用指针接收器
func (e *MaxTurnsExceededError) Error() string {
    return fmt.Sprintf("...")
}

// ❌ 避免：值接收器（会复制整个结构体）
func (e MaxTurnsExceededError) Error() string {
    return fmt.Sprintf("...")
}
```

**原因**: 错误通常作为指针返回（`return &MaxTurnsExceededError{...}`），使用指针接收器更一致。

---

### 📖 参考代码位置

- **现有错误定义**: `error.go:6-12`
  ```go
  var (
      ErrMCPServerNotInitialized = errors.New("...")
      ErrMCPAgentRequired        = errors.New("...")
  )
  ```

- **Guardrail 结构**: `guardrail.go:13-37`
  ```go
  type GuardrailFunctionOutput struct {
      TripwireTriggered bool
      OutputInfo        any
  }
  ```

---

## 📝 阶段 2: 基础框架

### 🎯 目标
实现 `run()` 方法的整体结构，**不包含具体实现**，只搭建框架。

### 📋 任务清单

#### 任务 2.1: 实现主函数框架

**目标**: 实现 `run()` 方法的整体流程，用 `// TODO` 注释标记待实现部分。

**整体流程**:
```
1. 初始化 RunResult
2. 运行输入 Guardrails（如果有）
3. 主循环 (最多 MaxTurns 次)
   3.1 调用 LLM
   3.2 处理响应
   3.3 执行工具调用
   3.4 检查是否完成
   3.5 保存到 Session
4. 检查是否超过最大循环
5. 运行输出 Guardrails
6. 返回结果
```

**框架模板**:
```go
func (r Runner) run(ctx context.Context, startingAgent *Agent, input Input) (*RunResult, error) {
	// 1. 初始化结果
	result := &RunResult{
		Input:        CopyInput(input),
		NewItems:     []RunItem{},
		RawResponses: []ModelResponse{},
	}

	// 2. 运行输入 Guardrails
	// TODO: 实现输入 guardrail 逻辑

	// 3. 初始化循环变量
	currentAgent := startingAgent
	turnCount := uint64(0)
	maxTurns := r.Config.MaxTurns
	if maxTurns == 0 {
		maxTurns = DefaultMaxTurns
	}

	// 4. 主循环
	for turnCount < maxTurns {
		turnCount++

		// 4.1 调用 LLM
		// TODO: 实现 LLM 调用逻辑

		// 4.2 处理响应
		// TODO: 解析 LLM 响应

		// 4.3 执行工具调用
		// TODO: 执行工具并收集结果

		// 4.4 检查是否完成
		// TODO: 检查是否有最终输出

		// 4.5 保存到 Session
		// TODO: 将新项保存到 Session

		// 4.6 检查是否需要继续循环
		// TODO: 如果没有待处理的工具调用，退出循环
	}

	// 5. 检查是否超过最大循环次数
	if turnCount >= maxTurns {
		return nil, &MaxTurnsExceededError{MaxTurns: maxTurns}
	}

	// 6. 运行输出 Guardrails
	// TODO: 实现输出 guardrail 逻辑

	// 7. 设置最后的 Agent
	result.LastAgent = currentAgent

	return result, nil
}
```

**🔨 动手任务**:
1. 删除当前的 `return nil, nil`
2. 按照上面的框架填写代码
3. 保留所有 `// TODO` 注释

**✅ 验收标准**:
- [ ] 函数结构清晰，分为 7 个步骤
- [ ] 每个 TODO 位置都有注释说明要做什么
- [ ] 代码能编译通过
- [ ] `MaxTurns` 的默认值处理正确

---

#### 任务 2.2: 理解数据流

**问题**: 画出数据在 `run()` 方法中的流动图。

```
Input (用户输入)
  ↓
[输入 Guardrails]
  ↓
┌─────────────────────────┐
│  主循环 (最多 N 次)      │
│  ┌──────────────────┐   │
│  │ 1. 调用 LLM      │   │
│  │ 2. 获取响应      │   │
│  │ 3. 执行工具      │   │
│  │ 4. 检查完成      │   │
│  │ 5. 保存历史      │   │
│  └──────────────────┘   │
└─────────────────────────┘
  ↓
[输出 Guardrails]
  ↓
RunResult (最终输出)
```

**思考题**:
1. 为什么需要 `turnCount` 计数器？
2. 什么情况下会退出主循环？
3. `RunResult.NewItems` 和 `RunResult.RawResponses` 有什么区别？

---

### 🎓 知识点

#### 循环控制
```go
// ✅ 使用 for 条件循环
for turnCount < maxTurns {
    // 循环体
    turnCount++
}

// ❌ 避免无限循环
for {
    // 没有退出条件
}
```

#### 默认值处理
```go
// ✅ 使用零值检查
if maxTurns == 0 {
    maxTurns = DefaultMaxTurns
}

// ❌ 避免硬编码
maxTurns = 10  // 不灵活
```

---

## 📝 阶段 3: LLM 调用

### 🎯 目标
实现与 OpenAI API 的集成，调用 LLM 并获取响应。

### 📋 背景知识

#### OpenAI 的两种 API

NVGo 支持两种调用方式：

1. **Responses API** (推荐，如果设置了 Prompt)
   - 路径: `/v1/responses`
   - 使用预设的 Prompt 模板
   - 支持 Prompt 变量替换

2. **Chat Completions API** (传统方式)
   - 路径: `/v1/chat/completions`
   - 传统的消息数组格式

**选择逻辑**:
```go
if agent.Prompt != nil {
    // 使用 Responses API
} else {
    // 使用 Chat Completions API
}
```

---

### 📋 任务清单

#### 任务 3.1: 准备 LLM 调用参数

**需要收集的信息**:
1. **Model 名称** - 从 RunConfig 或 Agent 获取
2. **Instructions** - 系统提示词
3. **Tools** - 工具列表（包括 MCP 工具）
4. **ModelSettings** - 模型参数（temperature 等）
5. **历史消息** - 从 Session 加载（如果有）

**实现思路**:

```go
// 在主循环内部，替换 "// TODO: 实现 LLM 调用逻辑"

// 3.1.1 确定使用的模型
model := r.Config.Model
if model == "" {
    model = currentAgent.Model
}

// 3.1.2 获取 Instructions
var instructions string
if currentAgent.Instructions != nil {
    var err error
    instructions, err = currentAgent.Instructions.GetInstructions(ctx, currentAgent)
    if err != nil {
        return nil, fmt.Errorf("get instructions: %w", err)
    }
}

// 3.1.3 获取工具列表
var tools []Tool
// TODO: 从 Agent.MCPServers 获取 MCP 工具
// TODO: 合并 Agent 自身的工具（如果有）

// 3.1.4 合并 ModelSettings
modelSettings := currentAgent.ModelSettings.Resolve(r.Config.ModelSettings)

// 3.1.5 加载历史消息（从 Session）
var historyItems []responses.ResponseInputItemUnionParam
if r.Config.Session != nil {
    // TODO: 调用 Session.GetItems
}

// 3.1.6 构建当前输入
// TODO: 将 input 转换为 ResponseInputItemUnionParam
```

**🔨 动手任务**:
1. 在主循环中实现参数收集
2. 使用 `fmt.Errorf` 包装错误
3. 注意处理 nil 值（如 `Instructions`, `Session`）

**✅ 验收标准**:
- [ ] 模型名称优先使用 RunConfig，fallback 到 Agent
- [ ] Instructions 正确获取并处理错误
- [ ] 代码能编译通过

---

#### 任务 3.2: 调用 LLM

这部分比较复杂，我们**暂时使用占位符**，下一阶段再详细实现。

**占位符代码**:
```go
// 调用 LLM
var modelResponse ModelResponse
// TODO: 根据 agent.Prompt 是否为 nil 选择 API
// - 如果有 Prompt: 使用 Responses API
// - 如果没有: 使用 Chat Completions API
modelResponse = ModelResponse{
    Output: []responses.ResponseOutputItemUnion{
        // TODO: 临时返回空响应
    },
}

// 记录响应
result.RawResponses = append(result.RawResponses, modelResponse)
```

**🔨 动手任务**:
暂时只添加占位符，确保代码结构正确。

---

#### 任务 3.3: 辅助函数 - 获取 MCP 工具

**目标**: 实现一个辅助函数，从 MCP 服务器获取所有工具。

**函数签名**:
```go
// 在 run() 方法外部定义
func getMCPTools(ctx context.Context, agent *Agent, strict bool) ([]Tool, error) {
    if len(agent.MCPServers) == 0 {
        return nil, nil
    }

    // 提示: 使用 mcp.go 中的 GetAllFunctionTools
    return GetAllFunctionTools(ctx, agent.MCPServers, strict, agent)
}
```

**🔨 动手任务**:
1. 在 `runner.go` 文件末尾添加这个辅助函数
2. 理解 `GetAllFunctionTools` 的作用（在 mcp.go 中）

**✅ 验收标准**:
- [ ] 函数能正确调用 `GetAllFunctionTools`
- [ ] 处理空 MCPServers 的情况
- [ ] 返回类型正确

---

### 🎓 知识点

#### 错误包装
```go
// ✅ 推荐：使用 fmt.Errorf 和 %w
instructions, err := currentAgent.Instructions.GetInstructions(ctx, currentAgent)
if err != nil {
    return nil, fmt.Errorf("get instructions: %w", err)
}

// ❌ 避免：丢失上下文
if err != nil {
    return nil, err  // 不知道哪里出错
}
```

#### 配置合并
```go
// Resolve 方法会将 override 的非零值覆盖到当前设置
merged := agent.ModelSettings.Resolve(runConfig.ModelSettings)
```

---

## 📝 阶段 4: 工具调用

### 🎯 目标
处理 LLM 返回的 function calls，执行工具并将结果反馈。

### 📋 任务清单

#### 任务 4.1: 解析工具调用

**背景**: LLM 的响应中可能包含多种输出类型：
- `ResponseMessage` - 普通消息
- `ResponseFunctionCall` - 工具调用
- `ResponseHandoff` - Agent 切换
- 其他...

**实现思路**:
```go
// 在主循环中，替换 "// TODO: 解析 LLM 响应"

for _, outputItem := range modelResponse.Output {
    switch item := outputItem.(type) {
    case *responses.ResponseMessage:
        // TODO: 处理普通消息

    case *responses.ResponseFunctionCall:
        // TODO: 处理工具调用

    case *responses.ResponseHandoff:
        // TODO: 处理 Agent 切换

    default:
        // TODO: 其他类型
    }
}
```

**🔨 动手任务**:
添加类型判断的框架，每个 case 先用 TODO 标记。

---

#### 任务 4.2: 执行单个工具

**目标**: 实现一个辅助函数执行单个工具调用。

**函数签名**:
```go
func executeTool(ctx context.Context, tool Tool, arguments string) (any, error) {
    // 1. 转换为 FunctionTool 类型
    funcTool, ok := tool.(FunctionTool)
    if !ok {
        return nil, fmt.Errorf("tool is not a FunctionTool")
    }

    // 2. 检查工具是否启用
    if funcTool.IsEnabled != nil {
        // TODO: 调用 IsEnabled 接口检查
    }

    // 3. 执行工具
    result, err := funcTool.OnInvokeTool(ctx, arguments)

    // 4. 处理错误
    if err != nil {
        // TODO: 使用 FailureErrorFunction 或 DefaultToolErrorFunction
    }

    return result, nil
}
```

**🔨 动手任务**:
实现这个辅助函数，参考 `tool.go` 中的定义。

**✅ 验收标准**:
- [ ] 能正确执行 FunctionTool.OnInvokeTool
- [ ] 处理 IsEnabled 检查
- [ ] 正确使用错误处理函数

---

#### 任务 4.3: 查找工具

**目标**: 根据工具名称找到对应的工具对象。

**函数签名**:
```go
func findTool(tools []Tool, name string) (Tool, bool) {
    for _, t := range tools {
        if t.ToolName() == name {
            return t, true
        }
    }
    return nil, false
}
```

**🔨 动手任务**:
实现这个简单的查找函数。

---

### 🎓 知识点

#### 类型断言和类型开关
```go
// 类型开关
switch v := someInterface.(type) {
case *ConcreteType1:
    // v 的类型是 *ConcreteType1
case *ConcreteType2:
    // v 的类型是 *ConcreteType2
default:
    // 其他类型
}

// 类型断言
funcTool, ok := tool.(FunctionTool)
if !ok {
    // 类型不匹配
}
```

---

## 📝 阶段 5: Guardrails

### 🎯 目标
实现输入和输出护栏的执行逻辑。

### 📋 任务清单

#### 任务 5.1: 运行输入 Guardrails

**位置**: 在主循环之前

**实现思路**:
```go
// 替换 "// TODO: 实现输入 guardrail 逻辑"

// 合并 Runner 和 Agent 的输入 Guardrails
inputGuardrails := append(r.Config.InputGuardrails, startingAgent.InputGuardrails...)

// 运行所有输入 Guardrails
for _, gr := range inputGuardrails {
    grResult, err := gr.Run(ctx, startingAgent, input)
    if err != nil {
        return nil, fmt.Errorf("input guardrail %q failed: %w", gr.Name, err)
    }

    result.InputGuardrailResults = append(result.InputGuardrailResults, grResult)

    // 检查是否触发 tripwire
    if grResult.TripwireTriggered {
        return nil, &GuardrailTripwireTriggeredError{
            GuardrailName: gr.Name,
            OutputInfo:    grResult.OutputInfo,
            IsInput:       true,
        }
    }
}
```

**🔨 动手任务**:
实现输入 guardrails 的执行和 tripwire 检查。

**✅ 验收标准**:
- [ ] 合并了 RunConfig 和 Agent 的 guardrails
- [ ] 使用 `InputGuardrail.Run()` 方法（参考 guardrail.go）
- [ ] 正确处理 tripwire 触发

---

#### 任务 5.2: 运行输出 Guardrails

**位置**: 在主循环之后，返回之前

**实现思路**: 类似输入 guardrails，但：
- 只在有 `FinalOutput` 时运行
- 使用 `OutputGuardrail.Run()` 方法
- `IsInput` 字段设为 `false`

**🔨 动手任务**:
参考任务 5.1，自己实现输出 guardrails。

---

## 📝 阶段 6: Session 集成

### 🎯 目标
将会话历史管理集成到 run() 方法中。

### 📋 任务清单

#### 任务 6.1: 加载历史消息

**位置**: 在主循环内部，调用 LLM 之前

**实现思路**:
```go
// 替换 "// TODO: 调用 Session.GetItems"

var historyItems []responses.ResponseInputItemUnionParam
if r.Config.Session != nil {
    items, err := r.Config.Session.GetItems(ctx, -1)  // -1 = 获取全部
    if err != nil {
        return nil, fmt.Errorf("load session history: %w", err)
    }
    historyItems = items
}
```

**🔨 动手任务**:
实现历史消息加载，理解 `GetItems` 的参数含义（参考 memory/session.go）。

---

#### 任务 6.2: 保存新消息

**位置**: 在主循环内部，处理完响应后

**实现思路**:
```go
// 替换 "// TODO: 将新项保存到 Session"

if r.Config.Session != nil && len(result.NewItems) > 0 {
    // 转换 NewItems 为 ResponseInputItemUnionParam
    var itemsToSave []responses.ResponseInputItemUnionParam
    for _, item := range result.NewItems {
        itemsToSave = append(itemsToSave, item.ToInputItem())
    }

    // 保存到 Session
    if err := r.Config.Session.AddItems(ctx, itemsToSave); err != nil {
        return nil, fmt.Errorf("save to session: %w", err)
    }
}
```

**🔨 动手任务**:
实现消息保存逻辑。

**注意**: `RunItem` 接口有 `ToInputItem()` 方法（参考 runner.go:10-13）。

---

## 📝 阶段 7: 测试验证

### 🎯 目标
编写测试用例验证 Runner 的基本功能。

### 📋 任务清单

#### 任务 7.1: 测试最大循环次数

**测试文件**: `runner_test.go`

**测试用例**:
```go
func TestRunner_MaxTurns(t *testing.T) {
    // 创建一个永远不会完成的 Agent
    // 验证是否在 MaxTurns 次后返回错误

    // TODO: 实现测试
}
```

---

#### 任务 7.2: 测试错误类型

**测试用例**:
```go
func TestMaxTurnsExceededError(t *testing.T) {
    err := &MaxTurnsExceededError{MaxTurns: 10}

    // 验证错误消息
    expected := "max turns exceeded: reached limit of 10 turns"
    if err.Error() != expected {
        t.Errorf("expected %q, got %q", expected, err.Error())
    }
}
```

---

## 📊 进度追踪

### 当前进度

- [x] 阶段 1: 准备工作 ✅ **已完成**
  - [x] 任务 1.1: 补充导入包
  - [x] 任务 1.2: 定义错误类型

- [x] 阶段 2: 基础框架 ✅ **已完成**
  - [x] 任务 2.1: 实现主函数框架
  - [x] 任务 2.2: 理解数据流

- [x] 阶段 3: LLM 调用 ✅ **已完成**
  - [x] 任务 3.1: 准备 LLM 调用参数
  - [x] 任务 3.2: 调用 LLM（占位符）
  - [x] 任务 3.3: 实现 getMCPTools 辅助函数

- [x] 阶段 4: 工具调用 ✅ **已完成**
  - [x] 任务 4.1: 解析工具调用
  - [x] 任务 4.2: 执行单个工具
  - [x] 任务 4.3: 查找工具

- [ ] 阶段 5: Guardrails
  - [ ] 任务 5.1: 运行输入 Guardrails
  - [ ] 任务 5.2: 运行输出 Guardrails

- [ ] 阶段 6: Session
  - [ ] 任务 6.1: 加载历史消息
  - [ ] 任务 6.2: 保存新消息

- [ ] 阶段 7: 测试
  - [ ] 任务 7.1: 测试最大循环次数
  - [ ] 任务 7.2: 测试错误类型

### 完成记录

**格式**: `[日期] 完成任务 X.Y - 备注`

```
[2025-11-11] 完成任务 1.1 - 添加了必要的导入包（context, fmt, errors, openai-go, option）
[2025-11-11] 完成任务 1.2 - 定义了两个错误类型（MaxTurnsExceededError, GuardrailTripwireTriggeredError）
[2025-11-11] 阶段 1 完成 - 代码通过 go fmt 格式化和 go build 编译检查
[2025-11-12] 完成任务 2.1 - 实现了主函数框架，添加了主循环的 6 个子步骤 TODO
[2025-11-12] 完成任务 2.2 - 添加了输出 Guardrails 的 TODO 标记
[2025-11-12] 阶段 2 完成 - run() 方法的 7 步框架结构清晰完整
[2025-11-12] 完成任务 3.3 - 实现了 getMCPTools 辅助函数（调用 GetAllFunctionTools）
[2025-11-12] 完成任务 3.1 - 实现了 LLM 调用参数收集（模型名、Instructions、工具列表、ModelSettings、历史消息）
[2025-11-12] 完成任务 3.2 - 添加了 LLM 调用占位符，记录响应到 RawResponses
[2025-11-12] 阶段 3 完成 - 代码编译通过，所有参数准备就绪
[2025-11-12] 完成任务 4.3 - 实现了 findTool 辅助函数（根据工具名查找工具）
[2025-11-12] 完成任务 4.2 - 实现了 executeTool 辅助函数（类型转换、IsEnabled检查、工具执行、错误处理）
[2025-11-12] 完成任务 4.1 - 实现了工具调用解析（使用 AsAny() 类型判断、处理 ResponseFunctionToolCall）
[2025-11-12] 阶段 4 完成 - 创建了 RunItemWrapper 包装类型，使用 ResponseInputItemParamOfFunctionCallOutput 创建输出
[2025-11-12] 代码编译通过 - 所有工具调用逻辑实现完毕

```

---

## 🔍 调试技巧

### 编译检查
```bash
# 检查语法错误
go build ./...

# 运行测试
go test -v

# 查看类型信息
go doc github.com/agent_go.Runner
```

### 常见错误

#### 1. 导入未使用
```
imported and not used: "fmt"
```
**解决**: 删除未使用的导入，或者使用 `_ "package"` 占位。

#### 2. 类型不匹配
```
cannot use ... (type X) as type Y
```
**解决**: 检查类型定义，使用类型断言或转换。

#### 3. 接口未实现
```
X does not implement Y (missing method Z)
```
**解决**: 为类型添加缺失的方法。

---

## 📚 参考资料

### 项目内部文档
- `plan/plan.md` - 整体项目分析
- `replace/MIGRATION_SSE_TO_STREAMABLE.md` - MCP 迁移指南

### 相关代码位置
- `agent.go` - Agent 定义
- `tool.go` - Tool 接口和 FunctionTool
- `mcp.go` - MCP 服务器和工具
- `guardrail.go` - Guardrail 接口
- `memory/session.go` - Session 接口
- `setting.go` - ModelSettings

### Go 语言资源
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go 标准库](https://pkg.go.dev/std)

### OpenAI SDK
- [openai-go GitHub](https://github.com/openai/openai-go)
- [Responses API 文档](https://platform.openai.com/docs/api-reference/responses)
- [Chat Completions API](https://platform.openai.com/docs/api-reference/chat)

---

## 💡 学习建议

### 逐步推进
1. ✅ **先搭框架，再填细节** - 不要一次做完所有
2. ✅ **频繁编译测试** - 每完成一小步就 `go build`
3. ✅ **阅读现有代码** - 参考项目中已有的实现
4. ✅ **理解而非死记** - 搞清楚为什么这样写

### 遇到问题时
1. 🔍 **查看错误消息** - Go 的错误消息很清晰
2. 📖 **阅读文档** - godoc, 官方文档
3. 🔬 **调试打印** - 使用 `fmt.Printf` 调试
4. 🤝 **寻求帮助** - 把错误贴给我分析

### 代码质量
1. ✨ **添加注释** - 说明复杂逻辑
2. 🧪 **编写测试** - 验证功能正确性
3. 🔧 **重构优化** - 第一版不用完美，可以后续改进
4. 📏 **遵循规范** - 参考项目现有代码风格

---

## 🎯 下一步行动

**现在开始**: 阶段 5 - Guardrails

### 阶段概述

阶段 5 实现输入和输出护栏（Guardrails）的执行逻辑，用于在运行前后对内容进行检查和过滤。

### 主要任务

#### 任务 5.1: 运行输入 Guardrails
在主循环之前（第 181-182 行），替换：
```go
// 2. 运行输入 Guardrails
// TODO: 实现输入 guardrail 逻辑
```

实现：
- 合并 RunConfig 和 Agent 的输入 Guardrails
- 运行所有输入 Guardrails
- 检查是否触发 tripwire
- 保存 guardrail 结果到 RunResult

#### 任务 5.2: 运行输出 Guardrails
在主循环之后（约第 217-218 行），替换：
```go
// 6. 运行输出 Guardrails
// TODO: 实现输出 guardrail 逻辑
```

实现：
- 合并 RunConfig 和 Agent 的输出 Guardrails
- 只在有 FinalOutput 时运行
- 检查是否触发 tripwire
- 保存 guardrail 结果到 RunResult

### 验收标准
- [ ] 输入 Guardrails 正确执行
- [ ] 输出 Guardrails 正确执行
- [ ] Tripwire 触发时返回正确的错误
- [ ] 代码能编译通过

### 📖 参考文档
详细实现指导请参考 md 文档第 714-773 行（阶段 5 的详细说明）。

**祝你学习愉快！** 🚀

---

**文档版本**: v1.4
**最后更新**: 2025-11-12（阶段 4 已完成，准备阶段 5）
**下次审查**: 完成阶段 5 任务后
