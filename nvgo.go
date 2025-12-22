// Package nvgo 是一个 Go 语言的 AI Agent 框架，参考 OpenAI Agents SDK 设计。
//
// nvgo 提供了构建智能 Agent 应用所需的全部功能，包括：
//   - Agent 定义与配置
//   - 工具调用（FunctionTool 和 MCP）
//   - 输入/输出护栏
//   - 会话管理
//   - 多 Agent 协作
//   - ReAct 模式
//   - 工具路由
//
// 快速开始：
//
//	client := openai.NewClient()
//	agent := nvgo.New("助手").
//		WithInstructions("你是一个友好的 AI 助手").
//		WithModel("gpt-4o-mini").
//		WithClient(client)
//
//	result, err := nvgo.Run(ctx, agent, "你好！")
//
// 更多示例请查看 examples 目录。
package nvgo

import (
	"nvgo/pkg/agent"
	"nvgo/pkg/memory"
	"nvgo/pkg/pattern"
	"nvgo/pkg/runner"
	"nvgo/pkg/tool"
	"nvgo/pkg/types"
	"nvgo/pkg/config"
)

// ========== Agent ==========

// Agent 代表一个 AI 模型的配置，包含指令、工具、护栏等。
type Agent = agent.Agent

// New 创建一个新的 Agent 实例。
func New(name string) *Agent {
	return agent.New(name)
}

// ========== Instructions ==========

// Instructions 定义如何获取 Agent 的系统提示词。
type Instructions = agent.Instructions

// InstructionsStr 是字符串形式的静态指令。
type InstructionsStr = agent.InstructionsStr

// InstructionsFunc 是函数形式的动态指令。
type InstructionsFunc = agent.InstructionsFunc

// StateProvider 提供动态状态供指令模板使用。
type StateProvider = agent.StateProvider

// MemoryStateProvider 是基于内存的状态提供者。
type MemoryStateProvider = agent.MemoryStateProvider

// NewMemoryStateProvider 创建一个新的内存状态提供者。
func NewMemoryStateProvider() *MemoryStateProvider {
	return agent.NewMemoryStateProvider()
}

// DynamicInstruction 支持模板化的动态指令。
type DynamicInstruction = agent.DynamicInstruction

// ========== Prompt (Responses API) ==========

// Prompt 配置 OpenAI Responses API 的提示参数。
type Prompt = agent.Prompt

// Prompter 定义如何获取 Prompt 配置。
type Prompter = agent.Prompter

// DynamicPromptFunction 是动态生成 Prompt 的函数。
type DynamicPromptFunction = agent.DynamicPromptFunction

// ========== ModelSettings ==========

// ModelSettings 包含模型参数（temperature、top_p 等）。
type ModelSettings = agent.ModelSettings

// ToolChoice 定义工具选择策略。
type ToolChoice = agent.ToolChoice

// ToolChoiceString 是字符串形式的工具选择。
type ToolChoiceString = agent.ToolChoiceString

// Truncation 定义消息截断策略。
type Truncation = agent.Truncation

// ========== Guardrails ==========

// GuardrailFunctionOutput 是护栏函数的输出。
type GuardrailFunctionOutput = agent.GuardrailFunctionOutput

// InputGuardrail 在处理用户输入前运行检查。
type InputGuardrail = agent.InputGuardrail

// InputGuardrailFunc 是输入护栏的函数签名。
type InputGuardrailFunc = agent.InputGuardrailFunc

// NewInputGuardrail 创建一个输入护栏。
func NewInputGuardrail(name string, fn InputGuardrailFunc) InputGuardrail {
	return agent.NewInputGuardrail(name, fn)
}

// InputGuardrailResult 是输入护栏的执行结果。
type InputGuardrailResult = agent.InputGuardrailResult

// OutputGuardrail 在生成最终输出后运行检查。
type OutputGuardrail = agent.OutputGuardrail

// OutputGuardrailFunc 是输出护栏的函数签名。
type OutputGuardrailFunc = agent.OutputGuardrailFunc

// NewOutputGuardrail 创建一个输出护栏。
func NewOutputGuardrail(name string, fn OutputGuardrailFunc) OutputGuardrail {
	return agent.NewOutputGuardrail(name, fn)
}

// OutputGuardrailResult 是输出护栏的执行结果。
type OutputGuardrailResult = agent.OutputGuardrailResult

// ========== OutputType ==========

// OutputTypeInterface 定义期望的输出格式。
type OutputTypeInterface = agent.OutputTypeInterface

// ========== MCP Config ==========

// MCPConfig 提供 MCP 服务器配置。
type MCPConfig = agent.MCPConfig

// PromptConfig 包含提示生成设置。
type PromptConfig = agent.PromptConfig

// ========== Runner ==========

// Runner 执行 Agent 的主循环。
type Runner = runner.Runner

// RunConfig 配置 Agent 执行行为。
type RunConfig = runner.RunConfig

// Run 使用默认 Runner 执行 Agent。
var Run = runner.Run

// RunResult 包含 Agent 执行的完整结果。
type RunResult = runner.RunResult

// RunItem 代表 Agent 执行过程中生成的一项内容。
type RunItem = runner.RunItem

// WrapRunItem 将 ResponseInputItemUnionParam 包装为 RunItem。
var WrapRunItem = runner.WrapRunItem

// Output 代表 Agent 的最终输出。
type Output = runner.Output

// Usage 跟踪 LLM 请求的 token 消耗。
type Usage = runner.Usage

// ModelResponse 包含单个 LLM 响应。
type ModelResponse = runner.ModelResponse

// MaxTurnsExceededError 在执行超过 MaxTurns 时返回。
type MaxTurnsExceededError = runner.MaxTurnsExceededError

// GuardrailTripwireTriggeredError 在护栏阻止执行时返回。
type GuardrailTripwireTriggeredError = runner.GuardrailTripwireTriggeredError

// DefaultMaxTurns 是默认的最大执行轮次。
const DefaultMaxTurns = runner.DefaultMaxTurns

// ========== Tool ==========

// Tool 定义 Agent 可调用的工具接口。
type Tool = tool.Tool

// FunctionTool 是函数形式的工具实现。
type FunctionTool = tool.FunctionTool

// ToolErrorFunction 处理工具执行错误。
type ToolErrorFunction = tool.ToolErrorFunction

// DefaultToolErrorFunction 是默认的工具错误处理函数。
var DefaultToolErrorFunction = tool.DefaultToolErrorFunction

// FunctionToolEnabler 定义工具是否启用的检查接口。
type FunctionToolEnabler = tool.FunctionToolEnabler

// ========== Tool Router ==========

// ToolRouter 动态选择相关工具。
type ToolRouter = tool.ToolRouter

// KeywordRouter 基于关键词匹配路由工具。
type KeywordRouter = tool.KeywordRouter

// ========== MCP ==========

// MCPServer 定义 Model Context Protocol 服务器接口。
type MCPServer = tool.MCPServer

// MCPToolFilter 过滤 MCP 工具。
type MCPToolFilter = tool.MCPToolFilter

// MCPToolFilterStatic 静态工具名过滤器。
type MCPToolFilterStatic = tool.MCPToolFilterStatic

// NewMCPToolFilterStatic 创建静态 MCP 工具过滤器。
var NewMCPToolFilterStatic = tool.NewMCPToolFilterStatic

// MCPToolFilterContext 提供过滤上下文。
type MCPToolFilterContext = tool.MCPToolFilterContext

// ToFunctionTool 将 MCP 工具转换为 FunctionTool。
var ToFunctionTool = tool.ToFunctionTool

// InvokeMCPTool 调用 MCP 工具。
var InvokeMCPTool = tool.InvokeMCPTool

// GetFunctionTools 从 MCP 服务器获取工具列表。
var GetFunctionTools = tool.GetFunctionTools

// GetAllFunctionTools 从多个 MCP 服务器获取所有工具。
var GetAllFunctionTools = tool.GetAllFunctionTools

// ApplyMCPToolFilter 应用工具过滤器。
var ApplyMCPToolFilter = tool.ApplyMCPToolFilter

// ========== Pattern: Agent-as-Tool ==========

// WrapAgentAsTool 将子 Agent 包装为 FunctionTool。
var WrapAgentAsTool = pattern.WrapAgentAsTool

// ========== Pattern: ReAct ==========

// DefaultReActInstruction 是标准的 ReAct 提示模板。
var DefaultReActInstruction = pattern.DefaultReActInstruction

// NewReActInstruction 创建自定义 ReAct 指令。
var NewReActInstruction = pattern.NewReActInstruction

// ReActStateProvider 跟踪 ReAct 执行状态。
type ReActStateProvider = pattern.ReActStateProvider

// NewReActStateProvider 创建新的 ReAct 状态提供者。
var NewReActStateProvider = pattern.NewReActStateProvider

// ========== Types ==========

// AgentLike 是跨包使用的最小 Agent 接口。
type AgentLike = types.AgentLike

// Input 代表 Agent 的输入。
type Input = types.Input

// InputString 是简单文本输入。
type InputString = types.InputString

// InputItems 是结构化输入项列表。
type InputItems = types.InputItems

// CopyInput 复制 Input 实例。
var CopyInput = types.CopyInput

  // ========== Config ==========

  // ModelConfig 存储单个模型的完整配置信息
  type ModelConfig = config.ModelConfig

  // ModelRegistry 管理多个模型配置
  type ModelRegistry = config.ModelRegistry

  // NewModelRegistry 创建新的模型注册表
  func NewModelRegistry() *ModelRegistry {
        return config.NewModelRegistry()
  }

  // LoadFromFile 从 JSON 文件加载多个模型配置
  var LoadFromFile = config.LoadFromFile

  // SaveToFile 将注册表中的所有配置保存到 JSON 文件
  var SaveToFile = config.SaveToFile

  // LoadOrCreate 尝试从文件加载配置，如果文件不存在则创建新注册表
  var LoadOrCreate = config.LoadOrCreate

// ========== Memory/Session ==========

// Session 管理对话历史。
type Session = memory.Session

// SQLiteSessionConfig 配置 SQLite 会话。
type SQLiteSessionConfig = memory.SQLiteSessionConfig

// SQLiteSession 是基于 SQLite 的会话实现。
type SQLiteSession = memory.SQLiteSession

// NewSQLiteSession 创建新的 SQLite 会话。
var NewSQLiteSession = memory.NewSQLiteSession

