package nvgo

import (
	"context"
	"encoding/json"
	"fmt"
)

func WrapAgentAsTool(agent *Agent, maxTurns uint64) FunctionTool {
	if maxTurns == 0 {
		maxTurns = DefaultMaxTurns
	}

	return FunctionTool{
		Name: fmt.Sprintf("Call_agent_%s", agent.Name),
		Description: fmt.Sprintf("Call_agent '%s' to handle a specific task. %s",
			agent.Name,
			getAgentDescription(agent)),
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

			runner := Runner{Config: RunConfig{
				MaxTurns: maxTurns,
			}}

			result, err := runner.Run(ctx, agent, params.Input)
			if err != nil {
				return nil, fmt.Errorf("run agent %q: %w", agent.Name, err)
			}

			return result.FinalOutput, err

		},
	}
}

func getAgentDescription(agent *Agent) string {
	if agent.Instructions == nil {
		return fmt.Sprintf("Use this to delegate tasks to the %s agent", agent.Name)
	}

	if instrStr, ok := agent.Instructions.(InstructionsStr); ok {
		str := instrStr.String()

		if len(str) > 100 {
			return str[:100] + "..."
		}
		return str
	}
	return fmt.Sprintf("use this to delegate tasks to the %s agent", agent.Name)
}
