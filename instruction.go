package nvgo

import (
	"context"
)

// InstructionsGetter interface is implemented by objects that can provide instructions to an Agent.
type InstructionsGetter interface {
	GetInstructions(context.Context, *Agent) (string, error)
}

// InstructionsStr satisfies InstructionsGetter providing a simple constant string value.
type InstructionsStr string

// GetInstructions returns the string value and always nil error.
func (s InstructionsStr) GetInstructions(context.Context, *Agent) (string, error) {
	return s.String(), nil
}

func (s InstructionsStr) String() string {
	return string(s)
}

// InstructionsFunc lets you implement a function that dynamically generates instructions for an Agent.
type InstructionsFunc func(context.Context, *Agent) (string, error)

// GetInstructions returns the string value and always nil error.
func (fn InstructionsFunc) GetInstructions(ctx context.Context, a *Agent) (string, error) {
	return fn(ctx, a)
}
