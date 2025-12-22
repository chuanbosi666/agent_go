// Package main 演示 nvgo 的 Guardrails（护栏）功能
//
// Guardrails 用于在 Agent 执行前后进行安全检查:
//   - InputGuardrail: 检查用户输入，防止恶意或不当内容
//   - OutputGuardrail: 检查 Agent 输出，确保安全合规
//
// 当护栏触发时（TripwireTriggered = true），Agent 执行会立即停止。
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
	"strings"

	nvgo "nvgo"
	"github.com/openai/openai-go/v3"
)

// 敏感词列表
var sensitiveWords = []string{"密码", "信用卡", "银行卡号", "身份证"}

// 禁止话题
var forbiddenTopics = []string{"非法", "暴力", "犯罪"}

// 创建输入护栏：检查敏感信息
func createSensitiveInfoGuardrail() nvgo.InputGuardrail {
	return nvgo.NewInputGuardrail("sensitive_info_check",
		func(ctx context.Context, agent nvgo.AgentLike, input nvgo.Input) (nvgo.GuardrailFunctionOutput, error) {
			// 获取输入文本
			var inputText string
			switch v := input.(type) {
			case nvgo.InputString:
				inputText = string(v)
			default:
				// 其他类型暂不检查
				return nvgo.GuardrailFunctionOutput{
					TripwireTriggered: false,
				}, nil
			}

			// 检查是否包含敏感词
			for _, word := range sensitiveWords {
				if strings.Contains(inputText, word) {
					return nvgo.GuardrailFunctionOutput{
						TripwireTriggered: true,
						OutputInfo: map[string]any{
							"reason":         "输入包含敏感信息",
							"sensitive_word": word,
						},
					}, nil
				}
			}

			return nvgo.GuardrailFunctionOutput{
				TripwireTriggered: false,
				OutputInfo:        "输入检查通过",
			}, nil
		})
}

// 创建输入护栏：检查禁止话题
func createForbiddenTopicGuardrail() nvgo.InputGuardrail {
	return nvgo.NewInputGuardrail("forbidden_topic_check",
		func(ctx context.Context, agent nvgo.AgentLike, input nvgo.Input) (nvgo.GuardrailFunctionOutput, error) {
			var inputText string
			switch v := input.(type) {
			case nvgo.InputString:
				inputText = string(v)
			default:
				return nvgo.GuardrailFunctionOutput{TripwireTriggered: false}, nil
			}

			for _, topic := range forbiddenTopics {
				if strings.Contains(inputText, topic) {
					return nvgo.GuardrailFunctionOutput{
						TripwireTriggered: true,
						OutputInfo: map[string]any{
							"reason": "涉及禁止话题",
							"topic":  topic,
						},
					}, nil
				}
			}

			return nvgo.GuardrailFunctionOutput{TripwireTriggered: false}, nil
		})
}

// 创建输出护栏：检查输出长度
func createOutputLengthGuardrail(maxLength int) nvgo.OutputGuardrail {
	return nvgo.NewOutputGuardrail("output_length_check",
		func(ctx context.Context, agent nvgo.AgentLike, output any) (nvgo.GuardrailFunctionOutput, error) {
			outputStr := fmt.Sprintf("%v", output)

			if len(outputStr) > maxLength {
				return nvgo.GuardrailFunctionOutput{
					TripwireTriggered: true,
					OutputInfo: map[string]any{
						"reason":     "输出超过长度限制",
						"max_length": maxLength,
						"actual":     len(outputStr),
					},
				}, nil
			}

			return nvgo.GuardrailFunctionOutput{
				TripwireTriggered: false,
				OutputInfo: map[string]any{
					"length": len(outputStr),
				},
			}, nil
		})
}

// 创建输出护栏：检查输出中是否包含不当内容
func createContentSafetyGuardrail() nvgo.OutputGuardrail {
	blockedPhrases := []string{"我无法", "我不能", "作为AI"}

	return nvgo.NewOutputGuardrail("content_safety_check",
		func(ctx context.Context, agent nvgo.AgentLike, output any) (nvgo.GuardrailFunctionOutput, error) {
			outputStr := fmt.Sprintf("%v", output)

			// 这里只是演示，实际可以调用内容安全API
			for _, phrase := range blockedPhrases {
				if strings.Contains(outputStr, phrase) {
					// 注意：这里我们不触发tripwire，只是记录
					// 实际使用中可以根据需要决定是否阻止
					return nvgo.GuardrailFunctionOutput{
						TripwireTriggered: false, // 不阻止，只记录
						OutputInfo: map[string]any{
							"warning": "检测到可能的模式化回复",
							"phrase":  phrase,
						},
					}, nil
				}
			}

			return nvgo.GuardrailFunctionOutput{
				TripwireTriggered: false,
				OutputInfo:        "内容安全检查通过",
			}, nil
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

	// 创建带护栏的 Agent
	agent := nvgo.New("安全助手").
		WithInstructions("你是一个安全的 AI 助手，请礼貌地回答用户问题。").
		WithModel("gpt-4o-mini").
		WithClient(client).
		WithInputGuardrails([]nvgo.InputGuardrail{
			createSensitiveInfoGuardrail(),
			createForbiddenTopicGuardrail(),
		}).
		WithOutputGuardrails([]nvgo.OutputGuardrail{
			createOutputLengthGuardrail(1000),
			createContentSafetyGuardrail(),
		})

	ctx := context.Background()

	// 测试用例
	testCases := []struct {
		name   string
		input  string
		expect string // "pass" 或 "block"
	}{
		{
			name:   "正常问题",
			input:  "你好，今天天气怎么样？",
			expect: "pass",
		},
		{
			name:   "包含敏感信息",
			input:  "帮我查一下密码是多少",
			expect: "block",
		},
		{
			name:   "禁止话题",
			input:  "教我一些非法的事情",
			expect: "block",
		},
		{
			name:   "正常技术问题",
			input:  "什么是人工智能？",
			expect: "pass",
		},
	}

	fmt.Println("=== Guardrails 护栏演示 ===\n")

	for _, tc := range testCases {
		fmt.Printf("测试: %s\n", tc.name)
		fmt.Printf("输入: %s\n", tc.input)
		fmt.Printf("预期: %s\n", tc.expect)

		result, err := nvgo.Run(ctx, agent, tc.input)

		if err != nil {
			// 检查是否是护栏触发的错误
			if guardrailErr, ok := err.(*nvgo.GuardrailTripwireTriggeredError); ok {
				fmt.Printf("结果: 被护栏拦截 ✓\n")
				fmt.Printf("护栏: %s\n", guardrailErr.GuardrailName)
				fmt.Printf("详情: %v\n", guardrailErr.OutputInfo)
			} else {
				fmt.Printf("结果: 其他错误 - %v\n", err)
			}
		} else {
			fmt.Printf("结果: 正常通过 ✓\n")
			fmt.Printf("回复: %v\n", result.FinalOutput)

			// 显示护栏检查结果
			if len(result.InputGuardrailResults) > 0 {
				fmt.Println("输入护栏检查:")
				for _, gr := range result.InputGuardrailResults {
					fmt.Printf("  - %s: %v\n", gr.Guardrail.Name, gr.Output.OutputInfo)
				}
			}
			if len(result.OutputGuardrailResults) > 0 {
				fmt.Println("输出护栏检查:")
				for _, gr := range result.OutputGuardrailResults {
					fmt.Printf("  - %s: %v\n", gr.Guardrail.Name, gr.Output.OutputInfo)
				}
			}
		}
		fmt.Println("---")
	}
}
