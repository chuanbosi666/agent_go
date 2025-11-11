package nvgo_test

import (
	"context"
	"testing"

	"github.com/agent_go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMCPToolFilterStatic(t *testing.T) {
	tests := []struct {
		name     string
		allowed  []string
		blocked  []string
		expected nvgo.MCPToolFilter
		ok       bool
	}{
		{
			name:    "no lists",
			allowed: nil,
			blocked: nil,
			ok:      false,
		},
		{
			name:    "empty lists",
			allowed: []string{},
			blocked: []string{},
			ok:      false,
		},
		{
			name:    "allowed list",
			allowed: []string{"tool1"},
			ok:      true,
		},
		{
			name:    "blocked list",
			blocked: []string{"tool2"},
			ok:      true,
		},
		{
			name:    "both lists",
			allowed: []string{"tool1"},
			blocked: []string{"tool2"},
			ok:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, ok := nvgo.NewMCPToolFilterStatic(tt.allowed, tt.blocked)
			assert.Equal(t, tt.ok, ok)
			if tt.ok {
				assert.IsType(t, &nvgo.MCPToolFilterStatic{}, filter)
				actual := filter.(*nvgo.MCPToolFilterStatic)
				assert.Equal(t, tt.allowed, actual.AllowedToolNames)
				assert.Equal(t, tt.blocked, actual.BlockedToolNames)
			} else {
				assert.Nil(t, filter)
			}
		})
	}
}

func TestMCPToolFilterStatic_FilterMCPTool(t *testing.T) {
	tests := []struct {
		name     string
		allowed  []string
		blocked  []string
		toolName string
		expected bool
	}{
		{
			name:     "allowed tool",
			allowed:  []string{"allowed"},
			toolName: "allowed",
			expected: true,
		},
		{
			name:     "blocked tool",
			blocked:  []string{"blocked"},
			toolName: "blocked",
			expected: false,
		},
		{
			name:     "unknown tool with allowed list",
			allowed:  []string{"allowed"},
			toolName: "unknown",
			expected: false,
		},
		{
			name:     "unknown tool with blocked list",
			blocked:  []string{"blocked"},
			toolName: "unknown",
			expected: true,
		},
		{
			name:     "tool in both lists",
			allowed:  []string{"tool"},
			blocked:  []string{"tool"},
			toolName: "tool",
			expected: false, // blocked takes precedence implicitly
		},
		{
			name:     "no lists",
			toolName: "any",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := &nvgo.MCPToolFilterStatic{
				AllowedToolNames: tt.allowed,
				BlockedToolNames: tt.blocked,
			}
			tool := &mcp.Tool{Name: tt.toolName}
			include, err := filter.FilterMCPTool(context.Background(), nvgo.MCPToolFilterContext{}, tool)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, include)
		})
	}
}

func TestApplyMCPToolFilter(t *testing.T) {
	tools := []*mcp.Tool{
		{Name: "tool1"},
		{Name: "tool2"},
		{Name: "tool3"},
	}

	tests := []struct {
		name     string
		filter   nvgo.MCPToolFilter
		expected []string
	}{
		{
			name:     "no filter",
			filter:   nil,
			expected: []string{"tool1", "tool2", "tool3"},
		},
		{
			name:     "allow tool1",
			filter:   &nvgo.MCPToolFilterStatic{AllowedToolNames: []string{"tool1"}},
			expected: []string{"tool1"},
		},
		{
			name:     "block tool2",
			filter:   &nvgo.MCPToolFilterStatic{BlockedToolNames: []string{"tool2"}},
			expected: []string{"tool1", "tool3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := nvgo.ApplyMCPToolFilter(context.Background(), nvgo.MCPToolFilterContext{}, tt.filter, tools)
			assert.Len(t, filtered, len(tt.expected))
			for i, exp := range tt.expected {
				if i < len(filtered) {
					assert.Equal(t, exp, filtered[i].Name)
				}
			}
		})
	}
}
