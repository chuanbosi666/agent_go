package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sync"

	"nvgo/internal/strictschema"
	"nvgo/pkg/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go/v3/packages/param"
)

// MCPToolFilterContext provides context for filtering MCP tools.
type MCPToolFilterContext struct {
	Agent      types.AgentLike
	ServerName string
}

// MCPToolFilter determines which MCP tools to include.
type MCPToolFilter interface {
	FilterMCPTool(ctx context.Context, filterCtx MCPToolFilterContext, tool *mcp.Tool) (bool, error)
}

var _ MCPToolFilter = (*MCPToolFilterStatic)(nil)

// MCPToolFilterStatic filters tools by allow/block lists.
type MCPToolFilterStatic struct {
	AllowedToolNames []string
	BlockedToolNames []string
}

// NewMCPToolFilterStatic creates a static filter. Returns nil if no filter configured.
func NewMCPToolFilterStatic(allowed, blocked []string) (MCPToolFilter, bool) {
	if len(allowed) == 0 && len(blocked) == 0 {
		return nil, false
	}
	return &MCPToolFilterStatic{AllowedToolNames: allowed, BlockedToolNames: blocked}, true
}

func (f *MCPToolFilterStatic) FilterMCPTool(_ context.Context, _ MCPToolFilterContext, t *mcp.Tool) (bool, error) {
	if len(f.AllowedToolNames) > 0 && !slices.Contains(f.AllowedToolNames, t.Name) {
		return false, nil
	}
	if len(f.BlockedToolNames) > 0 && slices.Contains(f.BlockedToolNames, t.Name) {
		return false, nil
	}
	return true, nil
}

// ApplyMCPToolFilter filters tools using the provided filter.
func ApplyMCPToolFilter(ctx context.Context, filterCtx MCPToolFilterContext, filter MCPToolFilter, tools []*mcp.Tool) []*mcp.Tool {
	if filter == nil {
		return tools
	}
	var filtered []*mcp.Tool
	for _, tool := range tools {
		include, err := filter.FilterMCPTool(ctx, filterCtx, tool)
		if err != nil || !include {
			continue
		}
		filtered = append(filtered, tool)
	}
	return filtered
}

// GetAllFunctionTools retrieves tools from multiple MCP servers.
func GetAllFunctionTools(ctx context.Context, servers []MCPServer, strict bool, a types.AgentLike) ([]Tool, error) {
	var tools []Tool
	names := make(map[string]struct{})
	for _, s := range servers {
		stools, err := GetFunctionTools(ctx, s, strict, a)
		if err != nil {
			return nil, err
		}
		for _, t := range stools {
			name := t.ToolName()
			if _, ok := names[name]; ok {
				return nil, fmt.Errorf("duplicate tool name: %q", name)
			}
			names[name] = struct{}{}
		}
		tools = append(tools, stools...)
	}
	return tools, nil
}

// GetFunctionTools retrieves tools from a single MCP server.
func GetFunctionTools(ctx context.Context, server MCPServer, strict bool, a types.AgentLike) ([]Tool, error) {
	mtools, err := server.ListTools(ctx, a)
	if err != nil {
		return nil, err
	}
	ftools := make([]Tool, 0, len(mtools))
	for _, mt := range mtools {
		ft, err := ToFunctionTool(mt, server, strict)
		if err != nil {
			return nil, err
		}
		ftools = append(ftools, ft)
	}
	return ftools, nil
}

// ToFunctionTool converts MCP tool to FunctionTool.
func ToFunctionTool(tool *mcp.Tool, server MCPServer, strict bool) (FunctionTool, error) {
	invoke := func(ctx context.Context, args string) (any, error) {
		return InvokeMCPTool(ctx, server, tool, args)
	}
	schema := map[string]any{}
	if tool.InputSchema != nil {
		b, err := json.Marshal(tool.InputSchema)
		if err != nil {
			return FunctionTool{}, fmt.Errorf("marshal input schema: %w", err)
		}
		if err := json.Unmarshal(b, &schema); err != nil {
			return FunctionTool{}, fmt.Errorf("unmarshal input schema: %w", err)
		}
	}
	if _, ok := schema["properties"]; !ok {
		schema["properties"] = map[string]any{}
	}
	isStrict := false
	if strict {
		var err error
		schema, err = strictschema.EnsureStrictJSONSchema(schema)
		if err != nil {
			return FunctionTool{}, fmt.Errorf("strict schema: %w", err)
		}
		isStrict = true
	}
	return FunctionTool{
		Name:             tool.Name,
		Description:      tool.Description,
		ParamsJSONSchema: schema,
		OnInvokeTool:     invoke,
		StrictJSONSchema: param.NewOpt(isStrict),
	}, nil
}

// InvokeMCPTool invokes an MCP tool and returns JSON result.
func InvokeMCPTool(ctx context.Context, server MCPServer, tool *mcp.Tool, input string) (string, error) {
	var data map[string]any
	if input != "" {
		if err := json.Unmarshal([]byte(input), &data); err != nil {
			return "", fmt.Errorf("invalid input for %s: %w", tool.Name, err)
		}
	}
	res, err := server.CallTool(ctx, tool.Name, data)
	if err != nil {
		return "", fmt.Errorf("invoke %s: %w", tool.Name, err)
	}
	if server.UseStructuredContent() && res.StructuredContent != nil {
		b, err := json.Marshal(res.StructuredContent)
		if err != nil {
			return "", fmt.Errorf("marshal structured: %w", err)
		}
		return string(b), nil
	}
	var b []byte
	switch len(res.Content) {
	case 0:
		return "[]", nil
	case 1:
		b, err = json.Marshal(res.Content[0])
	default:
		b, err = json.Marshal(res.Content)
	}
	if err != nil {
		return "", fmt.Errorf("marshal content: %w", err)
	}
	return string(b), nil
}

// MCPServer defines MCP server operations.
// MCPServer defines the interface for Model Context Protocol servers.
type MCPServer interface {
	Connect(context.Context) error
	Cleanup(context.Context) error
	Name() string
	UseStructuredContent() bool
	ListTools(context.Context, types.AgentLike) ([]*mcp.Tool, error)
	CallTool(context.Context, string, map[string]any) (*mcp.CallToolResult, error)
	ListPrompts(context.Context) (*mcp.ListPromptsResult, error)
	GetPrompt(context.Context, string, map[string]string) (*mcp.GetPromptResult, error)
}

// MCP error definitions.
var (
	ErrMCPServerNotInitialized = errors.New("MCP server not initialized")
	ErrMCPAgentRequired        = errors.New("agent required for tool filtering")
)

var _ MCPServer = (*MCPServerWithClientSession)(nil)

// MCPServerWithClientSession implements MCPServer using MCP ClientSession.
type MCPServerWithClientSession struct {
	transport            mcp.Transport
	session              *mcp.ClientSession
	cleanupMu            sync.Mutex
	cacheToolsList       bool
	cacheDirty           bool
	toolsList            []*mcp.Tool
	toolFilter           MCPToolFilter
	name                 string
	useStructuredContent bool
}

type MCPServerWithClientSessionParams struct {
	Name                 string
	Transport            mcp.Transport
	CacheToolsList       bool
	ToolFilter           MCPToolFilter
	UseStructuredContent bool
}

// NewMCPServerWithClientSession creates a session-based MCP server.
func NewMCPServerWithClientSession(p MCPServerWithClientSessionParams) *MCPServerWithClientSession {
	return &MCPServerWithClientSession{
		transport:            p.Transport,
		cacheToolsList:       p.CacheToolsList,
		cacheDirty:           true,
		toolFilter:           p.ToolFilter,
		name:                 p.Name,
		useStructuredContent: p.UseStructuredContent,
	}
}

func (s *MCPServerWithClientSession) Connect(ctx context.Context) error {
	client := mcp.NewClient(&mcp.Implementation{Name: s.name}, nil)
	session, err := client.Connect(ctx, s.transport, nil)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	s.session = session
	return nil
}

func (s *MCPServerWithClientSession) Cleanup(ctx context.Context) error {
	s.cleanupMu.Lock()
	defer s.cleanupMu.Unlock()
	if s.session == nil {
		return nil
	}
	err := s.session.Close()
	s.session = nil
	return err
}

func (s *MCPServerWithClientSession) Name() string               { return s.name }
func (s *MCPServerWithClientSession) UseStructuredContent() bool { return s.useStructuredContent }

func (s *MCPServerWithClientSession) ListTools(ctx context.Context, a types.AgentLike) ([]*mcp.Tool, error) {
	if s.session == nil {
		return nil, ErrMCPServerNotInitialized
	}
	var tools []*mcp.Tool
	if s.cacheToolsList && !s.cacheDirty && len(s.toolsList) > 0 {
		tools = s.toolsList
	} else {
		res, err := s.session.ListTools(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("list tools: %w", err)
		}
		s.toolsList = res.Tools
		s.cacheDirty = false
		tools = s.toolsList
	}
	if s.toolFilter == nil {
		return tools, nil
	}
	if a == nil {
		return nil, ErrMCPAgentRequired
	}
	filterCtx := MCPToolFilterContext{Agent: a, ServerName: s.name}
	return ApplyMCPToolFilter(ctx, filterCtx, s.toolFilter, tools), nil
}

func (s *MCPServerWithClientSession) CallTool(ctx context.Context, name string, args map[string]any) (*mcp.CallToolResult, error) {
	if s.session == nil {
		return nil, ErrMCPServerNotInitialized
	}
	return s.session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
}

func (s *MCPServerWithClientSession) ListPrompts(ctx context.Context) (*mcp.ListPromptsResult, error) {
	if s.session == nil {
		return nil, ErrMCPServerNotInitialized
	}
	return s.session.ListPrompts(ctx, nil)
}

func (s *MCPServerWithClientSession) GetPrompt(ctx context.Context, name string, args map[string]string) (*mcp.GetPromptResult, error) {
	if s.session == nil {
		return nil, ErrMCPServerNotInitialized
	}
	return s.session.GetPrompt(ctx, &mcp.GetPromptParams{Name: name, Arguments: args})
}

// Run connects, executes fn, then cleans up.
func (s *MCPServerWithClientSession) Run(ctx context.Context, fn func(context.Context, *MCPServerWithClientSession) error) (err error) {
	if err := s.Connect(ctx); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer func() {
		if e := s.Cleanup(ctx); e != nil {
			err = errors.Join(e, fmt.Errorf("MCP server cleanup error: %w", e))
		}
	}()
	return fn(ctx, s)
}

// InvalidateToolsCache marks cache as dirty for next fetch.
func (s *MCPServerWithClientSession) InvalidateToolsCache() { s.cacheDirty = true }

// CommonMCPServerParams shared params for MCP server types.
type CommonMCPServerParams struct {
	CacheToolsList       bool
	Name                 string
	ToolFilter           MCPToolFilter
	UseStructuredContent bool
}

// MCPServerStdio is Stdio-based MCP server.
type MCPServerStdioParams struct {
	Transport *mcp.CommandTransport
	CommonMCPServerParams
}

type MCPServerStdio struct{ *MCPServerWithClientSession }

func NewMCPServerStdio(p MCPServerStdioParams) *MCPServerStdio {
	if p.Transport == nil {
		panic("transport required")
	}
	name := p.Name
	if name == "" {
		name = fmt.Sprintf("stdio: %s", p.Transport.Command.Path)
	}
	return &MCPServerStdio{
		MCPServerWithClientSession: NewMCPServerWithClientSession(MCPServerWithClientSessionParams{
			Name:                 name,
			Transport:            p.Transport,
			CacheToolsList:       p.CacheToolsList,
			ToolFilter:           p.ToolFilter,
			UseStructuredContent: p.UseStructuredContent,
		}),
	}
}

// MCPServerStreamableHTTP is HTTP-based MCP server.
type MCPServerStreamableHTTPParams struct {
	Transport *mcp.StreamableClientTransport
	CommonMCPServerParams
}

type MCPServerStreamableHTTP struct{ *MCPServerWithClientSession }

func NewMCPServerStreamableHTTP(p MCPServerStreamableHTTPParams) *MCPServerStreamableHTTP {
	if p.Transport == nil {
		panic("transport required")
	}
	name := p.Name
	if name == "" {
		name = fmt.Sprintf("streamable_http: %s", p.Transport.Endpoint)
	}
	return &MCPServerStreamableHTTP{
		MCPServerWithClientSession: NewMCPServerWithClientSession(MCPServerWithClientSessionParams{
			Name:                 name,
			Transport:            p.Transport,
			CacheToolsList:       p.CacheToolsList,
			ToolFilter:           p.ToolFilter,
			UseStructuredContent: p.UseStructuredContent,
		}),
	}
}
