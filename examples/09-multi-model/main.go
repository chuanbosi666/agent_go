// Package main 演示如何使用统一接口访问多种 AI 模型
//
// 通过 OpenRouter，你可以使用相同的代码访问：
//   - Claude (Anthropic)
//   - Gemini (Google)
//   - Grok (xAI)
//   - DeepSeek
//   - Kimi (Moonshot)
//   - GLM (智谱)
//   - GPT (OpenAI)
//   - 以及 100+ 其他模型
//
// 设置：
// 1. 注册 OpenRouter: https://openrouter.ai/
// 2. 获取 API Key: https://openrouter.ai/keys
// 3. 充值（按使用量计费）
//
// 环境变量:
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

// 支持的模型列表（只需要改名字，代码完全一样！）
var supportedModels = map[string]string{
	// Claude 系列（Anthropic）
	"claude-sonnet": "anthropic/claude-3.5-sonnet",   // 最强推理
	"claude-haiku":  "anthropic/claude-3.5-haiku",    // 快速便宜
	"claude-opus":   "anthropic/claude-3-opus",       // 旧版最强

	// Gemini 系列（Google）
	"gemini-flash": "google/gemini-flash-1.5",        // 快速
	"gemini-pro":   "google/gemini-pro-1.5",          // 平衡
	"gemini-exp":   "google/gemini-2.0-flash-exp",    // 实验版

	// Grok 系列（xAI/Twitter）
	"grok-beta": "x-ai/grok-beta",                    // Grok 最新版
	"grok-2":    "x-ai/grok-2-1212",                  // Grok 2

	// DeepSeek 系列（国产）
	"deepseek-chat": "deepseek/deepseek-chat",        // 对话模型
	"deepseek-coder": "deepseek/deepseek-coder",      // 代码模型

	// Kimi 系列（Moonshot AI，国产）
	"kimi": "moonshot/moonshot-v1-8k",                // Kimi 8K 上下文

	// GLM 系列（智谱 AI，国产）
	"glm-4": "zhipuai/glm-4",                         // GLM-4
	"glm-4-plus": "zhipuai/glm-4-plus",               // GLM-4 Plus

	// GPT 系列（OpenAI，作为对比）
	"gpt-4o":       "openai/gpt-4o",                  // GPT-4o
	"gpt-4o-mini":  "openai/gpt-4o-mini",             // GPT-4o Mini
}

func main() {
	// 读取配置
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置环境变量: export OPENAI_API_KEY=sk-or-v1-your-key")
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}

	// 创建统一客户端
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	ctx := context.Background()

	// 测试问题
	question := "用一句话介绍你自己的特点。"

	fmt.Println("=== 多模型测试：同一个问题，不同模型的回答 ===\n")

	// 要测试的模型（你可以根据需要修改这个列表）
	testModels := []string{
		"claude-sonnet",   // Claude
		"gemini-flash",    // Gemini
		"grok-2",          // Grok
		"deepseek-chat",   // DeepSeek
		"kimi",            // Kimi
		"glm-4",           // GLM
	}

	for _, modelKey := range testModels {
		modelName := supportedModels[modelKey]

		fmt.Printf("【%s】\n", modelKey)
		fmt.Printf("模型: %s\n", modelName)

		// 创建 Agent（只需要改模型名，其他都一样！）
		agent := nvgo.New(modelKey).
			WithInstructions("你是一个有特色的 AI 助手。").
			WithModel(modelName).  // 只有这一行不同
			WithClient(client)

		// 运行（完全相同的代码）
		result, err := nvgo.Run(ctx, agent, question)
		if err != nil {
			fmt.Printf("❌ 错误: %v\n\n", err)
			continue
		}

		// 显示结果
		fmt.Printf("回复: %v\n", result.FinalOutput)

		// Token 使用
		if len(result.RawResponses) > 0 && result.RawResponses[0].Usage != nil {
			usage := result.RawResponses[0].Usage
			fmt.Printf("Token: 输入=%d, 输出=%d, 总计=%d\n",
				usage.InputTokens, usage.OutputTokens, usage.TotalTokens)
		}

		fmt.Println()
	}

	fmt.Println("✅ 测试完成！")
	fmt.Println("💡 查看使用统计: https://openrouter.ai/activity")
}

// ========== 使用示例 ==========

// 示例 1：快速切换模型
func ExampleSwitchModel() {
	client := openai.NewClient(
		option.WithAPIKey("sk-or-v1-your-key"),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)

	// 今天用 Claude
	agentClaude := nvgo.New("助手").
		WithModel("anthropic/claude-3.5-sonnet").
		WithClient(client)

	// 明天换 Gemini，只改一行！
	agentGemini := nvgo.New("助手").
		WithModel("google/gemini-flash-1.5").
		WithClient(client)

	// 后天试试 DeepSeek
	agentDeepSeek := nvgo.New("助手").
		WithModel("deepseek/deepseek-chat").
		WithClient(client)

	_, _, _ = agentClaude, agentGemini, agentDeepSeek
}

// 示例 2：根据任务选择模型
func ExampleSelectByTask() {
	client := openai.NewClient(
		option.WithAPIKey("sk-or-v1-your-key"),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)

	// 复杂推理任务 → Claude Sonnet
	reasoningAgent := nvgo.New("推理专家").
		WithModel("anthropic/claude-3.5-sonnet").
		WithClient(client)

	// 快速对话 → Gemini Flash
	chatAgent := nvgo.New("聊天助手").
		WithModel("google/gemini-flash-1.5").
		WithClient(client)

	// 代码生成 → DeepSeek Coder
	codeAgent := nvgo.New("编程助手").
		WithModel("deepseek/deepseek-coder").
		WithClient(client)

	_, _, _ = reasoningAgent, chatAgent, codeAgent
}

// 示例 3：动态选择模型
func ExampleDynamicModel(modelName string) *nvgo.Agent {
	client := openai.NewClient(
		option.WithAPIKey("sk-or-v1-your-key"),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)

	// 从配置文件或用户输入读取模型名
	return nvgo.New("动态助手").
		WithModel(modelName). // 运行时决定
		WithClient(client)
}
