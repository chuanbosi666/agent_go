package pattern

import (
	"context"
	"encoding/json"
	"fmt"

	"nvgo/pkg/agent"
	"nvgo/pkg/runner"
	"nvgo/pkg/tool"
)

// DefaultMaxTurns for agent-as-tool execution.
const DefaultMaxTurns = 10

// WrapAgentAsTool converts an Agent into a callable FunctionTool.
// Enables multi-agent collaboration by allowing agents to call other agents.
func WrapAgentAsTool(a *agent.Agent, maxTurns uint64) tool.FunctionTool {
	if maxTurns == 0 {
			maxTurns = DefaultMaxTurns
	}

	return tool.FunctionTool{
			Name:        fmt.Sprintf("call_agent_%s", a.Name),
			Description: fmt.Sprintf("Call agent '%s' to handle a task. %s", a.Name, getAgentDescription(a)),
			ParamsJSONSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
							"input": map[string]any{
									"type":        "string",
									"description": "The input message to pass to the agent",
							},
					},
					"required": []string{"input"},
			},
			OnInvokeTool: func(ctx context.Context, arguments string) (any, error) {
					var params struct {
							Input string `json:"input"`
					}
					if err := json.Unmarshal([]byte(arguments), &params); err != nil {
							return nil, fmt.Errorf("unmarshal arguments: %w", err)
					}

					r := runner.Runner{Config: runner.RunConfig{
							MaxTurns: maxTurns,
					}}

					result, err := r.Run(ctx, a, params.Input)
					if err != nil {
							return nil, fmt.Errorf("run agent %q: %w", a.Name, err)
					}

					return result.FinalOutput, nil
			},
	}
}

// getAgentDescription extracts description from agent instructions.
func getAgentDescription(a *agent.Agent) string {
	if a.Instructions == nil {
			return fmt.Sprintf("Delegate tasks to the %s agent", a.Name)
	}

	if instrStr, ok := a.Instructions.(agent.InstructionsStr); ok {
			str := instrStr.String()
			if len(str) > 100 {
					return str[:100] + "..."
			}
			return str
	}
	return fmt.Sprintf("Delegate tasks to the %s agent", a.Name)
}