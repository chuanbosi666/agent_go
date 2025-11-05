package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformStringFunctionStyle(t *testing.T) {
	input := "Foo Bar 123?Baz Quux!"
	result := TransformStringFunctionStyle(input)
	assert.Equal(t, "foo_bar_123_baz_quux_", result)
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"snake_case", "snakeCase"},
		{"PascalCase", "pascalCase"},
		{"camelCase", "camelCase"},
		{"", ""},
		{"single", "single"},
		{"UPPER_CASE", "upperCase"},
		{"multiple_word_example", "multipleWordExample"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"camelCase", "camel_case"},
		{"PascalCase", "pascal_case"},
		{"snake_case", "snake_case"},
		{"", ""},
		{"single", "single"},
		{"HTTPRequest", "httprequest"}, // Note: consecutive capitals
		{"getHTTPResponse", "get_httpresponse"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
