package runner

import "errors"

var (
	// ErrMaxTurnsExceeded indicates the agent exceeded max turns.
	ErrMaxTurnsExceeded = errors.New("max turns exceeded")

	// ErrGuardrailTriggered indicates a guardrail blocked execution.
	ErrGuardrailTriggered = errors.New("guardrail triggered")

	// ErrNoToolFound indicates the requested tool was not found.
	ErrNoToolFound = errors.New("tool not found")

	// ErrInvalidInput indicates invalid input was provided.
	ErrInvalidInput = errors.New("invalid input")
)
