// Package main 演示 agentgo 的 ReAct 模式
//
// ReAct (Reasoning + Acting) 是一种让 Agent 按照
// "思考 -> 行动 -> 观察" 循环来解决问题的模式。
//
// 工作流程:
//  1. Thought: 分析当前情况，思考下一步
//  2. Action: 选择并执行一个工具
//  3. Observation: 观察工具返回的结果
//  4. 重复以上步骤直到任务完成
//  5. Final Answer: 给出最终答案
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
	"encoding/json"
	"fmt"
	"log"
	"os"

	agentgo "github.com/chuanbosi666/agent_go"
	"github.com/openai/openai-go/v3"
)

// 模拟的知识库数据
var knowledgeBase = map[string]string{
	"agentgo":    "agentgo 是一个 Go 语言的 AI Agent 框架，参考 OpenAI Agents SDK 设计。支持工具调用、Guardrails、Session 管理等功能。",
	"agent":   "Agent 是 AI 系统中的智能实体，能够感知环境、做出决策并执行动作来实现目标。",
	"react":   "ReAct 是 Reasoning + Acting 的缩写，是一种让 AI 系统交替进行推理和行动的方法论。",
	"tool":    "Tool（工具）是 Agent 可以调用的外部功能，如搜索、计算、API调用等。",
	"session": "Session（会话）用于管理 Agent 的对话历史，支持多轮对话和上下文保持。",
}

// 创建搜索工具
func createSearchTool() agentgo.FunctionTool {
	return agentgo.FunctionTool{
		Name:        "search_knowledge",
		Description: "在知识库中搜索相关信息。输入关键词，返回匹配的内容。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"keyword": map[string]any{
					"type":        "string",
					"description": "要搜索的关键词",
				},
			},
			"required": []string{"keyword"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Keyword string `json:"keyword"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, err
			}

			if info, ok := knowledgeBase[params.Keyword]; ok {
				return fmt.Sprintf("找到信息: %s", info), nil
			}
			return fmt.Sprintf("未找到关于 '%s' 的信息", params.Keyword), nil
		},
	}
}

// 创建列表工具
func createListTopicsTool() agentgo.FunctionTool {
	return agentgo.FunctionTool{
		Name:        "list_topics",
		Description: "列出知识库中所有可用的主题",
		ParamsJSONSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var topics []string
			for k := range knowledgeBase {
				topics = append(topics, k)
			}
			return fmt.Sprintf("可用主题: %v", topics), nil
		},
	}
}

// 创建总结工具
func createSummarizeTool() agentgo.FunctionTool {
	return agentgo.FunctionTool{
		Name:        "summarize",
		Description: "总结收集到的信息",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"points": map[string]any{
					"type":        "array",
					"items":       map[string]any{"type": "string"},
					"description": "要总结的要点列表",
				},
			},
			"required": []string{"points"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Points []string `json:"points"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, err
			}

			summary := "总结完成:\n"
			for i, point := range params.Points {
				summary += fmt.Sprintf("%d. %s\n", i+1, point)
			}
			return summary, nil
		},
	}
}

func main() {
	// 检查 API Key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置 OPENAI_API_KEY 环境变量")
	}

	// 创建客户端
	client := openai.NewClient()

	// 使用 agentgo 内置的 ReAct 指令
	// DefaultReActInstruction 提供了标准的 ReAct 格式提示词
	reactInstructions := agentgo.DefaultReActInstruction

	// 创建 ReAct Agent
	agent := agentgo.New("ReAct研究员").
		WithInstructionsGetter(reactInstructions).
		WithModel("gpt-4o-mini").
		WithClient(client).
		WithTools([]agentgo.FunctionTool{
			createSearchTool(),
			createListTopicsTool(),
			createSummarizeTool(),
		})

	// 配置 Runner
	runner := agentgo.Runner{
		Config: agentgo.RunConfig{
			MaxTurns: 10, // ReAct 可能需要多轮迭代
		},
	}

	ctx := context.Background()

	// 测试问题 - 需要多步推理的任务
	question := "请帮我了解 agentgo 框架，包括它是什么、有什么功能。先列出可用主题，然后搜索相关信息，最后总结。"

	fmt.Println("=== ReAct 模式演示 ===")
	fmt.Printf("问题: %s\n", question)
	fmt.Println("---")

	result, err := runner.Run(ctx, agent, question)
	if err != nil {
		log.Fatalf("运行失败: %v", err)
	}

	fmt.Println("\n=== 最终回答 ===")
	fmt.Printf("%v\n", result.FinalOutput)

	// 显示执行统计
	fmt.Printf("\n=== 执行统计 ===\n")
	fmt.Printf("工具调用次数: %d\n", len(result.NewItems)/2)
	fmt.Printf("LLM 调用次数: %d\n", len(result.RawResponses))

	// 统计 token 使用
	var totalTokens uint64
	for _, resp := range result.RawResponses {
		if resp.Usage != nil {
			totalTokens += resp.Usage.TotalTokens
		}
	}
	fmt.Printf("总 Token 消耗: %d\n", totalTokens)
}
