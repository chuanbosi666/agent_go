// Package main 演示 nvgo 的 Session 会话管理功能
//
// Session 用于保存和管理对话历史，实现多轮对话。
// nvgo 提供了 SQLiteSession 实现，支持内存存储和文件存储。
//
// 主要功能:
//   - 保存对话历史
//   - 支持多轮对话上下文
//   - 持久化存储（可选）
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
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	nvgo "nvgo"
	"nvgo/pkg/memory"
	"github.com/openai/openai-go/v3"
)

func main() {
	// 检查 API Key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置 OPENAI_API_KEY 环境变量")
	}

	ctx := context.Background()

	// 创建 Session
	// 使用内存存储（程序结束后数据丢失）
	// 如需持久化，可以设置 DBPath 为文件路径，如 "./chat_history.db"
	session, err := memory.NewSQLiteSession(ctx, memory.SQLiteSessionConfig{
		SessionID: "demo-session-001",
		DBPath:    ":memory:", // 内存存储
		// DBPath: "./chat_history.db", // 文件存储（持久化）
	})
	if err != nil {
		log.Fatalf("创建 Session 失败: %v", err)
	}
	defer session.Close()

	// 创建客户端
	client := openai.NewClient()

	// 创建 Agent
	agent := nvgo.New("记忆助手").
		WithInstructions(`你是一个有记忆的 AI 助手。
你可以记住用户之前说过的内容，并在后续对话中引用。
请用友好、自然的方式与用户交流。`).
		WithModel("gpt-4o-mini").
		WithClient(client)

	// 创建带 Session 的 Runner
	runner := nvgo.Runner{
		Config: nvgo.RunConfig{
			MaxTurns: 5,
			Session:  session, // 关键：设置 Session
		},
	}

	fmt.Println("=== Session 会话演示 ===")
	fmt.Println("这是一个多轮对话演示，AI 会记住你说过的内容。")
	fmt.Println("输入 'quit' 退出，输入 'clear' 清除历史")
	fmt.Println("---")

	// 预设的演示对话（非交互模式）
	demoMode := len(os.Args) > 1 && os.Args[1] == "--demo"

	if demoMode {
		// 演示模式：使用预设对话
		runDemoConversation(ctx, runner, agent)
	} else {
		// 交互模式：从标准输入读取
		runInteractiveConversation(ctx, runner, agent, session)
	}
}

// runDemoConversation 运行预设的演示对话
func runDemoConversation(ctx context.Context, runner nvgo.Runner, agent *nvgo.Agent) {
	conversations := []string{
		"你好，我叫小明",
		"我最喜欢的颜色是蓝色",
		"我住在北京",
		"你还记得我叫什么名字吗？",
		"我之前说我喜欢什么颜色？",
		"总结一下你对我的了解",
	}

	for i, input := range conversations {
		fmt.Printf("\n[第 %d 轮]\n", i+1)
		fmt.Printf("用户: %s\n", input)

		result, err := runner.Run(ctx, agent, input)
		if err != nil {
			log.Printf("运行失败: %v", err)
			continue
		}

		fmt.Printf("AI: %v\n", result.FinalOutput)
	}

	fmt.Println("\n=== 演示结束 ===")
	fmt.Println("可以看到，AI 能够记住之前的对话内容。")
}

// runInteractiveConversation 运行交互式对话
func runInteractiveConversation(ctx context.Context, runner nvgo.Runner, agent *nvgo.Agent, session *memory.SQLiteSession) {
	scanner := bufio.NewScanner(os.Stdin)
	turn := 0

	for {
		fmt.Print("\n你: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// 特殊命令
		switch strings.ToLower(input) {
		case "quit", "exit", "q":
			fmt.Println("再见！")
			return
		case "clear":
			if err := session.ClearSession(ctx); err != nil {
				log.Printf("清除历史失败: %v", err)
			} else {
				fmt.Println("[历史已清除]")
				// 重新确保 session 存在
				session, _ = memory.NewSQLiteSession(ctx, memory.SQLiteSessionConfig{
					SessionID: "demo-session-001",
					DBPath:    ":memory:",
				})
				runner.Config.Session = session
			}
			continue
		case "history":
			// 显示历史记录数量
			items, _ := session.GetItems(ctx, -1)
			fmt.Printf("[当前历史记录: %d 条]\n", len(items))
			continue
		}

		turn++
		fmt.Printf("[第 %d 轮对话]\n", turn)

		result, err := runner.Run(ctx, agent, input)
		if err != nil {
			log.Printf("运行失败: %v", err)
			continue
		}

		fmt.Printf("AI: %v\n", result.FinalOutput)

		// 显示 token 使用（可选）
		if len(result.RawResponses) > 0 && result.RawResponses[0].Usage != nil {
			usage := result.RawResponses[0].Usage
			fmt.Printf("[tokens: %d]\n", usage.TotalTokens)
		}
	}
}
