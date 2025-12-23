package agent

import (
	"context"

	"github.com/chuanbosi666/agent_go/pkg/types"
)

// GuardrailFunctionOutput represents the result of a guardrail check.
type GuardrailFunctionOutput struct {
	// TripwireTriggered halts agent execution when true.
	TripwireTriggered bool
	// OutputInfo contains optional check details.
	OutputInfo any
}

// InputGuardrailFunc validates input before processing.
type InputGuardrailFunc func(ctx context.Context, agent types.AgentLike, input types.Input) (GuardrailFunctionOutput, error)

// InputGuardrail validates input before agent processing.
type InputGuardrail struct {
	Name          string
	GuardrailFunc InputGuardrailFunc
}

// NewInputGuardrail creates an input guardrail.
func NewInputGuardrail(name string, fn InputGuardrailFunc) InputGuardrail {
	return InputGuardrail{Name: name, GuardrailFunc: fn}
}

// InputGuardrailResult contains input guardrail execution results.
type InputGuardrailResult struct {
	Guardrail InputGuardrail
	Output    GuardrailFunctionOutput
}

// Run executes the input guardrail.
func (g InputGuardrail) Run(ctx context.Context, agent types.AgentLike, input types.Input) (InputGuardrailResult, error) {
	output, err := g.GuardrailFunc(ctx, agent, input)
	if err != nil {
		return InputGuardrailResult{}, err
	}
	return InputGuardrailResult{Guardrail: g, Output: output}, nil
}

// OutputGuardrailFunc validates output before returning.
type OutputGuardrailFunc func(ctx context.Context, agent types.AgentLike, output any) (GuardrailFunctionOutput, error)

// OutputGuardrail validates agent output.
type OutputGuardrail struct {
	Name          string
	GuardrailFunc OutputGuardrailFunc
}

// NewOutputGuardrail creates an output guardrail.
func NewOutputGuardrail(name string, fn OutputGuardrailFunc) OutputGuardrail {
	return OutputGuardrail{Name: name, GuardrailFunc: fn}
}

// OutputGuardrailResult contains output guardrail execution results.
type OutputGuardrailResult struct {
	Guardrail   OutputGuardrail
	Agent       types.AgentLike
	AgentOutput any
	Output      GuardrailFunctionOutput
}

// Run executes the output guardrail.
func (g OutputGuardrail) Run(ctx context.Context, agent types.AgentLike, output any) (OutputGuardrailResult, error) {
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
