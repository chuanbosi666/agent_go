package memory

import (
	"context"

	"github.com/openai/openai-go/v2/responses"
)

// A Session stores conversation history for a specific session, allowing
// agents to maintain context without requiring explicit manual memory management.
type Session interface {
	// GetItems retrieves the conversation history for this session.
	//
	// `limit` is the maximum number of items to retrieve. If <= 0, retrieves all items.
	// When specified, returns the latest N items in chronological order.
	GetItems(ctx context.Context, limit int) ([]responses.ResponseInputItemUnionParam, error)

	// AddItems adds new items to the conversation history.
	AddItems(ctx context.Context, items []responses.ResponseInputItemUnionParam) error

	// PopItem removes and returns the most recent item from the session.
	// It returns nil if the session is empty.
	PopItem(context.Context) (*responses.ResponseInputItemUnionParam, error)

	// ClearSession clears all items for this session.
	ClearSession(context.Context) error
}
