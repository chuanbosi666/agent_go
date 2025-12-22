package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// NewTimeTool creates a tool that returns the current time.
func NewTimeTool() FunctionTool{
	return FunctionTool{
		Name:		"get_current_time",
		Description: "Get the current date and time",
		ParamsJSONSchema: map[string]any{
			"type":"object",
			"properties": map[string]any{
				"timezone" : map[string]any{
					"type":			"string",
					"description":	"Timezone name (e.g., 'Asia/Shanghai', 'UTC'). Default is local timezone.",
				},
			},
		},
		OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {  // 完整签名
                        var params struct {
                                Timezone string `json:"timezone"`
                        }
                        if arguments != "" {
                                json.Unmarshal([]byte(arguments), &params)
                        }

                        var loc *time.Location
                        var err error
                        if params.Timezone != "" {
                                loc, err = time.LoadLocation(params.Timezone)
                                if err != nil {
                                        return nil, fmt.Errorf("invalid timezone: %s", params.Timezone)
                                }
                        } else {
                                loc = time.Local
                        }

                        now := time.Now().In(loc)
                        return fmt.Sprintf("Current time: %s", now.Format("2006-01-02 15:04:05 MST")), nil       
                },
        }
  }
// NewWebFetchTool creates a tool that fetches content from a URL.
func NewWebFetchTool() FunctionTool {
	return FunctionTool{
			Name:        "web_fetch",
			Description: "Fetch content from a web URL. Returns the raw text content.",
			ParamsJSONSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
							"url": map[string]any{
									"type":        "string",
									"description": "The URL to fetch content from",
							},
					},
					"required": []string{"url"},
			},
			OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
					var params struct {
							URL string `json:"url"`
					}
					if err := json.Unmarshal([]byte(arguments), &params); err != nil {
							return nil, fmt.Errorf("invalid arguments: %w", err)
					}

					if params.URL == "" {
							return nil, fmt.Errorf("url is required")
					}

					// Validate URL
					if _, err := url.Parse(params.URL); err != nil {
							return nil, fmt.Errorf("invalid url: %w", err)
					}

					// Create request with timeout
					req, err := http.NewRequestWithContext(ctx, "GET", params.URL, nil)
					if err != nil {
							return nil, fmt.Errorf("create request: %w", err)
					}
					req.Header.Set("User-Agent", "NVGo-Agent/1.0")

					client := &http.Client{Timeout: 30 * time.Second}
					resp, err := client.Do(req)
					if err != nil {
							return nil, fmt.Errorf("fetch failed: %w", err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
							return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
					}

					// Limit response size to 100KB
					body, err := io.ReadAll(io.LimitReader(resp.Body, 100*1024))
					if err != nil {
							return nil, fmt.Errorf("read body: %w", err)
					}

					return string(body), nil
			},
	}
}

// NewCalculatorTool creates a simple calculator tool.
func NewCalculatorTool() FunctionTool {
	return FunctionTool{
			Name:        "calculator",
			Description: "Perform basic math operations: add, subtract, multiply, divide",
			ParamsJSONSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
							"operation": map[string]any{
									"type":        "string",
									"enum":        []string{"add", "subtract", "multiply", "divide"},        
									"description": "The math operation to perform",
							},
							"a": map[string]any{
									"type":        "number",
									"description": "First operand",
							},
							"b": map[string]any{
									"type":        "number",
									"description": "Second operand",
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
							return nil, fmt.Errorf("invalid arguments: %w", err)
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
									return "Error: division by zero", nil
							}
							result = params.A / params.B
					default:
							return nil, fmt.Errorf("unknown operation: %s", params.Operation)
					}

					return fmt.Sprintf("%.6g", result), nil
			},
	}
}