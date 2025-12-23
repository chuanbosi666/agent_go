# 全自动软件开发多智能体系统

本示例展示如何使用 agentgo 框架构建一个完整的多智能体开发团队。

## 系统架构

```
                     ┌─────────────────────────────────────┐
                     │           Manager Agent             │
                     │      （总管理者，调度协调）          │
                     └──────────────────┬──────────────────┘
                                        │
         ┌──────────────────────────────┼──────────────────────────────┐
         │                              │                              │
         ▼                              ▼                              ▼
┌─────────────────┐            ┌─────────────────┐            ┌─────────────────┐
│   REQ-Agent     │            │   ARCH-Agent    │            │   CODE-Agent    │
│   产品经理       │ ────────> │   架构师         │ ────────> │   程序员         │
│   需求分析       │            │   系统设计       │            │   代码实现       │
└─────────────────┘            └─────────────────┘            └────────┬────────┘
                                                                       │
                                                                       ▼
                                                              ┌─────────────────┐
                                                              │   TEST-Agent    │
                                                              │   测试员         │
                                                              │   测试验证       │
                                                              └─────────────────┘
```

## 团队成员

| Agent | 角色 | 职责 | 工具 |
|-------|------|------|------|
| Manager | 总管理者 | 调度协调，流程控制 | 调用其他 Agent |
| REQ-Agent | 产品经理 | 需求分析，文档生成 | 无 |
| ARCH-Agent | 架构师 | 系统设计，技术选型 | 读取文件、搜索 |
| CODE-Agent | 程序员 | 代码实现 | 读写文件、搜索 |
| TEST-Agent | 测试员 | 测试编写、执行 | 读写文件、执行命令 |

## 工具集

### 文件操作 (tools/file.go)
- `read_file` - 读取文件内容
- `write_file` - 写入文件
- `list_dir` - 列出目录内容
- `append_file` - 追加文件内容

### 命令执行 (tools/exec.go)
- `exec_command` - 执行系统命令
- `go_test` - 运行 Go 测试
- `go_build` - 构建 Go 项目

### 搜索工具 (tools/search.go)
- `search_files` - 按 glob 模式搜索文件
- `search_content` - 搜索文件内容
- `find_symbol` - 查找符号定义

## 安全特性

- **目录沙箱**：所有文件操作限制在指定的项目目录内
- **路径校验**：防止 `../` 等路径穿越攻击
- **命令白名单**：只允许执行预定义的安全命令

## 使用方法

### 1. 设置环境变量

```bash
export OPENAI_API_KEY=your-api-key
export OPENAI_BASE_URL=your-api-base-url  # 可选
export OPENAI_MODEL=gpt-4o-mini           # 可选
```

### 2. 运行示例

```bash
cd examples/10-dev-team
go run main.go
```

### 3. 交互流程

```
请输入项目路径（留空使用当前目录）: /path/to/my-project
项目目录: /path/to/my-project
---
请输入您的需求: 创建一个简单的 HTTP 服务，包含健康检查接口

开始执行开发流程...
==================================================
```

## 示例输入

### 示例 1：创建新项目

```
请输入您的需求: 创建一个 TODO 应用的后端 API，支持增删改查操作
```

### 示例 2：修改现有项目

```
请输入您的需求: 为现有项目添加用户认证功能，使用 JWT
```

## 工作流程

1. **需求分析阶段**
   - Manager 调用 REQ-Agent
   - 输出：结构化需求文档

2. **架构设计阶段**
   - Manager 调用 ARCH-Agent
   - 输入：需求文档
   - 输出：架构设计文档

3. **代码实现阶段**
   - Manager 调用 CODE-Agent
   - 输入：需求 + 架构
   - 输出：实现的代码文件

4. **测试验证阶段**
   - Manager 调用 TEST-Agent
   - 输入：需求 + 代码
   - 输出：测试代码和结果

## 项目结构

```
10-dev-team/
├── main.go              # 入口 + Manager Agent
├── agents/
│   ├── req.go           # REQ-Agent（产品经理）
│   ├── arch.go          # ARCH-Agent（架构师）
│   ├── code.go          # CODE-Agent（程序员）
│   └── test.go          # TEST-Agent（测试员）
├── tools/
│   ├── file.go          # 文件操作工具
│   ├── exec.go          # 命令执行工具
│   └── search.go        # 搜索工具
└── README.md            # 本文件
```

## 扩展建议

### 添加更多 Agent

```go
// 创建新的专业 Agent
reviewAgent := agents.CreateReviewAgent(client, model, tools)

// 包装成工具
reviewTool := agentgo.WrapAgentAsTool(reviewAgent, 5)

// 添加到 Manager
managerAgent.WithTools([]tool.FunctionTool{
    reqTool,
    archTool,
    codeTool,
    testTool,
    reviewTool,  // 新增
})
```

### 自定义工具

```go
// 在 tools/ 目录下创建新工具
func CreateCustomTool() tool.FunctionTool {
    return tool.FunctionTool{
        Name:        "custom_tool",
        Description: "自定义工具描述",
        ParamsJSONSchema: map[string]any{...},
        OnInvokeTool: func(ctx context.Context, args string) (any, error) {
            // 实现逻辑
        },
    }
}
```

## 注意事项

1. **API 成本**：多 Agent 系统会产生较多 API 调用，注意控制成本
2. **执行时间**：完整流程可能需要几分钟，请耐心等待
3. **项目目录**：确保对项目目录有读写权限
4. **网络连接**：需要稳定的网络连接到 API 服务器

## 许可证

MIT License
