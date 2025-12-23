package runner

import (
	"context"
	"fmt"

	"github.com/chuanbosi666/agent_go/pkg/agent"
	"github.com/chuanbosi666/agent_go/pkg/memory"
	"github.com/chuanbosi666/agent_go/pkg/tool"
	"github.com/chuanbosi666/agent_go/pkg/types"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

// RunItem represents an item generated during Agent execution.
type RunItem interface {
	isRunItem()
	ToInputItem() responses.ResponseInputItemUnionParam
}

type RunItemWrapper struct {
	item responses.ResponseInputItemUnionParam
}

func (w RunItemWrapper) isRunItem() {}
func (w RunItemWrapper) ToInputItem() responses.ResponseInputItemUnionParam {
	return w.item
}

func WrapRunItem(item responses.ResponseInputItemUnionParam) RunItem {
	return RunItemWrapper{item: item}
}

// Output represents the final output from an Agent run.
type Output struct {
	Items        []responses.ResponseOutputItemUnion
	InputTokens  int64
	OutputTokens int64
}

// Usage tracks token consumption across LLM requests.
type Usage struct {
	Requests            uint64
	InputTokens         uint64
	InputTokensDetails  responses.ResponseUsageInputTokensDetails
	OutputTokens        uint64
	OutputTokensDetails responses.ResponseUsageOutputTokensDetails
	TotalTokens         uint64
}

// ModelResponse holds a single LLM response.
type ModelResponse struct {
	Output     []responses.ResponseOutputItemUnion
	Usage      *Usage
	ResponseID string
}

// RunResult contains the complete results of an Agent execution.
type RunResult struct {
	Input                  types.Input
	NewItems               []RunItem
	RawResponses           []ModelResponse
	FinalOutput            any
	InputGuardrailResults  []agent.InputGuardrailResult
	OutputGuardrailResults []agent.OutputGuardrailResult
	LastAgent              *agent.Agent
}

const DefaultMaxTurns = 10
const DefaultWorkflowName = "Agent workflow"

var DefaultRunner = Runner{}

// Runner executes Agents. The zero value is valid.
type Runner struct {
	Config RunConfig
}

// RunConfig configures Agent execution behavior.
type RunConfig struct {
	Model                string
	ModelSettings        agent.ModelSettings
	InputGuardrails      []agent.InputGuardrail
	OutputGuardrails     []agent.OutputGuardrail
	WorkflowName         string
	MaxTurns             uint64
	PreviousResponseID   string
	Session              memory.Session
	HandoffResolver      func(ctx context.Context, agentName string) (*agent.Agent, error)
	ToolRouter           tool.ToolRouter
	ToolRoutingThreshold int
}

func (o Output) TotalTokens() int64 {
	return o.InputTokens + o.OutputTokens
}

// Run executes the agent with a string input using DefaultRunner.
func Run(ctx context.Context, startingAgent *agent.Agent, input string) (*RunResult, error) {
	return DefaultRunner.Run(ctx, startingAgent, input)
}

// Run executes the agent with a string input.
func (r Runner) Run(ctx context.Context, startingAgent *agent.Agent, input string) (*RunResult, error) {
	return r.run(ctx, startingAgent, types.InputString(input))
}

// MaxTurnsExceededError is returned when execution exceeds MaxTurns.
type MaxTurnsExceededError struct {
	MaxTurns uint64
}

func (e *MaxTurnsExceededError) Error() string {
	return fmt.Sprintf("max turns exceeded: reached limit of %d turns", e.MaxTurns)
}

// GuardrailTripwireTriggeredError is returned when a guardrail blocks execution.
type GuardrailTripwireTriggeredError struct {
	GuardrailName string
	OutputInfo    any
	IsInput       bool
}

func (g *GuardrailTripwireTriggeredError) Error() string {
	if g.IsInput {
		return fmt.Sprintf("input guardrail '%s' triggered", g.GuardrailName)
	}
	return fmt.Sprintf("output guardrail '%s' triggered", g.GuardrailName)
}

// run is the core execution loop.
func (r Runner) run(ctx context.Context, startingAgent *agent.Agent, input types.Input) (*RunResult, error) {
	result := &RunResult{
		Input:        types.CopyInput(input),
		NewItems:     []RunItem{},
		RawResponses: []ModelResponse{},
	}

	// Run input guardrails
	inputGuardrails := append(r.Config.InputGuardrails, startingAgent.InputGuardrails...)
	for _, gr := range inputGuardrails {
		grResult, err := gr.Run(ctx, startingAgent, input)
		if err != nil {
			return nil, fmt.Errorf("input guardrail %q failed: %w", gr.Name, err)
		}
		result.InputGuardrailResults = append(result.InputGuardrailResults, grResult)
		if grResult.Output.TripwireTriggered {
			return nil, &GuardrailTripwireTriggeredError{
				GuardrailName: gr.Name,
				OutputInfo:    grResult.Output.OutputInfo,
				IsInput:       true,
			}
		}
	}

	currentAgent := startingAgent
	turnCount := uint64(0)
	maxTurns := r.Config.MaxTurns
	if maxTurns == 0 {
		maxTurns = DefaultMaxTurns
	}

	var accumulatedHistory []responses.ResponseInputItemUnionParam

	// Main execution loop
	for turnCount < maxTurns {
		turnCount++

		model := r.Config.Model
		if model == "" {
			model = currentAgent.Model
		}

		// Get instructions
		var instructions string
		if currentAgent.Instructions != nil {
			var err error
			instructions, err = currentAgent.Instructions.GetInstructions(ctx, currentAgent)
			if err != nil {
				return nil, fmt.Errorf("get instruction: %w", err)
			}
		}

		// Get tools (with optional routing)
		tools, err := getAgentTools(ctx, currentAgent, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get MCP tools: %w", err)
		}

		threshold := r.Config.ToolRoutingThreshold
		if threshold == 0 {
			threshold = 5
		}
		if r.Config.ToolRouter != nil && len(tools) > threshold {
			routedTools, routeErr := r.Config.ToolRouter.RouteTools(ctx, input, tools)
			if routeErr == nil {
				tools = routedTools
			}
		}

		modelsettings := currentAgent.ModelSettings.Resolve(r.Config.ModelSettings)

		// Load session history
		var historyItems []responses.ResponseInputItemUnionParam
		if r.Config.Session != nil {
			items, err := r.Config.Session.GetItems(ctx, -1)
			if err != nil {
				return nil, fmt.Errorf("load session history: %w", err)
			}
			historyItems = items
		}

		var modelResponse ModelResponse

		// Choose API path: Responses API or Chat Completions API
		if currentAgent.Prompt != nil {
			// Responses API path (OpenAI only)
			modelResponse, err = r.callResponsesAPI(ctx, currentAgent, model, instructions, tools, modelsettings, historyItems, accumulatedHistory, input, turnCount)
			if err != nil {
				return nil, err
			}
		} else {
			// Chat Completions API path (OpenAI-compatible)
			modelResponse, err = r.callChatCompletionsAPI(ctx, currentAgent, model, instructions, modelsettings, input, turnCount, result)
			if err != nil {
				return nil, err
			}
		}

		result.RawResponses = append(result.RawResponses, modelResponse)

		// Process tool calls
		for _, outputItem := range modelResponse.Output {
			switch item := outputItem.AsAny().(type) {
			case responses.ResponseOutputMessage:
				// Message output handled below
			case responses.ResponseFunctionToolCall:
				t, found := FindTool(tools, item.Name)
				if !found {
					errorOutput := responses.ResponseInputItemParamOfFunctionCallOutput(
						item.CallID,
						fmt.Sprintf("Tool %s not found", item.Name),
					)
					result.NewItems = append(result.NewItems, WrapRunItem(errorOutput))
					continue
				}

				toolResult, err := executeTool(ctx, currentAgent, t, item.Arguments)
				if err != nil {
					errorOutput := responses.ResponseInputItemParamOfFunctionCallOutput(
						item.CallID,
						fmt.Sprintf("Tool execution failed: %v", err),
					)
					result.NewItems = append(result.NewItems, WrapRunItem(errorOutput))
					continue
				}

				var outputStr string
				switch v := toolResult.(type) {
				case string:
					outputStr = v
				default:
					outputStr = fmt.Sprintf("%v", v)
				}
				successOutput := responses.ResponseInputItemParamOfFunctionCallOutput(
					item.CallID,
					outputStr,
				)
				result.NewItems = append(result.NewItems, WrapRunItem(successOutput))
			}
		}

		// Extract final output from messages
		for _, outputItem := range modelResponse.Output {
			if msg, ok := outputItem.AsAny().(responses.ResponseOutputMessage); ok {
				result.FinalOutput = msg.Content
				break
			}
		}

		// Save tool results to session/history
		if len(result.NewItems) > 0 {
			var itemsToSave []responses.ResponseInputItemUnionParam
			for _, item := range result.NewItems {
				itemsToSave = append(itemsToSave, item.ToInputItem())
			}
			if r.Config.Session != nil {
				if err := r.Config.Session.AddItems(ctx, itemsToSave); err != nil {
					return nil, fmt.Errorf("save tool results to session: %w", err)
				}
			} else {
				accumulatedHistory = append(accumulatedHistory, itemsToSave...)
			}
		}

		// Save model output to session/history
		if len(modelResponse.Output) > 0 {
			var modelOutputItems []responses.ResponseInputItemUnionParam
			for _, outputItem := range modelResponse.Output {
				switch item := outputItem.AsAny().(type) {
				case responses.ResponseOutputMessage:
					outputParm := responses.ResponseOutputMessageParam{
						ID:      item.ID,
						Content: convertOutputContentToParam(item.Content),
						Status:  item.Status,
						Role:    item.Role,
						Type:    item.Type,
					}
					inputItem := responses.ResponseInputItemUnionParam{
						OfOutputMessage: &outputParm,
					}
					modelOutputItems = append(modelOutputItems, inputItem)
				case responses.ResponseFunctionToolCall:
					if r.Config.Session == nil {
						funcCallParam := responses.ResponseInputItemParamOfFunctionCall(
							item.Arguments,
							item.CallID,
							item.Name,
						)
						modelOutputItems = append(modelOutputItems, funcCallParam)
					}
				}
			}
			if len(modelOutputItems) > 0 {
				if r.Config.Session != nil {
					if err := r.Config.Session.AddItems(ctx, modelOutputItems); err != nil {
						return nil, fmt.Errorf("save model output to session: %w", err)
					}
				} else {
					accumulatedHistory = append(accumulatedHistory, modelOutputItems...)
				}
			}
		}

		if result.FinalOutput != nil {
			break
		}
	}

	if turnCount >= maxTurns {
		return nil, &MaxTurnsExceededError{MaxTurns: maxTurns}
	}

	if result.FinalOutput == nil {
		result.LastAgent = currentAgent
		return result, nil
	}

	// Run output guardrails
	outputGuardrails := append(r.Config.OutputGuardrails, currentAgent.OutputGuardrails...)
	for _, gr := range outputGuardrails {
		grResult, err := gr.Run(ctx, currentAgent, result.FinalOutput)
		if err != nil {
			return nil, fmt.Errorf("output guardrail %q failed: %w", gr.Name, err)
		}
		result.OutputGuardrailResults = append(result.OutputGuardrailResults, grResult)
		if grResult.Output.TripwireTriggered {
			return nil, &GuardrailTripwireTriggeredError{
				GuardrailName: gr.Name,
				OutputInfo:    grResult.Output.OutputInfo,
				IsInput:       false,
			}
		}
	}

	result.LastAgent = currentAgent
	return result, nil
}

// callResponsesAPI calls the OpenAI Responses API.
func (r Runner) callResponsesAPI(
	ctx context.Context,
	currentAgent *agent.Agent,
	model, instructions string,
	tools []tool.Tool,
	modelsettings agent.ModelSettings,
	historyItems, accumulatedHistory []responses.ResponseInputItemUnionParam,
	input types.Input,
	turnCount uint64,
) (ModelResponse, error) {
	promptParam, hasPrompt, err := agent.PromptUtil().ToModelInput(ctx, currentAgent.Prompt, currentAgent)
	if err != nil {
		return ModelResponse{}, fmt.Errorf("get prompt: %w", err)
	}
	if !hasPrompt {
		return ModelResponse{}, fmt.Errorf("prompt is required but not provided")
	}

	var allInputItems []responses.ResponseInputItemUnionParam
	if len(historyItems) > 0 {
		allInputItems = append(allInputItems, historyItems...)
	} else {
		allInputItems = append(allInputItems, accumulatedHistory...)
	}

	if turnCount == 1 {
		currentInputItems := InputToItems(input)
		allInputItems = append(allInputItems, currentInputItems...)
	}

	toolParams := ToolsToParams(tools)

	createParams := responses.ResponseNewParams{
		Model:  model,
		Prompt: promptParam,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam(allInputItems),
		},
	}

	if instructions != "" {
		createParams.Instructions = param.NewOpt(instructions)
	}
	if len(toolParams) > 0 {
		createParams.Tools = toolParams
	}
	if modelsettings.Temperature.Valid() {
		createParams.Temperature = modelsettings.Temperature
	}
	if modelsettings.MaxTokens.Valid() {
		createParams.MaxOutputTokens = modelsettings.MaxTokens
	}
	if modelsettings.TopP.Valid() {
		createParams.TopP = modelsettings.TopP
	}

	resp, err := currentAgent.Client.Responses.New(ctx, createParams)
	if err != nil {
		return ModelResponse{}, fmt.Errorf("call responses API: %w", err)
	}

	return ModelResponse{
		Output:     resp.Output,
		ResponseID: resp.ID,
		Usage: &Usage{
			Requests:            1,
			InputTokens:         uint64(resp.Usage.InputTokens),
			InputTokensDetails:  resp.Usage.InputTokensDetails,
			OutputTokens:        uint64(resp.Usage.OutputTokens),
			OutputTokensDetails: resp.Usage.OutputTokensDetails,
			TotalTokens:         uint64(resp.Usage.TotalTokens),
		},
	}, nil
}

// callChatCompletionsAPI calls the OpenAI-compatible Chat Completions API.
func (r Runner) callChatCompletionsAPI(
	ctx context.Context,
	currentAgent *agent.Agent,
	model, instructions string,
	modelsettings agent.ModelSettings,
	input types.Input,
	turnCount uint64,
	result *RunResult,
) (ModelResponse, error) {
	var messages []openai.ChatCompletionMessageParamUnion

	if instructions != "" {
		messages = append(messages, openai.SystemMessage(instructions))
	}

	if turnCount == 1 {
		switch v := input.(type) {
		case types.InputString:
			messages = append(messages, openai.UserMessage(string(v)))
		}
	}

	chatParams := openai.ChatCompletionNewParams{
		Model:    model,
		Messages: messages,
	}
	if modelsettings.Temperature.Valid() {
		chatParams.Temperature = modelsettings.Temperature
	}
	if modelsettings.MaxTokens.Valid() {
		chatParams.MaxTokens = modelsettings.MaxTokens
	}
	if modelsettings.TopP.Valid() {
		chatParams.TopP = modelsettings.TopP
	}

	chatresp, err := currentAgent.Client.Chat.Completions.New(ctx, chatParams)
	if err != nil {
		return ModelResponse{}, fmt.Errorf("call chat completions API: %w", err)
	}

	modelResponse := ModelResponse{
		Output: []responses.ResponseOutputItemUnion{},
		Usage: &Usage{
			Requests:     1,
			InputTokens:  uint64(chatresp.Usage.PromptTokens),
			OutputTokens: uint64(chatresp.Usage.CompletionTokens),
			TotalTokens:  uint64(chatresp.Usage.TotalTokens),
		},
	}

	if len(chatresp.Choices) > 0 {
		choice := chatresp.Choices[0]
		if choice.Message.Content != "" {
			result.FinalOutput = choice.Message.Content
		}
	}

	return modelResponse, nil
}

// Helper functions

func getAgentTools(ctx context.Context, a *agent.Agent, strict bool) ([]tool.Tool, error) {
	var allTools []tool.Tool

	if len(a.MCPServers) > 0 {
		mcpTools, err := tool.GetAllFunctionTools(ctx, a.MCPServers, strict, a)
		if err != nil {
			return nil, err
		}
		allTools = append(allTools, mcpTools...)
	}

	for _, t := range a.Tools {
		allTools = append(allTools, t)
	}

	return allTools, nil
}

func FindTool(tools []tool.Tool, name string) (tool.Tool, bool) {
	for _, t := range tools {
		if t.ToolName() == name {
			return t, true
		}
	}
	return nil, false
}

func executeTool(ctx context.Context, a *agent.Agent, t tool.Tool, arguments string) (any, error) {
	funcTool, ok := t.(tool.FunctionTool)
	if !ok {
		return nil, fmt.Errorf("tool is not a FunctionTool")
	}

	if funcTool.IsEnabled != nil {
		enabled, err := funcTool.IsEnabled.IsEnabled(ctx, a)
		if err != nil {
			return nil, fmt.Errorf("check tool enabled: %w", err)
		}
		if !enabled {
			return nil, fmt.Errorf("tool %s is disabled", funcTool.ToolName())
		}
	}

	result, err := funcTool.OnInvokeTool(ctx, arguments)
	if err != nil {
		if funcTool.FailureErrorFunction != nil {
			errorFunc := *funcTool.FailureErrorFunction
			val, _ := errorFunc(ctx, err)
			return val, nil
		}
		val, _ := tool.DefaultToolErrorFunction(ctx, err)
		return val, nil
	}
	return result, nil
}

func InputToItems(input types.Input) []responses.ResponseInputItemUnionParam {
	switch v := input.(type) {
	case types.InputString:
		return []responses.ResponseInputItemUnionParam{
			responses.ResponseInputItemParamOfMessage(
				string(v),
				responses.EasyInputMessageRole(responses.ResponseInputMessageItemRoleUser)),
		}
	case types.InputItems:
		return []responses.ResponseInputItemUnionParam(v)
	default:
		panic(fmt.Errorf("unexpected Input type %T", v))
	}
}

func ToolsToParams(tools []tool.Tool) []responses.ToolUnionParam {
	if len(tools) == 0 {
		return nil
	}

	params := make([]responses.ToolUnionParam, 0, len(tools))
	for _, t := range tools {
		funcTool, ok := t.(tool.FunctionTool)
		if !ok {
			continue
		}
		toolParam := responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        funcTool.Name,
				Description: param.NewOpt(funcTool.Description),
				Parameters:  funcTool.ParamsJSONSchema,
				Strict:      funcTool.StrictJSONSchema,
			},
		}
		params = append(params, toolParam)
	}
	return params
}

func convertOutputContentToParam(content []responses.ResponseOutputMessageContentUnion) []responses.ResponseOutputMessageContentUnionParam {
	var result []responses.ResponseOutputMessageContentUnionParam
	for _, c := range content {
		switch item := c.AsAny().(type) {
		case responses.ResponseOutputText:
			textParam := item.ToParam()
			param := responses.ResponseOutputMessageContentUnionParam{
				OfOutputText: &textParam,
			}
			result = append(result, param)
		}
	}
	return result
}
