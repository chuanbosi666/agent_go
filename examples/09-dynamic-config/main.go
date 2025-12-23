package main

import (
	"context"
	"fmt"
	"log"

	agentgo "github.com/chuanbosi666/agent_go"

	"github.com/openai/openai-go/v3/packages/param"
)

func main() {
	fmt.Println("=== github.com/chuanbosi666/agent_go 动态模型配置示例 ===")

	// ========== 方式 1：手动注册配置 ==========
	fmt.Println("【方式 1】手动注册多个模型配置")

	// 创建注册表
	registry := agentgo.NewModelRegistry()

	// 注册配置 1：OpenAI
	registry.Registry(agentgo.ModelConfig{
			Name:    "openai",
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "your-openai-key",
			Model:   "gpt-4o-mini",
	})

	// 注册配置 2：本地 Ollama
	registry.Registry(agentgo.ModelConfig{
			Name:    "ollama",
			BaseURL: "http://localhost:11434/v1",
			APIKey:  "ollama",
			Model:   "llama3",
	})

	// 注册配置 3：自定义服务
	registry.Registry(agentgo.ModelConfig{
			Name:    "custom",
			BaseURL: "https://your-api.example.com/v1",
			APIKey:  "your-custom-key",
			Model:   "custom-model",
	})

	// 列出所有配置
	fmt.Println("\n已注册的配置:")
	for i, name := range registry.List() {
			config, _ := registry.Get(name)
			fmt.Printf("  %d. %s - %s (%s)\n", i+1, name, config.Model, config.BaseURL)
	}

	// ========== 方式 2：从配置创建 Agent ==========
	fmt.Println("\n【方式 2】使用配置创建 Agent")

	ctx := context.Background()

	// 使用 "openai" 配置创建 Agent
	agent1, err := registry.CreateAgent(
			"openai",
			"OpenAI助手",
			"你是一个友好的 AI 助手",
	)
	if err != nil {
			log.Printf("创建 Agent 失败: %v (可能 API Key 未配置)\n", err)
	} else {
			fmt.Printf("✓ 创建 Agent: %s (模型: gpt-4o-mini)\n", agent1.Name)
	}

	// 使用 "ollama" 配置创建另一个 Agent
	agent2, err := registry.CreateAgent(
			"ollama",
			"Ollama助手",
			"你是本地运行的 AI 助手",
	)
	if err != nil {
			log.Printf("创建 Agent 失败: %v (可能 Ollama 未启动)\n", err)
	} else {
			fmt.Printf("✓ 创建 Agent: %s (模型: llama3)\n", agent2.Name)
	}

	// ========== 方式 3：高级配置 ==========
	fmt.Println("\n【方式 3】使用高级配置（添加工具、护栏等）")

	agent3, err := registry.CreateAgentWithOptions("openai", func(a *agentgo.Agent) *agentgo.Agent {
			return a.
					WithInstructions("你是一个专业的编程助手").
					WithModelSettings(agentgo.ModelSettings{
							Temperature: param.NewOpt(0.7),
							MaxTokens:   param.NewOpt[int64](2000),
					})
			// 可以继续添加：
			// .WithTools(myTools)
			// .WithInputGuardrails(myGuardrails)
	})

	if err != nil {
			log.Printf("创建高级 Agent 失败: %v\n", err)
	} else {
			fmt.Printf("✓ 创建高级 Agent: %s\n", agent3.Name)
	}

	// ========== 方式 4：保存配置到文件 ==========
	fmt.Println("\n【方式 4】保存配置到 JSON 文件")

	err = agentgo.SaveToFile(registry, "models.json")
	if err != nil {
			log.Printf("保存失败: %v\n", err)
	} else {
			fmt.Println("✓ 配置已保存到 models.json")
	}

	// ========== 方式 5：从文件加载配置 ==========
	fmt.Println("\n【方式 5】从 JSON 文件加载配置")

	registry2, err := agentgo.LoadFromFile("models.json")
	if err != nil {
			log.Printf("加载失败: %v\n", err)
	} else {
			fmt.Printf("✓ 成功加载 %d 个配置\n", registry2.Count())
			for _, name := range registry2.List() {
					config, _ := registry2.Get(name)
					fmt.Printf("  - %s: %s\n", name, config.Model)
			}
	}

	// ========== 方式 6：实际使用（如果有可用的 API Key） ==========
	fmt.Println("\n【方式 6】实际测试（需要有效的 API Key）")

	// 如果你有 OpenAI API Key，可以测试
	if agent1 != nil {
			result, err := agentgo.Run(ctx, agent1, "用一句话介绍你自己")
			if err != nil {
					log.Printf("运行失败: %v\n", err)
			} else {
					fmt.Printf("✓ 回复: %v\n", result.FinalOutput)
			}
	}

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 提示:")
	fmt.Println("  - 修改 API Key 后可以实际运行")
	fmt.Println("  - 查看生成的 models.json 了解配置格式")
	fmt.Println("  - 可以创建自己的配置文件复用")
}