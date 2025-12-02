package nvgo

import (
	"context"
	"sort"
	"strings"
)

type ToolRouter interface {
	RouteTools(ctx context.Context, input Input, tool []Tool) ([]Tool, error)
}

type KeywordRouter struct {
	ToolKeywords map[string][]string
	TopN         int
}

func (r *KeywordRouter) RouteTools(ctx context.Context, input Input, tools []Tool) ([]Tool, error) {
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

func extractText(input Input) string {
	switch v := input.(type) {
	case InputString:
		return string(v)
	case InputItems:
		var texts []string
		for _, item := range v {
			if item.OfMessage != nil {
				texts = append(texts, item.OfMessage.Content.OfString.String())
			}
		}
		return strings.Join(texts, " ")
	}
	return ""
}
