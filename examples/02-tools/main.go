// Package main 演示如何在 agentgo 中使用工具（Tools）
//
// 本示例展示如何定义 FunctionTool 并让 Agent 调用它们。
// 工具让 Agent 能够执行具体操作，如计算、查询数据等。
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
	"math"
	"os"
	"time"

	agentgo "github.com/chuanbosi666/agent_go"
	"github.com/openai/openai-go/v3"
)

// 定义工具 1: 获取当前时间
func createGetTimeTool() agentgo.FunctionTool {
	return agentgo.FunctionTool{
		Name:        "get_current_time",
		Description: "获取当前的日期和时间",
		ParamsJSONSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			now := time.Now()
			return fmt.Sprintf("当前时间是: %s", now.Format("2006-01-02 15:04:05")), nil
		},
	}
}

// 定义工具 2: 计算器
func createCalculatorTool() agentgo.FunctionTool {
	return agentgo.FunctionTool{
		Name:        "calculator",
		Description: "执行数学计算。支持加减乘除和幂运算。",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"operation": map[string]any{
					"type":        "string",
					"enum":        []string{"add", "subtract", "multiply", "divide", "power"},
					"description": "要执行的运算类型",
				},
				"a": map[string]any{
					"type":        "number",
					"description": "第一个操作数",
				},
				"b": map[string]any{
					"type":        "number",
					"description": "第二个操作数",
				},
			},
			"required": []string{"operation", "a", "b"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				Operation string  `json:"operation"`
				A         float64 `json:"a"`
				B         float64 `json:"b"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			var result float64
			switch params.Operation {
			case "add":
				result = params.A + params.B
			case "subtract":
				result = params.A - params.B
			case "multiply":
				result = params.A * params.B
			case "divide":
				if params.B == 0 {
					return "错误: 除数不能为零", nil
				}
				result = params.A / params.B
			case "power":
				result = math.Pow(params.A, params.B)
			default:
				return fmt.Sprintf("未知操作: %s", params.Operation), nil
			}

			return fmt.Sprintf("%v %s %v = %v", params.A, params.Operation, params.B, result), nil
		},
	}
}

// 定义工具 3: 天气查询（模拟）
func createWeatherTool() agentgo.FunctionTool {
	return agentgo.FunctionTool{
		Name:        "get_weather",
		Description: "查询指定城市的天气信息",
		ParamsJSONSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"city": map[string]any{
					"type":        "string",
					"description": "要查询天气的城市名称",
				},
			},
			"required": []string{"city"},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
			var params struct {
				City string `json:"city"`
			}
			if err := json.Unmarshal([]byte(arguments), &params); err != nil {
				return nil, fmt.Errorf("参数解析失败: %w", err)
			}

			// 模拟天气数据
			weatherData := map[string]string{
				"北京": "晴天，温度 15°C，湿度 45%",
				"上海": "多云，温度 18°C，湿度 60%",
				"广州": "小雨，温度 22°C，湿度 80%",
				"深圳": "阴天，温度 20°C，湿度 70%",
			}

			if weather, ok := weatherData[params.City]; ok {
				return fmt.Sprintf("%s天气: %s", params.City, weather), nil
			}
			return fmt.Sprintf("抱歉，暂无 %s 的天气信息", params.City), nil
		},
	}
}

func main() {
	// 检查 API Key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请设置 OPENAI_API_KEY 环境变量")
	}

	// 创建 OpenAI 客户端
	client := openai.NewClient()

	// 创建带工具的 Agent
	agent := agentgo.New("工具助手").
		WithInstructions(`你是一个智能助手，可以使用以下工具帮助用户:
1. get_current_time: 获取当前时间
2. calculator: 进行数学计算
3. get_weather: 查询天气

请根据用户需求选择合适的工具。`).
		WithModel("gpt-4o-mini").
		WithClient(client).
		WithTools([]agentgo.FunctionTool{
			createGetTimeTool(),
			createCalculatorTool(),
			createWeatherTool(),
		})

	// 配置 Runner
	runner := agentgo.Runner{
		Config: agentgo.RunConfig{
			MaxTurns: 5, // 最多执行 5 轮
		},
	}

	ctx := context.Background()

	// 测试用例
	testCases := []string{
		"现在几点了？",
		"帮我计算 123 乘以 456",
		"北京今天天气怎么样？",
		"2 的 10 次方是多少？",
	}

	for i, question := range testCases {
		fmt.Printf("\n=== 测试 %d: %s ===\n", i+1, question)

		result, err := runner.Run(ctx, agent, question)
		if err != nil {
			log.Printf("运行失败: %v", err)
			continue
		}

		fmt.Printf("回复: %v\n", result.FinalOutput)

		// 显示工具调用信息
		if len(result.NewItems) > 0 {
			fmt.Printf("工具调用次数: %d\n", len(result.NewItems))
		}
	}
}
