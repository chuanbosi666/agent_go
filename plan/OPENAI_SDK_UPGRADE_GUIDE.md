# OpenAI Go SDK 版本升级指南

## 概述

本文档提供从 `github.com/openai/openai-go/v2 v2.1.1` 到 `github.com/openai/openai-go/v3 v3.7.0` 的升级指导。升级基于仓库分析，重点关注破坏性变更和新特性。您的项目是 `github.com/agent-go`，涉及 Agent 框架、模型设置、工具调用等。

### 主要变更
- **移除**：`InputAudio` 相关功能和类型
- **新增**：更严格的类型检查、改进的错误处理
- **改进**：Structured Outputs 支持、Azure 集成增强

升级目标：
- 提升类型安全性、错误处理和 API 对齐
- 兼容 OpenAI API 的最新规范
- 获得更好的性能和稳定性

**⚠️ 重要**：升级前请完整备份代码和 `go.mod` 文件。建议在单独的分支进行升级测试。

## 前提条件

- **Go 版本**：1.25（与当前 `go.mod` 一致）
- **当前依赖**：确认使用 `github.com/openai/openai-go/v2 v2.1.1`
- **测试覆盖率**：确保单元测试覆盖主要功能模块
- **时间估算**：预计需要 2-4 小时完成迁移和测试

## 升级前检查清单

- [ ] 备份当前代码和 `go.mod`
- [ ] 创建新的分支用于升级
- [ ] 运行现有测试套件确保通过
- [ ] 记录当前使用的 OpenAI API 功能
- [ ] 检查是否有音频相关功能（将被移除）

## 总体升级步骤

### 步骤 1：更新依赖

```bash
# 1. 创建升级分支
git checkout -b upgrade/openai-sdk-v3

# 2. 更新 go.mod 中的依赖
go mod edit -require=github.com/openai/openai-go/v3@v3.7.0

# 3. 移除旧版本依赖
go mod edit -droprequire=github.com/openai/openai-go/v2

# 4. 下载新依赖并整理
go get github.com/openai/openai-go/v3@v3.7.0
go mod tidy

# 5. 验证更新
go list -m github.com/openai/openai-go/v3
```

### 步骤 2：批量更新导入路径

使用以下命令批量替换：

```bash
# 查找所有需要更新的 Go 文件
grep -r "github.com/openai/openai-go/v2" --include="*.go" .

# 使用 sed 批量替换（Linux/Mac）
find . -name "*.go" -type f -exec sed -i '' 's|github.com/openai/openai-go/v2|github.com/openai/openai-go/v3|g' {} \;

# 或使用 PowerShell（Windows）
Get-ChildItem -Recurse -Filter "*.go" | ForEach-Object {
    (Get-Content $_.FullName) -replace 'github.com/openai/openai-go/v2', 'github.com/openai/openai-go/v3' | Set-Content $_.FullName
}
```

需要更新的导入路径示例：
```go
// v2 → v3
import (
    "github.com/openai/openai-go/v3"
    "github.com/openai/openai-go/v3/option"
    "github.com/openai/openai-go/v3/responses"
    "github.com/openai/openai-go/v3/packages/param"
)
```

### 步骤 3：处理破坏性变更

#### 3.1 InputAudio 功能移除

如果代码中有音频处理，需要完全移除：

```go
// v2 代码（需要删除）
if audioContent, ok := content.(ResponseInputContentAudio); ok {
    // 处理音频逻辑 - 删除这部分
}

// v3 中不再支持音频输入，使用纯文本
content := ResponseInputContentText{
    Type: "text",
    Text: yourTextContent,
}
```

#### 3.2 类型系统更新

v3 对类型检查更严格：

```go
// v2 代码
params := responses.ResponseNewParams{
    Model: openai.F("gpt-4"),
    Messages: openai.F([]openai.ChatCompletionMessageUnionParam{
        openai.UserMessage("Hello"),
    }),
}

// v3 代码（更严格的类型）
params := responses.ResponseNewParams{
    Model: openai.F("gpt-4"),
    Messages: openai.F([]openai.ChatCompletionMessageParam{
        openai.UserMessage("Hello"),
    }),
}
```

### 步骤 4：构建和逐步修复

```bash
# 1. 尝试构建
go build ./...

# 2. 查看编译错误
# 3. 逐个文件修复（见下一节）
# 4. 重新构建验证
```

## 文件特定修改

### agent.go

**导入更新**：
```go
import (
    "github.com/openai/openai-go/v3"
    "github.com/openai/openai-go/v3/option"
    "github.com/openai/openai-go/v3/responses"
)
```

**Azure 集成更新**（如果使用）：
```go
type Agent struct {
    Client     openai.Client
    // ... 其他字段
}

// WithClient 方法保持不变，但需要注�� v3 的初始化方式
func (a *Agent) WithClient(client openai.Client) *Agent {
    a.Client = client
    return a
}

// Azure 客户端初始化示例
func NewAzureClient(endpoint, apiKey string, scopes []string) openai.Client {
    opts := []option.RequestOption{
        option.WithAPIKey(apiKey),
        option.WithBaseURL(endpoint),
    }

    // v3 中 Azure 的正确配置方式
    if len(scopes) > 0 {
        // 注意：v3 可能不直接支持 scopes，需要检查文档
        // 这里提供框架，实际实现可能需要调整
    }

    return openai.NewClient(opts...)
}
```

**验证测试**：
```go
func TestAgentClientCompatibility(t *testing.T) {
    client := openai.NewClient()
    agent := NewAgent().WithClient(client)

    // 测试基本功能
    assert.NotNil(t, agent.Client)
}
```

### setting.go

**完整更新示例**：
```go
package main

import (
    "github.com/openai/openai-go/v3"
    "github.com/openai/openai-go/v3/option"
    "github.com/openai/openai-go/v3/responses"
    "github.com/openai/openai-go/v3/packages/param"
)

type ModelSettings struct {
    Model       string   `json:"model"`
    Temperature float64  `json:"temperature,omitempty"`
    MaxTokens   int      `json:"max_tokens,omitempty"`
    // 注意：AzureScopes 需要根据 v3 实际支持情况添加
}

func (ms ModelSettings) CustomizeResponsesRequest(
    ctx context.Context,
    params *responses.ResponseNewParams,
    opts []option.RequestOption,
) (*responses.ResponseNewParams, []option.RequestOption, error) {

    // 设置基本参数
    if ms.Model != "" {
        params.Model = openai.F(ms.Model)
    }

    if ms.Temperature > 0 {
        params.Temperature = openai.F(ms.Temperature)
    }

    if ms.MaxTokens > 0 {
        params.MaxTokens = openai.F(int64(ms.MaxTokens))
    }

    // 移除任何音频相关设置（v3 不支持）

    return params, opts, nil
}
```

### tool.go

**Structured Outputs 更新**：
```go
type ToolDefinition struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description,omitempty"`
    Parameters  map[string]interface{} `json:"parameters,omitempty"`

    // v3 中的 strict mode 支持
    Strict *bool `json:"strict,omitempty"`
}

// v3 中使用 param.Opt 更严格
func (td *ToolDefinition) ToOpenAITool() openai.ChatCompletionToolParam {
    params := openai.FunctionParameters{
        Type:       openai.F("object"),
        Properties: td.Parameters,
    }

    tool := openai.ChatCompletionToolParam{
        Type: openai.F("function"),
        Function: openai.F(openai.ChatCompletionToolParamFunction{
            Name:        openai.F(td.Name),
            Description: openai.F(td.Description),
            Parameters:  openai.F(params),
        }),
    }

    // v3 strict mode
    if td.Strict != nil && *td.Strict {
        // 根据 v3 文档设置 strict 模式
        // 具体实现取决于 v3 的 API
    }

    return tool
}
```

### runner.go

**处理 Union 类型更新**：
```go
func (r *Runner) processInput(input interface{}) error {
    // v3 中移除了音频支持
    switch v := input.(type) {
    case string:
        // 处理纯文本
        return r.processText(v)
    case []string:
        // 处理文本数组
        return r.processTexts(v)
    // 移除 audio 相关的 case
    default:
        return fmt.Errorf("unsupported input type: %T", v)
    }
}
```

### mcp.go

**简化 strict schema 处理**：
```go
// v3 提供了更好的内置支持
func (m *MCPHandler) GenerateStrictSchema() (interface{}, error) {
    // 利用 v3 的内置 strict mode
    schema := map[string]interface{}{
        "type": "object",
        "strict": true,  // v3 中直接支持
        "properties": m.Properties,
        "required": m.Required,
    }

    return schema, nil
}
```

### sqlite.go 和 session.go

**JSON 处理更新**：
```go
type SessionStorage struct {
    db *sql.DB
}

func (s *SessionStorage) SaveSession(session *Session) error {
    // v3 中不再有 InputAudio 字段
    data, err := json.Marshal(session)
    if err != nil {
        return err
    }

    // 处理旧数据中的音频字段（如果存���）
    var sessionData map[string]interface{}
    if err := json.Unmarshal(data, &sessionData); err == nil {
        // 移除 audio 相关字段
        delete(sessionData, "inputAudio")
        delete(sessionData, "audio_content")

        // 重新序列化
        data, err = json.Marshal(sessionData)
        if err != nil {
            return err
        }
    }

    _, err = s.db.Exec("INSERT INTO sessions (id, data) VALUES (?, ?)",
        session.ID, data)
    return err
}
```

## 验证和测试流程

### 构建验证

```bash
# 1. 基础构建
go build ./...

# 2. 运行测试
go test ./...

# 3. 运行竞态检测
go test -race ./...

# 4. 代码质量检查
go vet ./...
golangci-lint run  # 如果使用 golangci-lint
```

### 功能验证清单

- [ ] Agent 创建和初始化正常
- [ ] 基本对话功能工作
- [ ] 工具调用正常
- [ ] 响应解析正确
- [ ] 错误处理机制有效
- [ ] 会话存储和读取正常
- [ ] MCP 集成工作
- [ ] 无音频相关错误

### 性能基准测试

```go
// 创建 benchmarks_test.go
func BenchmarkAgentRun(b *testing.B) {
    agent := NewAgent()
    ctx := context.Background()

    for i := 0; i < b.N; i++ {
        _, err := agent.Run(ctx, "Test message")
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkToolCall(b *testing.B) {
    tool := &ToolDefinition{Name: "test_tool"}

    for i := 0; i < b.N; i++ {
        _ = tool.ToOpenAITool()
    }
}
```

运行基准测试：
```bash
go test -bench=. -benchmem
```

## 潜在风险与缓解策略

| 风险类型 | 具体问题 | 缓解措施 |
|---------|---------|---------|
| **编译错误** | 类型不匹配、导入路径错误 | 使用 `go mod tidy` 清理依赖，逐个修复编译错误 |
| **运行时错误** | 音频功能调用、Union 类型错误 | 移除所有音频相关代码，添加类型断言检查 |
| **功能缺失** | v3 移除的 API | 查找替代方案或自定义实现 |
| **性能下降** | 新版本性能变化 | 运行基准测试对比，优化热点代码 |
| **依赖冲突** | 第三方包仍依赖 v2 | 更新或替换有冲突的包 |

## 回滚计划

如果升级遇到无法解决的问题：

```bash
# 1. 回滚 go.mod
git checkout go.mod go.sum

# 2. 恢复 v2 依赖
go get github.com/openai/openai-go/v2@v2.1.1
go mod tidy

# 3. 恢复代码（如果有大量修改）
git checkout -- *.go

# 4. 重新构建
go build ./...
```

## 常见问题和解决方案

### Q1: 编译错误 "undefined: openai.ChatCompletionMessageUnionParam"

**原因**：v3 中移除了 Union 类型，改用更具体的类型。

**解决方案**：
```go
// v2 代码
messages: openai.F([]openai.ChatCompletionMessageUnionParam{...})

// v3 代码
messages: openai.F([]openai.ChatCompletionMessageParam{...})
```

### Q2: Azure 认证失败

**原因**：v3 中 Azure 配置方式可能发生变化。

**解决方案**：
```go
// 检查 v3 的 Azure 配置方式
client := openai.NewClient(
    option.WithAPIKey(apiKey),
    option.WithBaseURL("https://your-resource.openai.azure.com"),
    // v3 可能需要额外的 Azure 特定选项
)
```

### Q3: Strict JSON Schema 不工作

**原因**：v3 中 strict mode 的实现方式改变。

**解决方案**：
```go
// v3 中正确设置 strict mode
schema := openai.ChatCompletionToolParam{
    Type: openai.F("function"),
    Function: openai.F(openai.ChatCompletionToolParamFunction{
        Name: openai.F("my_function"),
        Parameters: openai.F(openai.FunctionParameters{
            Type:       openai.F("object"),
            Strict:     openai.F(true),  // v3 中的位置
            Properties: properties,
            Required:   required,
        }),
    }),
}
```

## 测试计划

### 1. 单元测试
```bash
# 运行所有单元测试
go test -v -short ./...

# 测试覆盖率
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 2. 集成测试
```bash
# 运行集成测试
go test -v -tags=integration ./...
```

### 3. 端到端测试
```bash
# 运行完整 workflow 测试
go test -v -tags=e2e ./...
```

### 4. 压力测试
```bash
# 长时间运行测试
go test -v -timeout=30m -count=1 ./...
```

## 升级完成后的验证

### 功能验证脚本
```bash
#!/bin/bash
# verify_upgrade.sh

echo "=== 验证升级 ==="

# 1. 检查依赖
echo "1. 检查 OpenAI SDK 版本..."
go list -m github.com/openai/openai-go/v3 | grep v3.7.0

# 2. 构建
echo "2. 构建项目..."
go build ./... || exit 1

# 3. 运行测试
echo "3. 运行测试..."
go test ./... || exit 1

# 4. 检查导入
echo "4. 检查是否还有 v2 引用..."
if grep -r "openai-go/v2" --include="*.go" .; then
    echo "错误：仍有 v2 引用"
    exit 1
fi

echo "✅ 升级验证通过！"
```

## 后续维护

1. **监控**：密切关注 v3 版本的更新和 bug 修复
2. **文档**：更新项目文档中的 SDK 使用示例
3. **CI/CD**：更新构建脚本以使用 v3
4. **团队培训**：分享 v3 的新特性和最佳实践

---

## 总结

升级到 OpenAI Go SDK v3 将带来：
- ✅ 更好的类型安全性和错误处理
- ✅ 改进的性能和稳定性
- ✅ 更好的 Structured Outputs 支持
- ✅ 与 OpenAI API 最新规范对齐

虽然需要处理一些破坏性变更，但长远来看，升级将使您的项目更加健壮和可维护。

如果遇到本指南未覆盖的问题，请参考：
- [OpenAI Go SDK v3 文档](https://pkg.go.dev/github.com/openai/openai-go/v3)
- [OpenAI API 文档](https://platform.openai.com/docs)
- [GitHub 仓库 issues](https://github.com/openai/openai-go/issues)