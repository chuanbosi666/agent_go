// Package main 演示 nvgo 的多 Agent 协作模式（Agent-as-Tool）
//
// 本示例展示如何使用 WrapAgentAsTool 将一个 Agent 包装成工具，
// 让主 Agent 可以调用其他专业 Agent 来完成特定任务。
//
// 架构:
//
//	主 Agent（协调者）
//	  ├── 数学专家 Agent（作为工具）
//	  ├── 翻译专家 Agent（作为工具）
//	  └── 写作专家 Agent（作为工具）
//
// 运行前请确保设置环境变量:
//
//	export OPENAI_API_KEY=sk-xxx
//
// 运行:
//
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	nvgo "nvgo"
	"github.com/openai/openai-go/v3"
)

// 创建数学专家 Agent
func createMathExpert(client openai.Client) *nvgo.Agent {
	return nvgo.New("数学专家").
		WithInstructions(`你是一位数学专家，擅长:
- 解决数学问题
- 解释数学概念
- 执行复杂计算
- 证明数学定理

请用清晰、步骤分明的方式回答数学问题。`).
		WithModel("gpt-4o-mini").
		WithClient(client)
}

// 创建翻译专家 Agent
func createTranslator(client openai.Client) *nvgo.Agent {
	return nvgo.New("翻译专家").
		WithInstructions(`你是一位专业翻译，精通中英日韩多种语言。

翻译原则:
1. 准确传达原文含义
2. 保持语言自然流畅
3. 适当本地化表达
4. 保留专业术语

请提供高质量的翻译结果。`).
		WithModel("gpt-4o-mini").
		WithClient(client)
}

// 创建写作专家 Agent
func createWriter(client openai.Client) *nvgo.Agent {
	return nvgo.New("写作专家").
		WithInstructions(`你是一位专业写作专家，擅长:
- 文章撰写和润色
- 创意写作
- 商务文案
- 学术写作

请提供高质量、结构清晰的文字内容。`).
		WithModel("gpt-4o-mini").
		WithClient(client)
}

// 创建协调者 Agent（主 Agent）
func createCoordinator(client openai.Client, mathExpert, translator, writer *nvgo.Agent) *nvgo.Agent {
	// 将专家 Agent 包装成工具
	mathTool := nvgo.WrapAgentAsTool(mathExpert, 3)
	translateTool := nvgo.WrapAgentAsTool(translator, 3)
	writerTool := nvgo.WrapAgentAsTool(writer, 3)

	return nvgo.New("协调者").
		WithInstructions(`你是一个智能协调者，可以调用以下专家来帮助用户:

1. 数学专家 (Call_agent_数学专家): 处理数学问题、计算、公式等
2. 翻译专家 (Call_agent_翻译专家): 处理翻译任务
3. 写作专家 (Call_agent_写作专家): 处理写作、润色、文案等

工作流程:
1. 分析用户需求
2. 判断需要哪个专家
3. 调用相应专家获取结果
4. 整合结果返回给用户

如果任务需要多个专家配合，请按顺序调用。`).
		WithModel("gpt-4o-mini").
		WithClient(client).
		WithTools([]nvgo.FunctionTool{
			mathTool,
			translateTool,
			writerTool,
		})
}

func main() {
	// 检查 API Key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置 OPENAI_API_KEY 环境变量")
	}

	// 创建客户端
	client := openai.NewClient()

	// 创建专家 Agent
	mathExpert := createMathExpert(client)
	translator := createTranslator(client)
	writer := createWriter(client)

	// 创建协调者（集成所有专家）
	coordinator := createCoordinator(client, mathExpert, translator, writer)

	// 配置 Runner
	runner := nvgo.Runner{
		Config: nvgo.RunConfig{
			MaxTurns: 10, // 允许更多轮次以完成多步骤任务
		},
	}

	ctx := context.Background()

	// 测试场景
	testCases := []struct {
		name  string
		query string
	}{
		{
			name:  "数学问题",
			query: "请帮我计算 1+2+3+...+100 的和，并解释计算方法",
		},
		{
			name:  "翻译任务",
			query: "请把这句话翻译成英文：'人工智能正在改变我们的生活方式'",
		},
		{
			name:  "写作任务",
			query: "请帮我写一段关于人工智能的简短介绍，100字左右",
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\n{'='*50}\n")
		fmt.Printf("=== %s ===\n", tc.name)
		fmt.Printf("问题: %s\n", tc.query)
		fmt.Println("---")

		result, err := runner.Run(ctx, coordinator, tc.query)
		if err != nil {
			log.Printf("运行失败: %v", err)
			continue
		}

		fmt.Printf("回复:\n%v\n", result.FinalOutput)

		// 显示调用链
		if len(result.NewItems) > 0 {
			fmt.Printf("\n[调用了 %d 次子 Agent]\n", len(result.NewItems)/2) // 每次调用产生2个item（调用+结果）
		}
	}
}
