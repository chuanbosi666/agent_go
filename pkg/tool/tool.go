package tool

import(
	"context"
	"fmt"

	"github.com/openai/openai-go/v3/packages/param"
)

// Tool is the interface that all tools must implement.
type Tool interface {
	ToolName() string
	isTool()
	GetName() string
	GetDescription() string
	GetParamsJSONSchema() map[string]any
	Invoke(ctx context.Context, args string) (any, error)
}

// ToolErrorFunction handles tool invocation errors and returns a value to send back to the LLM.
type ToolErrorFunction func(ctx context.Context, err error) (any, error)

// DefaultToolErrorFunction is used when FunctionTool doesn't specify its own error handler.
func DefaultToolErrorFunction(_ context.Context, err error) (any, error) {
	return fmt.Sprintf("An error occurred while running the tool. Please try again. Error: %s", err), nil    
}

// FunctionToolEnabler determines if a tool should be enabled.
type FunctionToolEnabler interface {
	IsEnabled(ctx context.Context, agent any) (bool, error)
}

var _ Tool = (*FunctionTool)(nil)

// FunctionTool is a Tool that wraps a callable function.
type FunctionTool struct {
	// Name is the tool name shown to the LLM.
	Name string

	// Description explains what the tool does (shown to LLM).
	Description string

	// ParamsJSONSchema defines the JSON schema for tool parameters.
	ParamsJSONSchema map[string]any

	// OnInvokeTool is called when the LLM invokes this tool.
	// Arguments is the JSON string from the LLM.
	OnInvokeTool func(ctx context.Context, arguments string) (any, error)

	// FailureErrorFunction handles errors (optional).
	// If nil, DefaultToolErrorFunction is used.
	FailureErrorFunction *ToolErrorFunction

	// StrictJSONSchema enables strict JSON schema mode.
	// When true, LLM guarantees valid JSON output.
	StrictJSONSchema param.Opt[bool]

	// IsEnabled optionally controls whether the tool is available.
	IsEnabled FunctionToolEnabler

}

// ToolName returns the tool's name.
func (t FunctionTool) ToolName() string {
	return t.Name
}

func (t FunctionTool) isTool() {}

func (t FunctionTool) GetName() string {
        return t.Name
  }

  // GetDescription returns the tool's description.
  func (t FunctionTool) GetDescription() string {
        return t.Description
  }

  // GetParamsJSONSchema returns the tool's parameter schema.
  func (t FunctionTool) GetParamsJSONSchema() map[string]any {
        return t.ParamsJSONSchema
  }

  // Invoke executes the tool with the given arguments.
  func (t FunctionTool) Invoke(ctx context.Context, args string) (any, error) {
        if t.OnInvokeTool == nil {
                return nil, fmt.Errorf("tool %q has no implementation", t.Name)
        }

        result, err := t.OnInvokeTool(ctx, args)
        if err != nil {
                // Use custom error handler if provided
                if t.FailureErrorFunction != nil {
                        return (*t.FailureErrorFunction)(ctx, err)
                }
                // Otherwise use default error handler
                return DefaultToolErrorFunction(ctx, err)
        }

        return result, nil
  }