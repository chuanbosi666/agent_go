// Package agent - setting.go provides model configuration parameters.
package agent

import (
	"context"
	"maps"
	"reflect"

	"github.com/openai/openai-go/v3"    
	"github.com/openai/openai-go/v3/option"
    "github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

// ModelSettings holds LLM configuration parameters.
// Not all models/providers support all parameters.
type ModelSettings struct{
	// Temperature controls randomness (0.0 = deterministic, 2.0 = very random).
	Temperature param.Opt[float64] `json:"temperature"`

	// TopP controls nucleus sampling (alternative to temperature).
    TopP param.Opt[float64] `json:"top_p"`

	// FrequencyPenalty reduces repetition of token sequences (-2.0 to 2.0).
    FrequencyPenalty param.Opt[float64] `json:"frequency_penalty"`

	// PresencePenalty reduces repetition of topics (-2.0 to 2.0).
    PresencePenalty param.Opt[float64] `json:"presence_penalty"`

	// ToolChoice controls which tool the model should use.
    ToolChoice ToolChoice `json:"tool_choice"`

	// ParallelToolCalls enables multiple tool calls in a single turn.
    ParallelToolCalls param.Opt[bool] `json:"parallel_tool_calls"`

	// Truncation strategy for long conversations.
    Truncation param.Opt[Truncation] `json:"truncation"`

	// MaxTokens is the maximum number of output tokens.
    MaxTokens param.Opt[int64] `json:"max_tokens"`

	// Reasoning configures reasoning model behavior (o1, o3, etc.).
    Reasoning openai.ReasoningParam `json:"reasoning"`
	
	// Metadata is optional key-value pairs included with the request.
    Metadata map[string]string `json:"metadata"`

	// Store enables storage of responses for later retrieval.
    Store param.Opt[bool] `json:"store"`

	// IncludeUsage includes token usage in response.
    IncludeUsage param.Opt[bool] `json:"include_usage"`

	// ResponseInclude specifies additional data to include in response.
    ResponseInclude []responses.ResponseIncludable `json:"response_include"`

	// ExtraQuery adds custom query parameters to the request.
    ExtraQuery map[string]string `json:"extra_query"`
	// ExtraHeaders adds custom headers to the request.
    ExtraHeaders map[string]string `json:"extra_headers"`
	
	// CustomizeResponsesRequest allows full customization of Responses API calls.
    CustomizeResponsesRequest func(context.Context, *responses.ResponseNewParams, []option.RequestOption) (*responses.ResponseNewParams, []option.RequestOption, error) `json:"-"`

    // CustomizeChatCompletionsRequest allows full customization of Chat API calls.
    CustomizeChatCompletionsRequest func(context.Context, *openai.ChatCompletionNewParams,
  []option.RequestOption) (*openai.ChatCompletionNewParams, []option.RequestOption, error) `json:"-"`
}

// ToolChoice specifies which tool the model should use.
type ToolChoice interface {
        isToolChoice()
}

// ToolChoiceString is a string-based tool choice.
type ToolChoiceString string

func (ToolChoiceString) isToolChoice()     {}
func (tc ToolChoiceString) String() string { return string(tc) }

// Predefined tool choice values.
const (
    ToolChoiceAuto     ToolChoiceString = "auto"     // Model decides
    ToolChoiceRequired ToolChoiceString = "required" // Must use a tool
    ToolChoiceNone     ToolChoiceString = "none"     // No tool use
)

type Truncation string

const (
	TruncationAuto     Truncation = "auto"
	TruncationDisabled Truncation = "disabled"
)

 // Resolve merges override settings into this instance.
  // Non-zero values in override take precedence.
func (ms ModelSettings) Resolve(override ModelSettings) ModelSettings {
        newSettings := ms
        resolveOpt(&newSettings.Temperature, override.Temperature)
        resolveOpt(&newSettings.TopP, override.TopP)
        resolveOpt(&newSettings.FrequencyPenalty, override.FrequencyPenalty)
        resolveOpt(&newSettings.PresencePenalty, override.PresencePenalty)
        resolveAny(&newSettings.ToolChoice, override.ToolChoice)
        resolveOpt(&newSettings.ParallelToolCalls, override.ParallelToolCalls)
        resolveOpt(&newSettings.Truncation, override.Truncation)
        resolveOpt(&newSettings.MaxTokens, override.MaxTokens)
        resolveAny(&newSettings.Reasoning, override.Reasoning)
        resolveMap(&newSettings.Metadata, override.Metadata)
        resolveOpt(&newSettings.Store, override.Store)
        resolveOpt(&newSettings.IncludeUsage, override.IncludeUsage)
        resolveAny(&newSettings.ResponseInclude, override.ResponseInclude)
        resolveMap(&newSettings.ExtraQuery, override.ExtraQuery)
        resolveMap(&newSettings.ExtraHeaders, override.ExtraHeaders)
        resolveAny(&newSettings.CustomizeResponsesRequest, override.CustomizeResponsesRequest)
        resolveAny(&newSettings.CustomizeChatCompletionsRequest, override.CustomizeChatCompletionsRequest)
        return newSettings
}

// resolveOpt overwrites base with override if override is valid.
func resolveOpt[T comparable](base *param.Opt[T], override param.Opt[T]) {
        if override.Valid() {
                *base = override
        }
}

// resolveAny overwrites base with override if override is non-zero.
func resolveAny[T any](base *T, override T) {
    v := reflect.ValueOf(override)
    if v.Kind() != reflect.Invalid && !v.IsZero() {
            *base = override
    }
}

// resolveMap overwrites base with a clone of override if override is non-empty.
  func resolveMap[M ~map[K]V, K comparable, V any](base *M, override M) {
    if len(override) > 0 {
        *base = maps.Clone(override)
    }
}