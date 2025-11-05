package nvgo

import (
	"context"
)

// GuardrailFunctionOutput represents the output of a guardrail function.
type GuardrailFunctionOutput struct {
	// TripwireTriggered indicates whether the guardrail's tripwire was activated.
	// If true, the agent's execution will be halted.
	TripwireTriggered bool

	// OutputInfo contains optional information about the guardrail's results.
	// For example, it could include details about performed checks.
	OutputInfo any
}

// InputGuardrail represents a check that runs in parallel with agent execution.
// It can validate input messages or take control if unexpected input is detected.
type InputGuardrail struct {
	// Name identifies the guardrail for tracing purposes.
	Name string

	// GuardrailFunc processes the input and returns the guardrail result.
	GuardrailFunc func(ctx context.Context, agent *Agent, input Input) (GuardrailFunctionOutput, error)
}

// NewInputGuardrail creates a new InputGuardrail with the specified name and guardrail function.
func NewInputGuardrail(name string, guardrailFunc func(ctx context.Context, agent *Agent, input Input) (GuardrailFunctionOutput, error)) InputGuardrail {
	return InputGuardrail{
		Name:          name,
		GuardrailFunc: guardrailFunc,
	}
}

// InputGuardrailResult contains the results of an input guardrail execution.
type InputGuardrailResult struct {
	// Guardrail is the guardrail that was executed.
	Guardrail InputGuardrail

	// Output contains the guardrail function's results.
	Output GuardrailFunctionOutput
}

// Run executes the input guardrail with the provided context, agent, and input.
func (g InputGuardrail) Run(ctx context.Context, agent *Agent, input Input) (InputGuardrailResult, error) {
	output, err := g.GuardrailFunc(ctx, agent, input)
	if err != nil {
		return InputGuardrailResult{}, err
	}
	return InputGuardrailResult{
		Guardrail: g,
		Output:    output,
	}, nil
}

// OutputGuardrail represents a check that validates the final output of an agent.
type OutputGuardrail struct {
	// Name identifies the guardrail for tracing purposes.
	Name string

	// GuardrailFunc processes the agent output and returns the guardrail result.
	GuardrailFunc func(ctx context.Context, agent *Agent, output any) (GuardrailFunctionOutput, error)
}

// NewOutputGuardrail creates a new OutputGuardrail with the specified name and guardrail function.
func NewOutputGuardrail(name string, guardrailFunc func(ctx context.Context, agent *Agent, output any) (GuardrailFunctionOutput, error)) OutputGuardrail {
	return OutputGuardrail{
		Name:          name,
		GuardrailFunc: guardrailFunc,
	}
}

// OutputGuardrailResult contains the results of an output guardrail execution.
type OutputGuardrailResult struct {
	// Guardrail is the guardrail that was executed.
	Guardrail OutputGuardrail

	// Agent is the agent whose output was checked.
	Agent *Agent

	// AgentOutput contains the agent's output that was validated.
	AgentOutput any

	// Output contains the guardrail function's results.
	Output GuardrailFunctionOutput
}

// Run executes the output guardrail with the provided context, agent, and output.
func (g OutputGuardrail) Run(ctx context.Context, agent *Agent, output any) (OutputGuardrailResult, error) {
	result, err := g.GuardrailFunc(ctx, agent, output)
	if err != nil {
		return OutputGuardrailResult{}, err
	}
	return OutputGuardrailResult{
		Guardrail:   g,
		Agent:       agent,
		AgentOutput: output,
		Output:      result,
	}, nil
}
