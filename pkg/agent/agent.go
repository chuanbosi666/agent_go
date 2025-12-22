package agent

import (
	"nvgo/pkg/tool"

	"github.com/openai/openai-go/v3"
)


// PromptConfig holds prompt generation settings.
type PromptConfig struct {
	StateProvider StateProvider
}


// MCPConfig provides configuration for MCP servers.
type MCPConfig struct {
	// ConvertSchemasToStrict attempts to convert MCP schemas to strict-mode (best-effort).
	ConvertSchemasToStrict bool
}

// Agent represents an AI model configured with instructions, tools, guardrails and more.
type Agent struct {
	// Name is the agent identifier.
	Name string

	// Instructions is the system prompt defining agent behavior.
	Instructions Instructions

	// Prompt enables dynamic configuration via OpenAI Responses API (OpenAI-only).
	// For OpenAI-compatible endpoints, use Chat Completions API instead.
	Prompt Prompter

	// Model specifies the LLM model name (e.g., "gpt-4o", "gpt-4o-mini").
	Model string

	// Client is the OpenAI client for API calls.
	Client openai.Client

	// ModelSettings contains model parameters (temperature, top_p, etc.).
	ModelSettings ModelSettings

	// MCPServers lists Model Context Protocol servers providing tools to the agent.
	// You must call Connect() before use and Cleanup() when done.
	MCPServers []tool.MCPServer

	// MCPConfig provides optional configuration for MCP servers.
	MCPConfig MCPConfig

	// InputGuardrails are checks run before generating responses (first agent only).
	InputGuardrails []InputGuardrail

	// OutputGuardrails are checks run on final outputs.
	OutputGuardrails []OutputGuardrail

	// Tools is the list of function tools available to this agent.
	Tools []tool.FunctionTool

	// OutputType describes the expected output format (defaults to plain text).
	OutputType OutputTypeInterface
}

// Implement types.AgentLike interface
func (a *Agent) GetName() string  { return a.Name }
func (a *Agent) GetModel() string { return a.Model }

// New creates a new Agent with the given name.
func New(name string) *Agent {
	return &Agent{Name: name}
}

// WithInstructions sets static instructions (system prompt).
func (a *Agent) WithInstructions(instr string) *Agent {
	a.Instructions = InstructionsStr(instr)
	return a
}

// WithInstructionsFunc sets dynamic instructions via function.
func (a *Agent) WithInstructionsFunc(fn InstructionsFunc) *Agent {
	a.Instructions = fn
	return a
}

// WithInstructionsGetter sets a custom InstructionsGetter.
func (a *Agent) WithInstructionsGetter(g Instructions) *Agent {
	a.Instructions = g
	return a
}

// WithPrompt sets the prompt configuration.
func (a *Agent) WithPrompt(prompt Prompter) *Agent {
	a.Prompt = prompt
	return a
}

// WithModel sets the LLM model name.
func (a *Agent) WithModel(model string) *Agent {
	a.Model = model
	return a
}

// WithClient sets the OpenAI client.
func (a *Agent) WithClient(client openai.Client) *Agent {
	a.Client = client
	return a
}

// WithModelSettings sets model parameters.
func (a *Agent) WithModelSettings(settings ModelSettings) *Agent {
	a.ModelSettings = settings
	return a
}

// WithMCPServers sets the MCP server list.
func (a *Agent) WithMCPServers(mcpServers []tool.MCPServer) *Agent {
	a.MCPServers = mcpServers
	return a
}

// AddMCPServer appends an MCP server.
func (a *Agent) AddMCPServer(mcpServer tool.MCPServer) *Agent {
	a.MCPServers = append(a.MCPServers, mcpServer)
	return a
}

// WithMCPConfig sets MCP configuration.
func (a *Agent) WithMCPConfig(mcpConfig MCPConfig) *Agent {
	a.MCPConfig = mcpConfig
	return a
}

// WithInputGuardrails sets input guardrails.
func (a *Agent) WithInputGuardrails(gr []InputGuardrail) *Agent {
	a.InputGuardrails = gr
	return a
}

// AddInputGuardrail appends an input guardrail.
func (a *Agent) AddInputGuardrail(gr InputGuardrail) *Agent {
	a.InputGuardrails = append(a.InputGuardrails, gr)
	return a
}

// WithOutputGuardrails sets output guardrails.
func (a *Agent) WithOutputGuardrails(gr []OutputGuardrail) *Agent {
	a.OutputGuardrails = gr
	return a
}

// AddOutputGuardrail appends an output guardrail.
func (a *Agent) AddOutputGuardrail(gr OutputGuardrail) *Agent {
	a.OutputGuardrails = append(a.OutputGuardrails, gr)
	return a
}

// WithOutputType sets the expected output type.
func (a *Agent) WithOutputType(outputType OutputTypeInterface) *Agent {
	a.OutputType = outputType
	return a
}

// WithTools sets the function tools.
func (a *Agent) WithTools(tools []tool.FunctionTool) *Agent {
	a.Tools = tools
	return a
}

// AddTools appends multiple tools.
func (a *Agent) AddTools(tools []tool.FunctionTool) *Agent {
	a.Tools = append(a.Tools, tools...)
	return a
}
