package nvgo

import (
	"github.com/openai/openai-go/v3"
)

// MCPConfig provides configuration parameters for MCP servers.
type MCPConfig struct {
	// If true, we will attempt to convert the MCP schemas to strict-mode schemas.
	// This is a best-effort conversion, so some schemas may not be convertible.
	// Defaults to false.
	ConvertSchemasToStrict bool
}

// An Agent is an AI model configured with instructions, tools, guardrails, handoffs and more.
//
// We strongly recommend passing `Instructions`, which is the "system prompt" for the agent. In
// addition, you can pass `HandoffDescription`, which is a human-readable description of the
// agent, used when the agent is used inside tools/handoffs.
type Agent struct {
	// The name of the agent.
	Name string

	// Optional instructions for the agent. Will be used as the "system prompt" when this agent is
	// invoked. Describes what the agent should do, and how it responds.
	Instructions InstructionsGetter

	// Optional Prompter object. Prompts allow you to dynamically configure the instructions,
	// tools and other config for an agent outside your code.
	// Only usable with OpenAI models, using the Responses API.
	Prompt Prompter

	// The model implementation to use when invoking the LLM.
	Model string

	// The client to use when invoking the LLM.
	Client openai.Client

	// Configures model-specific tuning parameters (e.g. temperature, top_p).
	ModelSettings ModelSettings

	// Optional list of Model Context Protocol (https://modelcontextprotocol.io) servers that
	// the agent can use. Every time the agent runs, it will include tools from these servers in the
	// list of available tools.
	//
	// NOTE: You are expected to manage the lifecycle of these servers. Specifically, you must call
	// `MCPServer.Connect()` before passing it to the agent, and `MCPServer.Cleanup()` when the server is no
	// longer needed.
	MCPServers []MCPServer

	// Optional configuration for MCP servers.
	MCPConfig MCPConfig

	// A list of checks that run in parallel to the agent's execution, before generating a
	// response. Runs only if the agent is the first agent in the chain.
	InputGuardrails []InputGuardrail

	// A list of checks that run on the final output of the agent, after generating a response.
	// Runs only if the agent produces a final output.
	OutputGuardrails []OutputGuardrail

	Tools []FunctionTool

	// Optional output type describing the output. If not provided, the output will be a simple string.
	OutputType OutputTypeInterface
}

// New creates a new Agent with the given name.
//
// The returned Agent can be further configured using the builder methods.
func New(name string) *Agent {
	return &Agent{Name: name}
}

// WithInstructions sets the Agent instructions.
func (a *Agent) WithInstructions(instr string) *Agent {
	a.Instructions = InstructionsStr(instr)
	return a
}

// WithInstructionsFunc sets dynamic instructions using an InstructionsFunc.
func (a *Agent) WithInstructionsFunc(fn InstructionsFunc) *Agent {
	a.Instructions = fn
	return a
}

// WithInstructionsGetter sets custom instructions implementing InstructionsGetter.
func (a *Agent) WithInstructionsGetter(g InstructionsGetter) *Agent {
	a.Instructions = g
	return a
}

// WithPrompt sets the agent's static or dynamic prompt.
func (a *Agent) WithPrompt(prompt Prompter) *Agent {
	a.Prompt = prompt
	return a
}

// WithModel sets the model to use by name.
func (a *Agent) WithModel(model string) *Agent {
	a.Model = model
	return a
}

// WithClient sets the client to use.
// use openai
func (a *Agent) WithClient(client openai.Client) *Agent {
	a.Client = client
	return a
}

// WithModelSettings sets model-specific settings.
func (a *Agent) WithModelSettings(settings ModelSettings) *Agent {
	a.ModelSettings = settings
	return a
}

// WithMCPServers sets the list of MCP servers available to the agent.
func (a *Agent) WithMCPServers(mcpServers []MCPServer) *Agent {
	a.MCPServers = mcpServers
	return a
}

// AddMCPServer appends an MCP server to the agent's MCP server list.
func (a *Agent) AddMCPServer(mcpServer MCPServer) *Agent {
	a.MCPServers = append(a.MCPServers, mcpServer)
	return a
}

// Add MCPStdioServer appends an MCP server to the agent's MCP server list.
// func (a *Agent) AddMCPStdioServer() *Agent {
// 	return
// }

// WithMCPConfig sets the agent's MCP configuration.
func (a *Agent) WithMCPConfig(mcpConfig MCPConfig) *Agent {
	a.MCPConfig = mcpConfig
	return a
}

// WithInputGuardrails sets the input guardrails.
func (a *Agent) WithInputGuardrails(gr []InputGuardrail) *Agent {
	a.InputGuardrails = gr
	return a
}

// AddInputGuardrail appends an input guardrail to the agent's input guardrails list.
func (a *Agent) AddInputGuardrail(gr InputGuardrail) *Agent {
	a.InputGuardrails = append(a.InputGuardrails, gr)
	return a
}

// WithOutputGuardrails sets the output guardrails.
func (a *Agent) WithOutputGuardrails(gr []OutputGuardrail) *Agent {
	a.OutputGuardrails = gr
	return a
}

// AddOutputGuardrail appends an output guardrail to the agent's output guardrails list.
func (a *Agent) AddOutputGuardrail(gr OutputGuardrail) *Agent {
	a.OutputGuardrails = append(a.OutputGuardrails, gr)
	return a
}

// WithOutputType sets the output type.
func (a *Agent) WithOutputType(outputType OutputTypeInterface) *Agent {
	a.OutputType = outputType
	return a
}

func (a *Agent) WithTools(tools []FunctionTool) *Agent {
	a.Tools = tools
	return a
}

func (a *Agent) AddTools(tools []FunctionTool) *Agent {
	a.Tools = append(a.Tools, tools...)
	return a
}
