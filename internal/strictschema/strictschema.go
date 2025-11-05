package strictschema

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// EnsureStrictJSONSchema mutates the given JSON schema to ensure it conforms
// to the `strict` standard that the OpenAI API expects.
func EnsureStrictJSONSchema(schema map[string]any) (map[string]any, error) {
	if len(schema) == 0 {
		return map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"properties":           map[string]any{},
			"required":             []string{},
		}, nil
	}
	return ensureStrictJSONSchema(schema, nil, schema)
}

func ensureStrictJSONSchema(schema any, path []string, root map[string]any) (map[string]any, error) {
	js, ok := schema.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected map[string]any at path %q, got %T", strings.Join(path, "/"), schema)
	}

	for _, defKey := range []string{"$defs", "definitions"} {
		if defs, ok := js[defKey].(map[string]any); ok {
			for name, def := range defs {
				_, err := ensureStrictJSONSchema(def, append(slices.Clone(path), defKey, name), root)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if typ, _ := js["type"].(string); typ == "object" {
		ap, ok := js["additionalProperties"]
		if !ok {
			js["additionalProperties"] = false
		} else if ap != false && !maps.Equal(ap.(map[string]any), map[string]any{"not": map[string]any{}}) {
			return nil, errors.New("additionalProperties must be false for object types in strict schemas")
		}
	}

	if props, ok := js["properties"].(map[string]any); ok {
		keys := slices.Collect(maps.Keys(props))
		sort.Strings(keys)
		js["required"] = keys

		newProps := make(map[string]any, len(props))
		for _, key := range keys {
			var err error
			newProps[key], err = ensureStrictJSONSchema(props[key], append(slices.Clone(path), "properties", key), root)
			if err != nil {
				return nil, err
			}
		}
		js["properties"] = newProps
	}

	if items, ok := js["items"].(map[string]any); ok {
		var err error
		js["items"], err = ensureStrictJSONSchema(items, append(slices.Clone(path), "items"), root)
		if err != nil {
			return nil, err
		}
	}

	if anyOf, ok := js["anyOf"].([]any); ok {
		newAnyOf := make([]any, len(anyOf))
		for i, v := range anyOf {
			var err error
			newAnyOf[i], err = ensureStrictJSONSchema(v, append(slices.Clone(path), "anyOf", strconv.Itoa(i)), root)
			if err != nil {
				return nil, err
			}
		}
		js["anyOf"] = newAnyOf
	}

	if allOf, ok := js["allOf"].([]any); ok {
		if len(allOf) == 1 {
			res, err := ensureStrictJSONSchema(allOf[0], append(slices.Clone(path), "allOf", "0"), root)
			if err != nil {
				return nil, err
			}
			delete(js, "allOf")
			maps.Copy(js, res)
		} else {
			newAllOf := make([]any, len(allOf))
			for i, v := range allOf {
				var err error
				newAllOf[i], err = ensureStrictJSONSchema(v, append(slices.Clone(path), "allOf", strconv.Itoa(i)), root)
				if err != nil {
					return nil, err
				}
			}
			js["allOf"] = newAllOf
		}
	}

	if d, ok := js["default"]; ok && d == nil {
		delete(js, "default")
	}

	if ref, ok := js["$ref"]; ok && len(js) > 1 {
		refStr, ok := ref.(string)
		if !ok {
			return nil, fmt.Errorf("non-string $ref at path %q: got %T", strings.Join(path, "/"), ref)
		}
		resolved, err := resolveJSONSchemaRef(root, refStr)
		if err != nil {
			return nil, err
		}
		delete(js, "$ref")
		for k, v := range resolved {
			if _, exists := js[k]; !exists {
				js[k] = v
			}
		}
		return ensureStrictJSONSchema(js, path, root)
	}

	return js, nil
}

func resolveJSONSchemaRef(root map[string]any, ref string) (map[string]any, error) {
	if !strings.HasPrefix(ref, "#/") {
		return nil, fmt.Errorf("unexpected $ref format: expected `#/` prefix in $ref value %q", ref)
	}

	path := strings.Split(ref[2:], "/")
	resolved := root

	for _, key := range path {
		entry, ok := resolved[key]
		if !ok {
			return nil, fmt.Errorf("missing key %q while resolving $ref %q", key, ref)
		}
		resolved, ok = entry.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("non-object entry at key %q while resolving $ref %q: got %T", key, ref, entry)
		}
	}

	return resolved, nil
}
