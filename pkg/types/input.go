package types

import (
	"fmt"
	"slices"

	"github.com/openai/openai-go/v3/responses"
)

// Input represents the input to an Agent. It can be either a simple string
// or a list of structured input items.
type Input interface {
	isInput()
	ToInputItems() []responses.ResponseInputItemUnionParam
}

// ItemsInput is a list of structured input items (backwards compatibility).
type ItemsInput []responses.ResponseInputItemUnionParam

func (i ItemsInput) isInput() {}

// ToInputItems implements the Input interface.
func (i ItemsInput) ToInputItems() []responses.ResponseInputItemUnionParam {
	return []responses.ResponseInputItemUnionParam(i)
}

// InputString is a simple text input.
type InputString string

func (InputString) isInput() {}

// ToInputItems converts the string to a message input item.
func (s InputString) ToInputItems() []responses.ResponseInputItemUnionParam {
	return []responses.ResponseInputItemUnionParam{
		responses.ResponseInputItemParamOfMessage(
			string(s),
			responses.EasyInputMessageRole(responses.ResponseInputMessageItemRoleUser)),
	}
}

// String returns the underlying string value.
func (s InputString) String() string { return string(s) }

// InputItems is a list of structured input items (messages, tool results, etc.).
type InputItems []responses.ResponseInputItemUnionParam

func (InputItems) isInput() {}

// ToInputItems implements the Input interface.
func (items InputItems) ToInputItems() []responses.ResponseInputItemUnionParam {
	return []responses.ResponseInputItemUnionParam(items)
}

// Copy returns a shallow copy of the input items.
func (items InputItems) Copy() InputItems {
	return slices.Clone(items)
}

// CopyInput creates a copy of the input.
// For InputString, returns the same value (strings are immutable).
// For InputItems, returns a cloned slice.
func CopyInput(input Input) Input {
	switch v := input.(type) {
	case InputString:
		return v
	case InputItems:
		return v.Copy()
	default:
		panic(fmt.Errorf("unexpected Input type %T", v))
	}
}
