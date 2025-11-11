package strictschema_test

import (
	"fmt"
	"testing"

	"github.com/agent_go/internal/strictschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptySchemaHasAdditionalPropertiesFalse(t *testing.T) {
	t.Parallel()
	schemas := []map[string]any{
		nil,
		{},
	}
	for _, schema := range schemas {
		t.Run(fmt.Sprintf("%#v", schema), func(t *testing.T) {
			strictSchema, err := strictschema.EnsureStrictJSONSchema(schema)
			require.NoError(t, err)
			assert.Contains(t, strictSchema, "additionalProperties")
			assert.Equal(t, false, strictSchema["additionalProperties"])
		})
	}
}

func TestObjectWithoutAdditionalProperties(t *testing.T) {
	t.Parallel()
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"a": map[string]any{"type": "string"},
		},
	}
	result, err := strictschema.EnsureStrictJSONSchema(schema)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"a"},
		"properties": map[string]any{
			"a": map[string]any{"type": "string"},
		},
	}, result)
}

func TestArrayItemsProcessingAndDefaultRemoval(t *testing.T) {
	t.Parallel()
	schema := map[string]any{
		"type": "array",
		"items": map[string]any{
			"type":    "number",
			"default": nil,
		},
	}
	result, err := strictschema.EnsureStrictJSONSchema(schema)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"type":  "array",
		"items": map[string]any{"type": "number"},
	}, result)
}

func TestAnyOfProcessing(t *testing.T) {
	t.Parallel()
	schema := map[string]any{
		"anyOf": []any{
			map[string]any{"type": "object", "properties": map[string]any{"a": map[string]any{"type": "string"}}},
			map[string]any{"type": "number", "default": nil},
		},
	}
	result, err := strictschema.EnsureStrictJSONSchema(schema)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"anyOf": []any{
			map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"required":             []string{"a"},
				"properties":           map[string]any{"a": map[string]any{"type": "string"}},
			},
			map[string]any{"type": "number"},
		},
	}, result)
}

func TestAllOfSingleEntryMerging(t *testing.T) {
	t.Parallel()
	schema := map[string]any{
		"type": "object",
		"allOf": []any{
			map[string]any{"properties": map[string]any{"a": map[string]any{"type": "boolean"}}},
		},
	}
	result, err := strictschema.EnsureStrictJSONSchema(schema)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"a"},
		"properties":           map[string]any{"a": map[string]any{"type": "boolean"}},
	}, result)
}

func TestDefaultRemovalOnNonObject(t *testing.T) {
	t.Parallel()
	schema := map[string]any{"type": "string", "default": nil}
	result, err := strictschema.EnsureStrictJSONSchema(schema)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"type": "string"}, result)
}

func TestRefExpansion(t *testing.T) {
	t.Parallel()
	schema := map[string]any{
		"definitions": map[string]any{"refObj": map[string]any{"type": "string", "default": nil}},
		"type":        "object",
		"properties":  map[string]any{"a": map[string]any{"$ref": "#/definitions/refObj", "description": "desc"}},
	}
	result, err := strictschema.EnsureStrictJSONSchema(schema)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"definitions":          map[string]any{"refObj": map[string]any{"type": "string"}},
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"a"},
		"properties":           map[string]any{"a": map[string]any{"type": "string", "description": "desc"}},
	}, result)
}

func TestRefNoExpansionWhenAlone(t *testing.T) {
	t.Parallel()
	schema := map[string]any{"$ref": "#/definitions/refObj"}
	result, err := strictschema.EnsureStrictJSONSchema(schema)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"$ref": "#/definitions/refObj"}, result)
}

func TestInvalidRefFormat(t *testing.T) {
	t.Parallel()
	schema := map[string]any{
		"type":       "object",
		"properties": map[string]any{"a": map[string]any{"$ref": "invalid", "description": "desc"}},
	}
	_, err := strictschema.EnsureStrictJSONSchema(schema)
	assert.Error(t, err)
}
