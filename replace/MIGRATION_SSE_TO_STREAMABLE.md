# 从 MCPServerSSE 迁移到 MCPServerStreamableHTTP 指南

## 📋 概述

本指南说明如何将 MCP 服务器从已废弃的 **SSE (Server-Sent Events)** 传输方式迁移到推荐的 **Streamable HTTP** 传输方式。

## ⚠️ 为什么要迁移？

- **SSE 已废弃**: `MCPServerSSE` 在代码中标记为 deprecated
- **Streamable HTTP 更现代**: 提供更好的性能和可靠性
- **长期支持**: Streamable HTTP 是 MCP SDK 推荐的传输方式

## 🔄 迁移步骤

### 需要修改的文件

在 NVGo 项目中，只有以下位置使用了 MCP 服务器：

1. **[mcp.go](e:\Lab\work\develop\write_agent\nvgo-main\mcp.go)** - MCP 服务器定义（框架代码，无需修改）
2. **您的应用代码** - 任何创建 MCP 服务器的地方

### 修改您的应用代码

#### 方案 A: 如果您直接使用 SSE 传输

**修改前 (SSE)**:
```go
import (
    "github.com/demo/nvgo"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

// 创建 SSE 传输
sseTransport := &mcp.SSEClientTransport{
    Endpoint: "https://your-mcp-server.com/sse",
}

// 创建 SSE 服务器
mcpServer := nvgo.NewMCPServerSSE(nvgo.MCPServerSSEParams{
    Transport: sseTransport,
    CommonMCPServerParams: nvgo.CommonMCPServerParams{
        Name:           "my-mcp-server",
        CacheToolsList: true,
        ToolFilter:     nil,
        UseStructuredContent: false,
    },
})
```

**修改后 (Streamable HTTP)**:
```go
import (
    "github.com/demo/nvgo"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

// 创建 Streamable HTTP 传输
streamableTransport := &mcp.StreamableClientTransport{
    Endpoint: "https://your-mcp-server.com/streamable",  // ⚠️ 注意：URL 可能需要修改
}

// 创建 Streamable HTTP 服务器
mcpServer := nvgo.NewMCPServerStreamableHTTP(nvgo.MCPServerStreamableHTTPParams{
    Transport: streamableTransport,
    CommonMCPServerParams: nvgo.CommonMCPServerParams{
        Name:           "my-mcp-server",
        CacheToolsList: true,
        ToolFilter:     nil,
        UseStructuredContent: false,
    },
})
```

#### 方案 B: 如果您使用配置化的方式

**修改前 (SSE)**:
```go
func createMCPServer(config ServerConfig) nvgo.MCPServer {
    transport := &mcp.SSEClientTransport{
        Endpoint: config.Endpoint,
    }

    return nvgo.NewMCPServerSSE(nvgo.MCPServerSSEParams{
        Transport: transport,
        CommonMCPServerParams: nvgo.CommonMCPServerParams{
            Name: config.Name,
        },
    })
}
```

**修改后 (Streamable HTTP)**:
```go
func createMCPServer(config ServerConfig) nvgo.MCPServer {
    transport := &mcp.StreamableClientTransport{
        Endpoint: config.Endpoint,
    }

    return nvgo.NewMCPServerStreamableHTTP(nvgo.MCPServerStreamableHTTPParams{
        Transport: transport,
        CommonMCPServerParams: nvgo.CommonMCPServerParams{
            Name: config.Name,
        },
    })
}
```

---

## 📝 详细对比

### 类型对比

| 项目 | SSE (旧) | Streamable HTTP (新) |
|------|----------|----------------------|
| **参数类型** | `MCPServerSSEParams` | `MCPServerStreamableHTTPParams` |
| **传输类型** | `*mcp.SSEClientTransport` | `*mcp.StreamableClientTransport` |
| **服务器类型** | `*MCPServerSSE` | `*MCPServerStreamableHTTP` |
| **构造函数** | `NewMCPServerSSE()` | `NewMCPServerStreamableHTTP()` |

### 代码结构对比

#### SSE 结构
```go
// 参数结构
type MCPServerSSEParams struct {
    Transport *mcp.SSEClientTransport
    CommonMCPServerParams
}

// 服务器结构
type MCPServerSSE struct{ *MCPServerWithClientSession }

// 构造函数
func NewMCPServerSSE(p MCPServerSSEParams) *MCPServerSSE
```

#### Streamable HTTP 结构
```go
// 参数结构
type MCPServerStreamableHTTPParams struct {
    Transport *mcp.StreamableClientTransport
    CommonMCPServerParams
}

// 服务器结构
type MCPServerStreamableHTTP struct{ *MCPServerWithClientSession }

// 构造函数
func NewMCPServerStreamableHTTP(p MCPServerStreamableHTTPParams) *MCPServerStreamableHTTP
```

---

## 🔍 完整示例

### 示例 1: 基础用法

**SSE 版本**:
```go
package main

import (
    "context"
    "fmt"
    "github.com/demo/nvgo"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    "github.com/openai/openai-go/v2"
)

func main() {
    // 创建 SSE MCP 服务器
    sseServer := nvgo.NewMCPServerSSE(nvgo.MCPServerSSEParams{
        Transport: &mcp.SSEClientTransport{
            Endpoint: "https://api.example.com/mcp/sse",
        },
        CommonMCPServerParams: nvgo.CommonMCPServerParams{
            Name: "example-sse-server",
        },
    })

    // 连接服务器
    if err := sseServer.Connect(context.Background()); err != nil {
        panic(err)
    }
    defer sseServer.Cleanup(context.Background())

    // 创建 Agent
    client := openai.NewClient()
    agent := nvgo.New("assistant").
        WithModel("gpt-4").
        WithClient(client).
        AddMCPServer(sseServer)

    // 运行
    result, err := nvgo.Run(context.Background(), agent, "Hello!")
    if err != nil {
        panic(err)
    }
    fmt.Println(result.FinalOutput)
}
```

**Streamable HTTP 版本**:
```go
package main

import (
    "context"
    "fmt"
    "github.com/demo/nvgo"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    "github.com/openai/openai-go/v2"
)

func main() {
    // 创建 Streamable HTTP MCP 服务器
    streamableServer := nvgo.NewMCPServerStreamableHTTP(nvgo.MCPServerStreamableHTTPParams{
        Transport: &mcp.StreamableClientTransport{
            Endpoint: "https://api.example.com/mcp/streamable",  // ⚠️ URL 可能不同
        },
        CommonMCPServerParams: nvgo.CommonMCPServerParams{
            Name: "example-streamable-server",
        },
    })

    // 连接服务器
    if err := streamableServer.Connect(context.Background()); err != nil {
        panic(err)
    }
    defer streamableServer.Cleanup(context.Background())

    // 创建 Agent
    client := openai.NewClient()
    agent := nvgo.New("assistant").
        WithModel("gpt-4").
        WithClient(client).
        AddMCPServer(streamableServer)

    // 运行
    result, err := nvgo.Run(context.Background(), agent, "Hello!")
    if err != nil {
        panic(err)
    }
    fmt.Println(result.FinalOutput)
}
```

---

### 示例 2: 带工具过滤

**SSE 版本**:
```go
// 创建工具过滤器
toolFilter, _ := nvgo.NewMCPToolFilterStatic(
    []string{"allowed_tool_1", "allowed_tool_2"},  // 白名单
    []string{"blocked_tool"},                       // 黑名单
)

sseServer := nvgo.NewMCPServerSSE(nvgo.MCPServerSSEParams{
    Transport: &mcp.SSEClientTransport{
        Endpoint: "https://api.example.com/mcp/sse",
    },
    CommonMCPServerParams: nvgo.CommonMCPServerParams{
        Name:           "filtered-sse-server",
        CacheToolsList: true,
        ToolFilter:     toolFilter,
    },
})
```

**Streamable HTTP 版本**:
```go
// 创建工具过滤器（完全相同）
toolFilter, _ := nvgo.NewMCPToolFilterStatic(
    []string{"allowed_tool_1", "allowed_tool_2"},  // 白名单
    []string{"blocked_tool"},                       // 黑名单
)

streamableServer := nvgo.NewMCPServerStreamableHTTP(nvgo.MCPServerStreamableHTTPParams{
    Transport: &mcp.StreamableClientTransport{
        Endpoint: "https://api.example.com/mcp/streamable",
    },
    CommonMCPServerParams: nvgo.CommonMCPServerParams{
        Name:           "filtered-streamable-server",
        CacheToolsList: true,
        ToolFilter:     toolFilter,
    },
})
```

---

### 示例 3: 使用结构化内容

**SSE 版本**:
```go
sseServer := nvgo.NewMCPServerSSE(nvgo.MCPServerSSEParams{
    Transport: &mcp.SSEClientTransport{
        Endpoint: "https://api.example.com/mcp/sse",
    },
    CommonMCPServerParams: nvgo.CommonMCPServerParams{
        Name:                 "structured-sse-server",
        UseStructuredContent: true,  // 启用结构化内容
    },
})
```

**Streamable HTTP 版本**:
```go
streamableServer := nvgo.NewMCPServerStreamableHTTP(nvgo.MCPServerStreamableHTTPParams{
    Transport: &mcp.StreamableClientTransport{
        Endpoint: "https://api.example.com/mcp/streamable",
    },
    CommonMCPServerParams: nvgo.CommonMCPServerParams{
        Name:                 "structured-streamable-server",
        UseStructuredContent: true,  // 启用结构化内容
    },
})
```

---

## ⚙️ 服务器端配置

### 如果您控制 MCP 服务器

您的 MCP 服务器需要支持 Streamable HTTP 协议。请参考：

1. **MCP 服务器文档**: https://modelcontextprotocol.io
2. **Go SDK 示例**: https://github.com/modelcontextprotocol/go-sdk

### 如果使用第三方 MCP 服务器

请确认服务器是否支持 Streamable HTTP 传输方式：

```bash
# 检查服务器是否支持 Streamable HTTP
curl -X POST https://your-server.com/streamable \
  -H "Content-Type: application/json" \
  -d '{"method": "initialize"}'
```

---

## 🔧 迁移检查清单

### 代码修改
- [ ] 将 `mcp.SSEClientTransport` 替换为 `mcp.StreamableClientTransport`
- [ ] 将 `MCPServerSSEParams` 替换为 `MCPServerStreamableHTTPParams`
- [ ] 将 `NewMCPServerSSE()` 替换为 `NewMCPServerStreamableHTTP()`
- [ ] 将 `*MCPServerSSE` 类型替换为 `*MCPServerStreamableHTTP`

### 配置修改
- [ ] 更新 Endpoint URL（从 `/sse` 改为 `/streamable`）
- [ ] 验证服务器支持 Streamable HTTP 协议
- [ ] 更新环境变量或配置文件中的 URL

### 测试
- [ ] 运行单元测试
- [ ] 测试服务器连接
- [ ] 测试工具列表获取
- [ ] 测试工具调用
- [ ] 测试错误处理

### 部署
- [ ] 更新文档
- [ ] 通知团队成员
- [ ] 部署到测试环境
- [ ] 监控错误日志
- [ ] 部署到生产环境

---

## 🐛 常见问题

### Q1: 迁移后连接失败怎么办？

**可能原因**:
1. Endpoint URL 不正确
2. 服务器不支持 Streamable HTTP

**解决方案**:
```go
// 添加错误处理
if err := streamableServer.Connect(ctx); err != nil {
    log.Printf("Failed to connect to MCP server: %v", err)
    // 检查 URL 是否正确
    // 检查服务器是否支持 Streamable HTTP
}
```

### Q2: 工具列表为空？

**可能原因**:
- 工具过滤器配置错误
- 服务器返回格式不兼容

**解决方案**:
```go
// 临时禁用缓存和过滤器进行调试
streamableServer := nvgo.NewMCPServerStreamableHTTP(nvgo.MCPServerStreamableHTTPParams{
    Transport: streamableTransport,
    CommonMCPServerParams: nvgo.CommonMCPServerParams{
        Name:           "debug-server",
        CacheToolsList: false,  // 禁用缓存
        ToolFilter:     nil,    // 禁用过滤
    },
})

// 手动列出工具
tools, err := streamableServer.ListTools(ctx, agent)
if err != nil {
    log.Printf("Error listing tools: %v", err)
} else {
    log.Printf("Found %d tools", len(tools))
    for _, tool := range tools {
        log.Printf("  - %s: %s", tool.Name, tool.Description)
    }
}
```

### Q3: 性能有差异吗？

**答案**:
- Streamable HTTP 通常比 SSE 更高效
- 支持更好的连接复用
- 降低延迟和资源消耗

**基准测试建议**:
```go
// 测试连接性能
start := time.Now()
err := streamableServer.Connect(ctx)
log.Printf("Connection time: %v", time.Since(start))

// 测试工具调用性能
start = time.Now()
result, err := streamableServer.CallTool(ctx, "tool_name", args)
log.Printf("Tool call time: %v", time.Since(start))
```

### Q4: 可以同时使用 SSE 和 Streamable HTTP 吗？

**答案**: 可以！它们都实现了 `MCPServer` 接口。

```go
// 同时使用两种传输方式
sseServer := nvgo.NewMCPServerSSE(...)
streamableServer := nvgo.NewMCPServerStreamableHTTP(...)

agent := nvgo.New("multi-transport-agent").
    AddMCPServer(sseServer).
    AddMCPServer(streamableServer)
```

---

## 📚 参考资料

### 相关代码位置

- **mcp.go:346-371** - SSE 服务器定义（已废弃）
- **mcp.go:373-398** - Streamable HTTP 服务器定义
- **mcp.go:189-310** - 基础 MCPServerWithClientSession 实现

### 外部链接

- [MCP 协议文档](https://modelcontextprotocol.io)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Streamable HTTP 规范](https://modelcontextprotocol.io/docs/specification/transport)

---

## 🎯 快速查找替换

如果您使用 IDE 或文本编辑器，可以使用以下正则表达式进行批量替换：

### 查找
```regex
NewMCPServerSSE\s*\(\s*nvgo\.MCPServerSSEParams
```

### 替换为
```regex
NewMCPServerStreamableHTTP(nvgo.MCPServerStreamableHTTPParams
```

### 同时替换传输类型

**查找**:
```regex
mcp\.SSEClientTransport
```

**替换为**:
```regex
mcp.StreamableClientTransport
```

---

## ✅ 迁移完成验证

运行以下检查确保迁移成功：

```bash
# 1. 搜索是否还有 SSE 引用
grep -r "MCPServerSSE" . --exclude-dir=vendor --exclude="*.md"

# 2. 搜索 SSEClientTransport
grep -r "SSEClientTransport" . --exclude-dir=vendor --exclude="*.md"

# 3. 运行测试
make test

# 4. 运行 linter
make lint
```

如果以上命令没有输出（或只在 mcp.go 框架代码中有输出），说明迁移成功！

---

**迁移版本**: v1.0
**最后更新**: 2025-10-28
**适用于**: NVGo v0.x
