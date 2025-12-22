package agent
import(
	"context"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

// Prompt defines configuration for OpenAI Responses API prompt management.
//This feature is only available with official OpenAI endpoints.
type Prompt struct{
	// ID is the unique identifier of the prompt.
	ID string
	// Version is the optional prompt version.
	Version param.Opt[string]
	// Variables are optional substitution values for the prompt template.
	Variables map[string]responses.ResponsePromptVariableUnionParam
}

// Prompt satisfies Prompter interface, returning itself.
func (p Prompt) Prompt(context.Context, *Agent) (Prompt, error) {
    return p, nil
}

// Prompter generates prompts dynamically.
type Prompter interface {
    Prompt(context.Context, *Agent) (Prompt, error)
}

// DynamicPromptFunction is a function that dynamically generates prompts.
type DynamicPromptFunction func(context.Context, *Agent) (Prompt, error)

// Prompt satisfies Prompter interface.
func (f DynamicPromptFunction) Prompt(ctx context.Context, agent *Agent) (Prompt, error) {
    return f(ctx, agent)
}

type promptUtil struct{}

// PromptUtil returns prompt utility functions.
func PromptUtil() promptUtil { 
	return promptUtil{} 
}

// ToModelInput converts Prompter to API request parameter.
func (promptUtil) ToModelInput(
        ctx context.Context,
        prompter Prompter,
        agent *Agent,
) (responses.ResponsePromptParam, bool, error) {
        if prompter == nil {
                return responses.ResponsePromptParam{}, false, nil
        }

        prompt, err := prompter.Prompt(ctx, agent)
        if err != nil {
                return responses.ResponsePromptParam{}, false, err
        }
        return responses.ResponsePromptParam{
                ID:        prompt.ID,
                Version:   prompt.Version,
                Variables: prompt.Variables,
        }, true, nil
}