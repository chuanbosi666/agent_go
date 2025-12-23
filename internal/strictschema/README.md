# StrictSchema

A Go library to convert JSON Schemas into a strict format compatible with OpenAI's API expectations for tool calls. It enforces `additionalProperties: false` on objects, sets `required` from property keys, resolves `$ref`s when needed, removes `null` defaults, and handles unions/intersections recursively.

## Usage

```go
import "github.com/demo/github.com/chuanbosi666/agent_go/internal/strictschema"

schema := map[string]any{ /* your JSON schema */ }
strict, err := strictschema.EnsureStrictJSONSchema(schema)
if err != nil {
    // handle error
}
// Use strict in OpenAI tool definitions
```

For empty schemas:

```go
strictschema.NewEmptyJSONSchema() // Returns a base strict object schema
```

## Examples

### Empty Schema

Before: `{}`
After:

```json
{
  "type": "object",
  "additionalProperties": false,
  "properties": {},
  "required": []
}
```

### Object without `additionalProperties`

Before:

```json
{
  "type": "object",
  "properties": {
    "a": { "type": "string" }
  }
}
```

After:

```json
{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "a": { "type": "string" }
  },
  "required": ["a"]
}
```

### Array with Default Null

Before:

```json
{
  "type": "array",
  "items": {
    "type": "number",
    "default": null
  }
}
```

After:

```json
{
  "type": "array",
  "items": { "type": "number" }
}
```

### Union (anyOf)

Before:

```json
{
  "anyOf": [
    {
      "type": "object",
      "properties": { "a": { "type": "string" } }
    },
    { "type": "number" }
  ]
}
```

After:

```json
{
  "anyOf": [
    {
      "type": "object",
      "additionalProperties": false,
      "properties": { "a": { "type": "string" } },
      "required": ["a"]
    },
    { "type": "number" }
  ]
}
```

### Reference Expansion

Before:

```json
{
  "definitions": {
    "refObj": { "type": "string", "default": null }
  },
  "type": "object",
  "properties": {
    "a": { "$ref": "#/definitions/refObj", "description": "desc" }
  }
}
```

After:

```json
{
  "definitions": { "refObj": { "type": "string" } },
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "a": { "type": "string", "description": "desc" }
  },
  "required": ["a"]
}
```

## Links

- [OpenAI Function Calling Guide](https://platform.openai.com/docs/guides/function-calling)
- [OpenAI Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs#supported-schemas)
