package nvgo

import (
	"fmt"
	"slices"

	"github.com/openai/openai-go/v2/responses"
)

// Input can be either a string or a list of TResponseInputItem.
type Input interface {
	isInput()
}

type InputString string

func (InputString) isInput()         {}
func (s InputString) String() string { return string(s) }

type InputItems []responses.ResponseInputItemUnionParam

func (InputItems) isInput() {}

func (items InputItems) Copy() InputItems {
	return slices.Clone(items)
}

func CopyInput(input Input) Input {
	switch v := input.(type) {
	case InputString:
		return v
	case InputItems:
		return v.Copy()
	default:
		// This would be an unrecoverable implementation bug, so a panic is appropriate.
		panic(fmt.Errorf("unexpected Input type %T", v))
	}
}
