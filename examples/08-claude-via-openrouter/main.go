// Package main 演示如何通过 OpenRouter 使用 Anthropic Claude 模型
//
// OpenRouter 提供统一的 OpenAI 兼容接口，支持 Claude、GPT、Gemini 等多种模型。
//
// 准备工作：
// 1. 注册 OpenRouter 账号: https://openrouter.ai/
// 2. 获取 API Key: https://openrouter.ai/keys
// 3. 充值一定金额（按使用量计费）
//
// 设置环境变量:
//
//	export OPENAI_API_KEY=sk-or-v1-your-openrouter-key
//	export OPENAI_BASE_URL=https://openrouter.ai/api/v1
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
	"github.com/openai/openai-go/v3/option"
)

func main() {
	// 方式 1：通过环境变量配置（推荐）
	// export OPENAI_API_KEY=sk-or-v1-your-openrouter-key
	// export OPENAI_BASE_URL=https://openrouter.ai/api/v1

	// 方式 2：直接在代码中配置
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = "sk-or-v1-your-openrouter-key-here" // 从 OpenRouter 获取
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}

	// 创建客户端
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	// 支持的 Claude 模型列表：
	// - anthropic/claude-3.5-sonnet        (最强，推荐)
	// - anthropic/claude-3.5-haiku         (快速，便宜)
	// - anthropic/claude-3-opus            (旧版最强)
	// - anthropic/claude-3-sonnet          (旧版平衡)
	// - anthropic/claude-3-haiku           (旧版快速)
	//
	// 完整模型列表: https://openrouter.ai/docs#models

	// 创建使用 Claude 的 Agent
	agent := nvgo.New("Claude助手").
		WithInstructions("你是 Claude，一个由 Anthropic 创建的 AI 助手。请用简洁的中文回答问题。").
		WithModel("anthropic/claude-3.5-sonnet"). // OpenRouter 的 Claude 模型名
		WithClient(client)

	// 运行测试
	ctx := context.Background()

	fmt.Println("=== 使用 Claude 3.5 Sonnet (通过 OpenRouter) ===\n")

	// 测试 1：简单对话
	fmt.Println("【测试 1】简单对话")
	result1, err := nvgo.Run(ctx, agent, "你好！请介绍一下你自己。")
	if err != nil {
		log.Fatalf("运行失败: %v", err)
	}
	fmt.Println("回复:", result1.FinalOutput)
	fmt.Println()

	// 测试 2：推理能力
	fmt.Println("【测试 2】推理能力")
	result2, err := nvgo.Run(ctx, agent, "如果 5 个人 5 天吃 5 个苹果，那么 10 个人 10 天吃多少个苹果？请一步步推理。")
	if err != nil {
		log.Fatalf("运行失败: %v", err)
	}
	fmt.Println("回复:", result2.FinalOutput)
	fmt.Println()

	// 显示 token 使用情况
	if len(result2.RawResponses) > 0 && result2.RawResponses[0].Usage != nil {
		usage := result2.RawResponses[0].Usage
		fmt.Printf("=== Token 使用 ===\n")
		fmt.Printf("输入 tokens: %d\n", usage.InputTokens)
		fmt.Printf("输出 tokens: %d\n", usage.OutputTokens)
		fmt.Printf("总计 tokens: %d\n", usage.TotalTokens)
	}

	fmt.Println("\n✅ Claude 模型测试成功！")
	fmt.Println("💡 提示: 你可以在 https://openrouter.ai/activity 查看使用统计")
}
