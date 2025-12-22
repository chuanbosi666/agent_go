// Package agent - output.go defines the output type interface for structured outputs.
package agent

import "context"

// OutputTypeInterface describes the expected output format from an Agent.
// It supports both plain text and structured JSON outputs.
type OutputTypeInterface interface {
    // IsPlainText returns true if output is plain text (not JSON).
    IsPlainText() bool

    // Name returns the output type name.
    Name() string

    // JSONSchema returns the JSON schema for structured output.
    // Only called if IsPlainText() returns false.
    JSONSchema() (map[string]any, error)

    // IsStrictJSONSchema returns true if using strict JSON schema mode.
    // Strict mode guarantees valid JSON but has feature constraints.
    // See: https://platform.openai.com/docs/guides/structured-outputs
    IsStrictJSONSchema() bool

    // ValidateJSON validates and parses a JSON string against the output type.
    // Returns the parsed object or error if invalid.
    // Only called if IsPlainText() returns false.
    ValidateJSON(ctx context.Context, jsonStr string) (any, error)
  }