// Package main 演示 agentgo 框架的基础用法
//
// 本示例展示如何创建一个简单的 AI Agent 并与其对话。
// 这是使用 agentgo 的最简单方式，适合快速入门。
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

	agentgo "github.com/chuanbosi666/agent_go"
	"github.com/openai/openai-go/v3"
)

func main() {
	// 1. 检查 API Key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置 OPENAI_API_KEY 环境变量")
	}

	// 2. 创建 OpenAI 客户端
	client := openai.NewClient()

	// 3. 创建 Agent
	// 使用 Builder 模式配置 Agent 的各项属性
	agent := agentgo.New("助手"). // Agent 名称
					WithInstructions("你是一个友好的 AI 助手，请用简洁的中文回答问题。"). // 系统提示词
					WithModel("gpt-4o-mini").                             // 使用的模型
					WithClient(client)                                    // OpenAI 客户端

	// 4. 创建 Runner 并执行
	ctx := context.Background()
	result, err := agentgo.Run(ctx, agent, "你好！请介绍一下你自己。")
	if err != nil {
		log.Fatalf("运行失败: %v", err)
	}

	// 5. 输出结果
	fmt.Println("=== Agent 回复 ===")
	fmt.Println(result.FinalOutput)

	// 6. 可选：查看 token 使用情况
	if len(result.RawResponses) > 0 && result.RawResponses[0].Usage != nil {
		usage := result.RawResponses[0].Usage
		fmt.Printf("\n=== Token 使用 ===\n")
		fmt.Printf("输入 tokens: %d\n", usage.InputTokens)
		fmt.Printf("输出 tokens: %d\n", usage.OutputTokens)
		fmt.Printf("总计 tokens: %d\n", usage.TotalTokens)
	}
}
