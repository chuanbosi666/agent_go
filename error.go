package nvgo

import "errors"

// MCP errors
var (
	// ErrMCPServerNotInitialized indicates that the server is not initialized.
	ErrMCPServerNotInitialized = errors.New("server not initialized: make sure you call `Connect()` first")
	// ErrMCPAgentRequired indicates that the agent is required for dynamic tool filtering.
	ErrMCPAgentRequired = errors.New("agent is required for dynamic tool filtering")
)
