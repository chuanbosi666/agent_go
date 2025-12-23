// Package main 演示全自动软件开发多智能体系统
//
// 本示例展示如何使用 agentgo 框架构建一个完整的开发团队，包括：
// - Manager Agent（总管理者）：调度各个专业 Agent
// - REQ-Agent（产品经理）：需求分析
// - ARCH-Agent（架构师）：系统设计
// - CODE-Agent（程序员）：代码实现
// - TEST-Agent（测试员）：测试验证
//
// 运行前请设置环境变量:
//
//	export OPENAI_API_KEY=your-api-key
//	export OPENAI_BASE_URL=your-api-base-url (可选)
//	export OPENAI_MODEL=gpt-4o-mini (可选)
//
// 运行:
//
//	go run main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	agentgo "github.com/chuanbosi666/agent_go"
	"github.com/chuanbosi666/agent_go/examples/10-dev-team/agents"
	"github.com/chuanbosi666/agent_go/examples/10-dev-team/tools"
	"github.com/chuanbosi666/agent_go/pkg/tool"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// ManagerAgentInstructions Manager Agent 的系统指令
const ManagerAgentInstructions = `你是一个软件开发团队的管理者，负责协调整个开发流程。

## 团队成员

你有以下专业 Agent 可以调用：

1. **call_agent_REQ-Agent（产品经理）**
   - 负责需求分析
   - 输入：用户的原始需求描述
   - 输出：结构化的需求文档

2. **call_agent_ARCH-Agent（架构师）**
   - 负责系统架构设计
   - 输入：需求文档
   - 输出：架构设计文档

3. **call_agent_CODE-Agent（程序员）**
   - 负责代码实现
   - 输入：需求文档 + 架构设计
   - 输出：实现的代码文件

4. **call_agent_TEST-Agent（测试员）**
   - 负责编写和运行测试
   - 输入：需求文档 + 代码文件
   - 输出：测试代码和测试结果

## 工作流程

请按照以下流程执行任务：

### 第一阶段：需求分析
1. 将用户需求传递给 REQ-Agent
2. 获取结构化的需求文档
3. 向用户确认需求理解是否正确

### 第二阶段：架构设计
1. 将需求文档传递给 ARCH-Agent
2. 获取架构设计文档
3. 确保设计满足需求

### 第三阶段：代码实现
1. 将需求和架构传递给 CODE-Agent
2. CODE-Agent 会创建必要的文件
3. 确保代码实现完整

### 第四阶段：测试验证
1. 将需求和代码传递给 TEST-Agent
2. TEST-Agent 编写测试并运行
3. 报告测试结果

## 输出格式

在每个阶段完成后，请输出简要总结：

---
## 阶段 N: [阶段名称]
### 执行情况
[简要描述完成的工作]

### 输出摘要
[关键输出内容摘要]

### 下一步
[下一步计划]
---

## 错误处理

如果某个阶段失败：
1. 分析失败原因
2. 尝试修正输入后重试（最多 2 次）
3. 如果仍然失败，向用户报告问题并请求指导

## 最终报告

所有阶段完成后，输出最终报告：

---
# 开发完成报告

## 项目概述
[项目名称和描述]

## 完成情况
- [x] 需求分析：完成
- [x] 架构设计：完成
- [x] 代码实现：完成
- [x] 测试验证：完成

## 生成的文件
[列出所有创建的文件]

## 测试结果
[测试通过/失败情况]

## 后续建议
[可能的改进建议]
---

## 重要提醒

1. 每个阶段的输出要作为下个阶段的输入
2. 保持上下文的连贯性
3. 如遇问题，及时向用户报告
4. 确保每个步骤都有明确的产出`

func main() {
	// 读取配置
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置 OPENAI_API_KEY 环境变量")
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	// 创建 OpenAI 客户端
	var client openai.Client
	if baseURL != "" {
		client = openai.NewClient(
			option.WithAPIKey(apiKey),
			option.WithBaseURL(baseURL),
		)
	} else {
		client = openai.NewClient(
			option.WithAPIKey(apiKey),
		)
	}

	// 获取项目路径
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入项目路径（留空使用当前目录）: ")
	projectPath, _ := reader.ReadString('\n')
	projectPath = strings.TrimSpace(projectPath)

	if projectPath == "" {
		var err error
		projectPath, err = os.Getwd()
		if err != nil {
			log.Fatalf("获取当前目录失败: %v", err)
		}
	}

	// 转换为绝对路径
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		log.Fatalf("获取绝对路径失败: %v", err)
	}
	projectPath = absPath

	// 确保项目目录存在
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		log.Fatalf("创建项目目录失败: %v", err)
	}

	fmt.Printf("项目目录: %s\n", projectPath)
	fmt.Println("---")

	// 初始化工具集
	fileTools := tools.NewFileTools(projectPath)
	execTools := tools.NewExecTools(projectPath)
	searchTools := tools.NewSearchTools(projectPath)

	// 收集所有工具
	allFileTools := fileTools.GetAllTools()
	allExecTools := execTools.GetAllTools()
	allSearchTools := searchTools.GetAllTools()

	// 架构师工具：文件读取 + 搜索
	archTools := []tool.FunctionTool{
		fileTools.CreateReadFileTool(),
		fileTools.CreateListDirTool(),
	}
	archTools = append(archTools, allSearchTools...)

	// 程序员工具：完整文件操作 + 搜索
	codeTools := append(allFileTools, allSearchTools...)

	// 测试员工具：文件操作 + 执行 + 搜索
	testTools := append(allFileTools, allExecTools...)
	testTools = append(testTools, allSearchTools...)

	// 创建各个专业 Agent
	reqAgent := agents.CreateREQAgent(client, model)
	archAgent := agents.CreateARCHAgent(client, model, archTools)
	codeAgent := agents.CreateCODEAgent(client, model, codeTools)
	testAgent := agents.CreateTESTAgent(client, model, testTools)

	// 将 Agent 包装成工具
	reqTool := agentgo.WrapAgentAsTool(reqAgent, 5)
	archTool := agentgo.WrapAgentAsTool(archAgent, 10)
	codeTool := agentgo.WrapAgentAsTool(codeAgent, 15)
	testTool := agentgo.WrapAgentAsTool(testAgent, 10)

	// 创建 Manager Agent
	managerAgent := agentgo.New("Manager").
		WithInstructions(ManagerAgentInstructions).
		WithModel(model).
		WithClient(client).
		WithTools([]tool.FunctionTool{
			reqTool,
			archTool,
			codeTool,
			testTool,
		})

	// 获取用户需求
	fmt.Print("请输入您的需求: ")
	requirement, _ := reader.ReadString('\n')
	requirement = strings.TrimSpace(requirement)

	if requirement == "" {
		log.Fatal("需求不能为空")
	}

	fmt.Println("\n开始执行开发流程...")
	fmt.Println("=" + strings.Repeat("=", 50))

	// 配置 Runner
	runner := agentgo.Runner{
		Config: agentgo.RunConfig{
			MaxTurns: 50, // 允许足够多的轮次完成整个流程
		},
	}

	// 执行
	ctx := context.Background()
	result, err := runner.Run(ctx, managerAgent, requirement)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	// 输出结果
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("最终输出:")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println(result.FinalOutput)

	// 统计信息
	fmt.Println("\n" + strings.Repeat("-", 50))
	fmt.Printf("执行统计:\n")
	fmt.Printf("  - Agent 调用次数: %d\n", countToolCalls(result))
	fmt.Printf("  - 总响应数: %d\n", len(result.RawResponses))
}

// countToolCalls 统计工具调用次数
// 每个工具调用通常产生2个 item（调用+结果），所以除以2
func countToolCalls(result *agentgo.RunResult) int {
	return len(result.NewItems) / 2
}
