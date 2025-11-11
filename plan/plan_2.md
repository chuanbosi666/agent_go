# NVGo 项目开发与迁移计划 v2.0

## 📋 项目概述

**NVGo** 是一个用 Go 语言编写的**优雅的多智能体工作流框架**，灵感来自 OpenAI 的 Agents Python SDK 和 NVIDIA 的 NeMo Agent Toolkit。

### 核心定位
- 提供多智能体编排能力
- 支持 Model Context Protocol (MCP) 集成
- 与 OpenAI API 深度集成
- 提供 guardrails（护栏）机制用于输入/输出验证

### 项目信息
- **仓库**: `github.com/demo/nvgo`
- **Go 版本**: 1.25
- **许可证**: MIT License (Copyright 2025 qntx.sol)
- **状态**: 🚧 早期开发阶段，核心功能未完成

---

## 🚀 第一阶段：依赖包升级与迁移

### 1. OpenAI SDK v2 → v3 升级

#### 📊 升级概览
- **当前版本**: `github.com/openai/openai-go/v2 v2.1.1`
- **目标版本**: `github.com/openai/openai-go/v3 v3.7.0`
- **影响文件**: 10+ 个核心文件

#### 🔧 主要变更

**包路径变更**:
```go
// 旧版本 (v2)
"github.com/openai/openai-go/v2"
"github.com/openai/openai-go/v2/responses"
"github.com/openai/openai-go/v2/packages/param"
"github.com/openai/openai-go/v2/option"

// 新版本 (v3)
"github.com/openai/openai-go/v3"
"github.com/openai/openai-go/v3/responses"
"github.com/openai/openai-go/v3/packages/param"
"github.com/openai/openai-go/v3/option"
```

**需要修改的文件**:
1. `agent.go` - 客户端类型定义
2. `runner.go` - responses 包导入
3. `memory/session.go` - responses 包导入
4. `memory/sqlite.go` - responses 包导入
5. `prompt.go` - param 包导入
6. `setting.go` - 多个包导入
7. `tool.go` - param 包导入
8. `mcp.go` - param 包导入
9. `input.go` - responses 包导入
10. `go.mod` - 依赖声明

#### ✅ 迁移步骤

**步骤 1: 更新 go.mod**
```bash
# 移除旧版本
go mod edit -droprequire github.com/openai/openai-go/v2

# 添加新版本
go get github.com/openai/openai-go/v3@v3.7.0

# 整理依赖
go mod tidy
```

**步骤 2: 批量替换导入路径**
```bash
# 使用 sed 或 IDE 批量替换
find . -name "*.go" -type f -exec sed -i 's|github.com/openai/openai-go/v2|github.com/openai/openai-go/v3|g' {} \;
```

**步骤 3: API 变更适配**

主要 API 变更可能包括：
- 客户端初始化方式
- 请求参数结构
- 响应处理方式
- 错误处理机制

**步骤 4: 测试验证**
```bash
# 运行所有测试
go test ./...

# 检查编译
go build ./...
```

### 2. MCPServerSSE → MCPServerStreamableHTTP 迁移

#### 📊 迁移概览
- **废弃类型**: `MCPServerSSE` (已标记为 deprecated)
- **新类型**: `MCPServerStreamableHTTP`
- **传输协议**: SSE → Streamable HTTP
- **详细指南**: 参考 [MIGRATION_SSE_TO_STREAMABLE.md](MIGRATION_SSE_TO_STREAMABLE.md)

#### 🔧 主要变更

**类型对比**:
| 项目 | 旧 (SSE) | 新 (Streamable HTTP) |
|------|-----------|----------------------|
| 参数类型 | `MCPServerSSEParams` | `MCPServerStreamableHTTPParams` |
| 传输类型 | `*mcp.SSEClientTransport` | `*mcp.StreamableClientTransport` |
| 构造函数 | `NewMCPServerSSE()` | `NewMCPServerStreamableHTTP()` |

**迁移模板**:
```go
// 旧代码
sseServer := nvgo.NewMCPServerSSE(nvgo.MCPServerSSEParams{
    Transport: &mcp.SSEClientTransport{
        Endpoint: "https://api.example.com/mcp/sse",
    },
    CommonMCPServerParams: nvgo.CommonMCPServerParams{
        Name: "my-server",
    },
})

// 新代码
streamableServer := nvgo.NewMCPServerStreamableHTTP(nvgo.MCPServerStreamableHTTPParams{
    Transport: &mcp.StreamableClientTransport{
        Endpoint: "https://api.example.com/mcp/streamable",
    },
    CommonMCPServerParams: nvgo.CommonMCPServerParams{
        Name: "my-server",
    },
})
```

#### ✅ 迁移步骤

**步骤 1: 代码搜索**
```bash
# 查找所有 SSE 相关代码
grep -r "MCPServerSSE" . --exclude-dir=vendor
grep -r "SSEClientTransport" . --exclude-dir=vendor
```

**步骤 2: 批量替换**
```bash
# 替换构造函数
find . -name "*.go" -type f -exec sed -i 's/NewMCPServerSSE/NewMCPServerStreamableHTTP/g' {} \;

# 替换参数类型
find . -name "*.go" -type f -exec sed -i 's/MCPServerSSEParams/MCPServerStreamableHTTPParams/g' {} \;

# 替换传输类型
find . -name "*.go" -type f -exec sed -i 's/SSEClientTransport/StreamableClientTransport/g' {} \;
```

**步骤 3: URL 更新**
- 将端点从 `/sse` 更新为 `/streamable`
- 验证服务器端支持

**步骤 4: 测试验证**
- 测试服务器连接
- 测试工具列表获取
- 测试工具调用功能

---

## 🏗️ 第二阶段：核心功能开发

### 1. 完成 Runner 实现

#### 📋 当前状态
- **文件**: `runner.go`
- **完成度**: 10%
- **核心功能**: 未实现

#### 🎯 需要实现的功能

**基础运行流程**:
```go
func Run(ctx context.Context, agent *Agent, input string) (*RunResult, error) {
    // 1. 准备输入
    // 2. 应用输入护栏
    // 3. 获取指令
    // 4. 调用 LLM
    // 5. 处理工具调用
    // 6. 应用输出护栏
    // 7. 返回结果
}
```

**实现优先级**:
1. 基础运行流程 (高)
2. 工具调用处理 (高)
3. MCP 集成 (高)
4. 错误处理 (中)
5. 流式响应 (中)
6. 回调支持 (低)

### 2. MCP 集成完善

#### 📋 当前状态
- **文件**: `mcp.go`
- **完成度**: 80%
- **缺失部分**: 实际连接逻辑

#### 🎯 需要完善的功能

**连接管理**:
- 自动重连机制
- 连接池管理
- 超时处理
- 错误恢复

**工具同步**:
- 动态工具列表更新
- 工具缓存管理
- 版本控制

### 3. Guardrails 实现

#### 📋 当前状态
- **输入护栏**: 接口定义完成
- **输出护栏**: 接口定义完成
- **实现**: 缺失

#### 🎯 需要实现的功能

**内置护栏类型**:
- 内容过滤护栏
- 长度限制护栏
- 格式验证护栏
- 敏感信息检测护栏

**护栏执行流程**:
```go
// 输入处理
for _, guardrail := range agent.InputGuardrails {
    if err := guardrail.Validate(ctx, input); err != nil {
        return nil, err
    }
}

// 输出处理
for _, guardrail := range agent.OutputGuardrails {
    if output, err := guardrail.Process(ctx, output); err != nil {
        return nil, err
    }
}
```

---

## 🔧 第三阶段：高级功能开发

### 1. 会话管理增强

#### 当前实现
- 内存存储: ✅ 完成
- SQLite 存储: ✅ 完成

#### 需要增强
- 会话限制策略
- 自动清理机制
- 会话导出/导入
- 分布式存储支持

### 2. 工具系统优化

#### 当前实现
- 基础工具接口: ✅ 完成
- MCP 工具集成: 🔄 进行中

#### 需要优化
- 工具执行超时控制
- 并发工具调用
- 工具链式调用
- 自定义工具注册

### 3. 输出类型系统

#### 当前实现
- 基础接口: ✅ 完成
- JSON 支持: 🔄 部分完成

#### 需要完善
- 严格 JSON 模式
- 流式输出支持
- 多种输出格式
- 输出验证

---

## 📊 项目时间线

### 第一周：依赖升级
- [ ] Day 1-2: OpenAI SDK v3 升级
- [ ] Day 3-4: MCPServer 迁移
- [ ] Day 5: 测试与验证

### 第二周：Runner 核心实现
- [ ] Day 1-2: 基础运行流程
- [ ] Day 3-4: 工具调用集成
- [ ] Day 5: 错误处理

### 第三周：MCP 完善
- [ ] Day 1-2: 连接管理
- [ ] Day 3-4: 工具同步
- [ ] Day 5: 集成测试

### 第四周：Guardrails 实现
- [ ] Day 1-2: 输入护栏
- [ ] Day 3-4: 输出护栏
- [ ] Day 5: 集成测试

### 第五-六周：高级功能
- [ ] 会话管理增强
- [ ] 工具系统优化
- [ ] 输出类型完善

### 第七-八周：测试与文档
- [ ] 单元测试覆盖
- [ ] 集成测试
- [ ] 文档编写
- [ ] 示例项目

---

## �� 测试策略

### 单元测试
```bash
# 运行所有单元测试
go test -v ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 集成测试
- OpenAI API 集成测试
- MCP 服务器集成测试
- 端到端工作流测试

### 性能测试
- 并发请求测试
- 内存使用测试
- 响应时间基准测试

---

## 📚 文档计划

### API 文档
- 使用 godoc 生成 API 文档
- 完善代码注释

### 用户指南
- 快速开始教程
- 高级用法指南
- 最佳实践

### 开发者文档
- 架构设计文档
- 贡献指南
- 发布流程

---

## 🚀 部署计划

### 版本规划
- **v0.1.0**: 基础功能 MVP
- **v0.2.0**: 完整 Runner 实现
- **v0.3.0**: Guardrails 支持
- **v1.0.0**: 生产就绪版本

### 发布准备
- [ ] CI/CD 流水线搭建
- [ ] 自动化测试集成
- [ ] 版本标签管理
- [ ] Release Note 编写

---

## 🔍 风险评估

### 技术风险
1. **OpenAI API 变更**: SDK 升级可能带来破坏性变更
   - 缓解: 仔细阅读迁移文档，充分测试

2. **MCP 协议变更**: 传输协议升级可能影响兼容性
   - 缓解: 保留向后兼容性，渐进式迁移

3. **性能瓶颈**: 多 Agent 并发可能影响性能
   - 缓解: 早期性能测试，优化关键路径

### 项目风险
1. **开发进度**: 核心功能实现复杂度可能超出预期
   - 缓解: 分阶段交付，MVP 优先

2. **测试覆盖**: 复杂交互场景测试难度大
   - 缓解: 自动化测试，模拟环境

---

## 📈 成功指标

### 功能指标
- [ ] 所有单元测试通过 (目标: >90% 覆盖率)
- [ ] 基础 Agent 工作流正常运行
- [ ] MCP 工具调用成功
- [ ] Guardrails 正确执行

### 性能指标
- [ ] 单次请求响应时间 < 2秒
- [ ] 支持 100+ 并发 Agent
- [ ] 内存使用 < 500MB (空载)

### 质量指标
- [ ] 零已知安全漏洞
- [ ] 代码审查通过率 100%
- [ ] 文档完整性 > 80%

---

## 📝 待办事项清单

### 立即行动 (本周)
1. ✅ 分析现有代码依赖
2. 🔄 创建迁移计划
3. ⏳ 升级 OpenAI SDK 到 v3
4. ⏳ 迁移 MCPServerSSE 到 StreamableHTTP
5. ⏳ 运行测试验证

### 短期目标 (2-4周)
1. ⏳ 完成 Runner 核心实现
2. ⏳ 实现 Guardrails 机制
3. ⏳ 完善 MCP 集成
4. ⏳ 编写基础测试

### 中期目标 (1-2月)
1. ⏳ 实现所有高级功能
2. ⏳ 完成测试覆盖
3. ⏳ 编写完整文档
4. ⏳ 发布 v0.1.0

### 长期目标 (3-6月)
1. ⏳ 生产环境验证
2. ⏳ 社区反馈收集
3. ⏳ 功能迭代优化
4. ⏳ 发布 v1.0.0

---

## 🤝 贡献指南

### 开发环境设置
```bash
# 克隆仓库
git clone https://github.com/demo/nvgo.git
cd nvgo

# 安装依赖
go mod download

# 运行测试
go test ./...
```

### 代码规范
- 遵循 Go 官方代码规范
- 使用 gofmt 格式化代码
- 使用 golint 检查代码质量
- 所有公共函数必须有注释

### 提交规范
- feat: 新功能
- fix: 修复 bug
- docs: 文档更新
- test: 测试相关
- refactor: 重构代码

---

**计划版本**: v2.0
**创建日期**: 2025-11-04
**最后更新**: 2025-11-04
**适用于**: NVGo v0.x

---

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 创建 Issue: https://github.com/demo/nvgo/issues
- 发起 Discussion: https://github.com/demo/nvgo/discussions
- ���箱: [待补充]

---

**记住：这是一个敏捷开发计划，会根据实际进展和需求变化进行调整。保持灵活性，持续改进！** 🚀