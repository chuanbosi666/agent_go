package nvgo_test

import (
	"context"
	"testing"
	"errors"

	"nvgo/pkg/types"
	"nvgo/pkg/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/google/jsonschema-go/jsonschema"
)

type MockMCPServer struct{
	name					string
	useStructuredContent 	bool
	tools					[]*mcp.Tool
	callToolFunc			func(ctx context.Context, name string, args map[string]any)(*mcp.CallToolResult,error)
	prompts					[]*mcp.Prompt
}

func (m *MockMCPServer) Connect(ctx context.Context) error{
	return nil
}
func (m *MockMCPServer) Cleanup(ctx context.Context) error{
	return nil
}

func(m *MockMCPServer) Name() string{
	if m.name == ""{
		return "mock-server"
	}
	return m.name
}

func (m *MockMCPServer) UseStructuredContent() bool {
        return m.useStructuredContent
}

func (m *MockMCPServer) ListTools(ctx context.Context, a types.AgentLike) ([]*mcp.Tool, error) {
        return m.tools, nil
}

func (m *MockMCPServer) CallTool(ctx context.Context, name string, args map[string]any) (*mcp.CallToolResult, error){
	if m.callToolFunc != nil{
		return m.callToolFunc(ctx, name, args)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "mock result"},
		},
	}, nil
}

func (m *MockMCPServer) ListPrompts(ctx context.Context) (*mcp.ListPromptsResult, error) {
        return &mcp.ListPromptsResult{Prompts: m.prompts}, nil
}

func (m *MockMCPServer) GetPrompt(ctx context.Context, name string, args map[string]string)(*mcp.GetPromptResult, error) {
        return &mcp.GetPromptResult{}, nil
}
func TestNewMCPToolFilterStatic(t *testing.T) {
	tests := []struct {
		name     string
		allowed  []string
		blocked  []string
		expected tool.MCPToolFilter
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
			filter, ok := tool.NewMCPToolFilterStatic(tt.allowed, tt.blocked)
			assert.Equal(t, tt.ok, ok)
			if tt.ok {
				assert.IsType(t, &tool.MCPToolFilterStatic{}, filter)
				actual := filter.(*tool.MCPToolFilterStatic)
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
			filter := &tool.MCPToolFilterStatic{
				AllowedToolNames: tt.allowed,
				BlockedToolNames: tt.blocked,
			}
			mcpTool := &mcp.Tool{Name: tt.toolName}
			include, err := filter.FilterMCPTool(context.Background(), tool.MCPToolFilterContext{}, mcpTool)
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
		filter   tool.MCPToolFilter
		expected []string
	}{
		{
			name:     "no filter",
			filter:   nil,
			expected: []string{"tool1", "tool2", "tool3"},
		},
		{
			name:     "allow tool1",
			filter:   &tool.MCPToolFilterStatic{AllowedToolNames: []string{"tool1"}},
			expected: []string{"tool1"},
		},
		{
			name:     "block tool2",
			filter:   &tool.MCPToolFilterStatic{BlockedToolNames: []string{"tool2"}},
			expected: []string{"tool1", "tool3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := tool.ApplyMCPToolFilter(context.Background(), tool.MCPToolFilterContext{}, tt.filter, tools)
			assert.Len(t, filtered, len(tt.expected))
			for i, exp := range tt.expected {
				if i < len(filtered) {
					assert.Equal(t, exp, filtered[i].Name)
				}
			}
		})
	}
}

func TestToFunctionTool(t *testing.T) {
      tests := []struct {
          name     string
          mcpTool  *mcp.Tool
          strict   bool
          wantName string
          wantDesc string
      }{
          {
              name: "basic tool with schema",
              mcpTool: &mcp.Tool{
                  Name:        "get_weather",
                  Description: "Get weather information",
                  InputSchema: &jsonschema.Schema{
                      Type: "object",
                      Properties: map[string]*jsonschema.Schema{
                          "city": {
                              Type:        "string",
                              Description: "City name",
                          },
                      },
                      Required: []string{"city"},
                  },
              },
              strict:   false,
              wantName: "get_weather",
              wantDesc: "Get weather information",
          },
          {
              name: "tool without schema",
              mcpTool: &mcp.Tool{
                  Name:        "no_params",
                  Description: "Tool with no parameters",
              },
              strict:   false,
              wantName: "no_params",
              wantDesc: "Tool with no parameters",
          },
          {
              name: "strict mode",
              mcpTool: &mcp.Tool{
                  Name: "strict_tool",
                  InputSchema: &jsonschema.Schema{
                      Type: "object",
                      Properties: map[string]*jsonschema.Schema{
                          "value": {Type: "string"},
                      },
                  },
              },
              strict:   true,
              wantName: "strict_tool",
          },
      }

      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              mockServer := &MockMCPServer{}

              ft, err := tool.ToFunctionTool(tt.mcpTool, mockServer, tt.strict)

              require.NoError(t, err)
              assert.Equal(t, tt.wantName, ft.Name)
              assert.Equal(t, tt.wantDesc, ft.Description)
              assert.NotNil(t, ft.OnInvokeTool)
              assert.NotNil(t, ft.ParamsJSONSchema)

              // 验证 schema 包含 properties
              props, ok := ft.ParamsJSONSchema["properties"]
              assert.True(t, ok, "schema should have properties")
              assert.NotNil(t, props)
          })
      }
  }

  func TestInvokeMCPTool(t *testing.T) {
      tests := []struct {
          name        string
          tool        *mcp.Tool
          input       string
          mockResult  *mcp.CallToolResult
          mockErr     error
          wantContain string
          wantErr     bool
          errContains string
      }{
          {
              name:  "successful call",
              tool:  &mcp.Tool{Name: "test_tool"},
              input: `{"key": "value"}`,
              mockResult: &mcp.CallToolResult{
                  Content: []mcp.Content{
                      &mcp.TextContent{Text: "success"},
                  },
              },
              wantContain: "success",
          },
          {
              name:  "empty input",
              tool:  &mcp.Tool{Name: "no_params"},
              input: "",
              mockResult: &mcp.CallToolResult{
                  Content: []mcp.Content{},
              },
              wantContain: "[]",
          },
          {
              name:        "invalid JSON input",
              tool:        &mcp.Tool{Name: "test_tool"},
              input:       "invalid json",
              wantErr:     true,
              errContains: "invalid input",
          },
          {
              name:        "server error",
              tool:        &mcp.Tool{Name: "error_tool"},
              input:       `{}`,
              mockErr:     errors.New("server error"),
              wantErr:     true,
              errContains: "invoke error_tool",
          },
          {
              name:  "multiple content items",
              tool:  &mcp.Tool{Name: "multi_tool"},
              input: `{}`,
              mockResult: &mcp.CallToolResult{
                  Content: []mcp.Content{
                      &mcp.TextContent{Text: "first"},
                      &mcp.TextContent{Text: "second"},
                  },
              },
              wantContain: "first",
          },
      }

      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              mockServer := &MockMCPServer{
                  callToolFunc: func(ctx context.Context, name string, args map[string]any) (*mcp.CallToolResult, error) {
                      if tt.mockErr != nil {
                          return nil, tt.mockErr
                      }
                      return tt.mockResult, nil
                  },
              }

              result, err := tool.InvokeMCPTool(context.Background(), mockServer, tt.tool, tt.input)

              if tt.wantErr {
                  require.Error(t, err)
                  if tt.errContains != "" {
                      assert.Contains(t, err.Error(), tt.errContains)
                  }
                  return
              }
              require.NoError(t, err)
              assert.Contains(t, result, tt.wantContain)
          })
      }
  }

func TestGetFunctionTools(t *testing.T) {
	t.Run("get tools from server", func(t *testing.T) {
		mockServer := &MockMCPServer{
			tools: []*mcp.Tool{
				{Name: "tool1", Description: "First tool"},
				{Name: "tool2", Description: "Second tool"},
			},
		}

		tools, err := tool.GetFunctionTools(context.Background(), mockServer, false, nil)

		require.NoError(t, err)
		assert.Len(t, tools, 2)
		assert.Equal(t, "tool1", tools[0].ToolName())
		assert.Equal(t, "tool2", tools[1].ToolName())
	})

	t.Run("empty tools list", func(t *testing.T) {
		mockServer := &MockMCPServer{
			tools: []*mcp.Tool{},
		}

		tools, err := tool.GetFunctionTools(context.Background(), mockServer, false, nil)

		require.NoError(t, err)
		assert.Len(t, tools, 0)
	})
}

func TestGetAllFunctionTools(t *testing.T) {
	t.Run("multiple servers", func(t *testing.T) {
		server1 := &MockMCPServer{
			name:  "server1",
			tools: []*mcp.Tool{{Name: "tool_a", Description: "Tool A"}},
		}
		server2 := &MockMCPServer{
			name:  "server2",
			tools: []*mcp.Tool{{Name: "tool_b", Description: "Tool B"}},
		}

		tools, err := tool.GetAllFunctionTools(
			context.Background(),
			[]tool.MCPServer{server1, server2},
			false,
			nil,
		)

		require.NoError(t, err)
		assert.Len(t, tools, 2)
	})

	t.Run("duplicate tool names error", func(t *testing.T) {
		server1 := &MockMCPServer{
			name:  "server1",
			tools: []*mcp.Tool{{Name: "duplicate"}},
		}
		server2 := &MockMCPServer{
			name:  "server2",
			tools: []*mcp.Tool{{Name: "duplicate"}},
		}

		_, err := tool.GetAllFunctionTools(
			context.Background(),
			[]tool.MCPServer{server1, server2},
			false,
			nil,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate tool name")
	})

	t.Run("empty servers list", func(t *testing.T) {
		tools, err := tool.GetAllFunctionTools(
			context.Background(),
			[]tool.MCPServer{},
			false,
			nil,
		)

		require.NoError(t, err)
		assert.Len(t, tools, 0)
	})
}