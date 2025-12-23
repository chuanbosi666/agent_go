package tool

import (
	"context"
	"sort"
	"strings"

	"github.com/chuanbosi666/agent_go/pkg/types"
)

// ToolRouter selects relevant tools based on user input.
type ToolRouter interface {
	RouteTools(ctx context.Context, input types.Input, tools []Tool) ([]Tool, error)
}

// KeywordRouter routes tools by matching keywords in user input.
type KeywordRouter struct {
	ToolKeywords map[string][]string // Map of tool name to associated keywords
	TopN         int                 // Max number of tools to return (default: 5)
}

// RouteTools scores and returns top N tools based on keyword matches.
func (r *KeywordRouter) RouteTools(ctx context.Context, input types.Input, tools []Tool) ([]Tool, error) {
	inputText := extractText(input)
	inputLower := strings.ToLower(inputText)

	type toolScore struct {
		tool  Tool
		score int
	}

	var scored []toolScore

	for _, tool := range tools {
		toolName := tool.ToolName()
		keywords := r.ToolKeywords[toolName]

		score := 0
		for _, keyword := range keywords {
			if strings.Contains(inputLower, strings.ToLower(keyword)) {
				score++
			}
		}
		scored = append(scored, toolScore{tool: tool, score: score})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	topN := r.TopN
	if topN <= 0 {
		topN = 5
	}
	if topN > len(scored) {
		topN = len(scored)
	}

	result := make([]Tool, topN)
	for i := 0; i < topN; i++ {
		result[i] = scored[i].tool
	}
	return result, nil
}

// extractText converts Input to plain text for keyword matching.
func extractText(input types.Input) string {
	if input == nil {
		return ""
	}
	items := input.ToInputItems()
	var texts []string
	for _, item := range items {
		if item.OfMessage != nil {
			texts = append(texts, item.OfMessage.Content.OfString.String())
		}
	}
	return strings.Join(texts, " ")
}
