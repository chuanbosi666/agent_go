// Package main 演示如何使用自定义 API endpoint
//
// 本示例展示如何连接非 OpenAI 的兼容 API（如 Azure OpenAI、本地模型等）。
//
// 运行前请设置环境变量:
//
//	export OPENAI_API_KEY=your-api-key
//	export OPENAI_BASE_URL=https://your-api-endpoint.com/v1
//
// 或者直接在代码中配置。
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
	"github.com/openai/openai-go/v3/option"
)

func main() {
	// 方式 1：通过环境变量配置（推荐）
	// export OPENAI_API_KEY=your-api-key
	// export OPENAI_BASE_URL=https://your-api-endpoint.com/v1
	// client := openai.NewClient()

	// 方式 2：直接在代码中配置
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = "your-api-key-here" // 或从配置文件读取
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://your-api-endpoint.com/v1" // 默认值
	}

	// 创建自定义客户端
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	// 创建 Agent
	agent := agentgo.New("助手").
		WithInstructions("你是一个友好的 AI 助手，请用简洁的中文回答问题。").
		WithModel("gpt-4o-mini"). // 使用你的模型名称
		WithClient(client)

	// 运行
	ctx := context.Background()
	result, err := agentgo.Run(ctx, agent, "你好！请介绍一下你自己。")
	if err != nil {
		log.Fatalf("运行失败: %v", err)
	}

	// 输出结果
	fmt.Println("=== Agent 回复 ===")
	fmt.Println(result.FinalOutput)

	// Token 使用情况
	if len(result.RawResponses) > 0 && result.RawResponses[0].Usage != nil {
		usage := result.RawResponses[0].Usage
		fmt.Printf("\n=== Token 使用 ===\n")
		fmt.Printf("输入 tokens: %d\n", usage.InputTokens)
		fmt.Printf("输出 tokens: %d\n", usage.OutputTokens)
		fmt.Printf("总计 tokens: %d\n", usage.TotalTokens)
	}
}
