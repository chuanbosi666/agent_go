package nvgo

import "context"

// OutputTypeInterface is implemented by an object that describes an agent's output type.
// Unless the output type is plain text (string), it captures the JSON schema of the output,
// as well as validating/parsing JSON produced by the LLM into the output type.
type OutputTypeInterface interface {
	// IsPlainText reports whether the output type is plain text (versus a JSON object).
	IsPlainText() bool

	// The Name of the output type.
	Name() string

	// JSONSchema returns the JSON schema of the output.
	// It will only be called if the output type is not plain text.
	JSONSchema() (map[string]any, error)

	// IsStrictJSONSchema reports whether the JSON schema is in strict mode.
	// Strict mode constrains the JSON schema features, but guarantees valid JSON.
	//
	// For more details, see https://platform.openai.com/docs/guides/structured-outputs#supported-schemas
	IsStrictJSONSchema() bool

	// ValidateJSON validates a JSON string against the output type.
	// You must return the validated object, or a `ModelBehaviorError` if the JSON is invalid.
	// It will only be called if the output type is not plain text.
	ValidateJSON(ctx context.Context, jsonStr string) (any, error)
}
