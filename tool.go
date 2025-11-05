package nvgo

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v2/packages/param"
)

// A Tool that can be used in an Agent.
type Tool interface {
	// ToolName returns the name of the tool.
	ToolName() string
	isTool()
}

// ToolErrorFunction is a callback that handles tool invocation errors and returns a value to be sent back to the LLM.
// If this function returns an error, it will be treated as a fatal error for the tool.
type ToolErrorFunction func(ctx context.Context, err error) (any, error)

// DefaultToolErrorFunction is the default handler used when a FunctionTool does not specify its own FailureErrorFunction.
// It returns a generic error message containing the original error string.
func DefaultToolErrorFunction(_ context.Context, err error) (any, error) {
	return fmt.Sprintf("An error occurred while running the tool. Please try again. Error: %s", err), nil
}

type FunctionToolEnabler interface {
	IsEnabled(ctx context.Context, agent *Agent) (bool, error)
}

var _ Tool = (*FunctionTool)(nil)

// FunctionTool is a Tool that wraps a function.
type FunctionTool struct {
	// The name of the tool, as shown to the LLM. Generally the name of the function.
	Name string

	// A description of the tool, as shown to the LLM.
	Description string

	// The JSON schema for the tool's parameters.
	ParamsJSONSchema map[string]any

	// A function that invokes the tool with the given context and parameters.
	//
	// The params passed are:
	// 	1. The tool run context.
	// 	2. The arguments from the LLM, as a JSON string.
	//
	// You must return a string representation of the tool output.
	// In case of errors, you can either return an error (which will cause the run to fail) or
	// return a string error message (which will be sent back to the LLM).
	OnInvokeTool func(ctx context.Context, arguments string) (any, error)

	// Optional error handling function. When the tool invocation returns an error,
	// this function is called with the original error and its return value is sent
	// back to the LLM. If not set, a default function returning a generic error
	// message is used. To disable error handling and propagate the original error,
	// explicitly set this to a pointer to a nil ToolErrorFunction.
	FailureErrorFunction *ToolErrorFunction

	// Whether the JSON schema is in strict mode.
	// We **strongly** recommend setting this to True, as it increases the likelihood of correct JSON input.
	// Defaults to true if omitted.
	StrictJSONSchema param.Opt[bool]

	// Optional flag reporting whether the tool is enabled.
	// It can be either a boolean or a function which allows you to dynamically
	// enable/disable a tool based on your context/state.
	// Default value, if omitted: true.
	IsEnabled FunctionToolEnabler
}

func (t FunctionTool) ToolName() string {
	return t.Name
}

func (t FunctionTool) isTool() {}
