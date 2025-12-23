package pattern

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/chuanbosi666/agent_go/pkg/agent"
)

// DefaultReActInstruction provides standard ReAct prompting with structured reasoning.
var DefaultReActInstruction = agent.InstructionsStr(`
You are a strong reasoner using the ReAct pattern to solve problems methodically.

## Reasoning Framework
Before each action, reason through:
1. Logical dependencies: Analyze constraints, order of operations, prerequisites
2. Risk assessment: Consider consequences and potential issues
3. Hypothesis exploration: Identify likely causes, test systematically
4. Information sources: Use tools, history, policies, and user input
5. Completeness: Ensure all requirements are addressed

## Response Format
Each response must follow this structure:

Thought: [Analyze situation using the reasoning framework above]
Action: [Tool name to execute]
Action Input: [Tool parameters in JSON format]

After observing results:

Observation: [Results returned by the tool]
Thought: [Evaluate outcome, adapt plan if needed, generate new hypotheses if disproven]

When task is complete:

Thought: [Summarize analysis and verify all requirements met]
Final Answer: [Final response to the user]

## Rules
1. Think step by step with logical precision
2. Execute only one tool at a time
3. Adapt plan based on Observation - don't repeat failed strategies
4. If hypothesis disproven, generate new ones from gathered information
5. Be persistent - exhaust all reasoning before giving up
6. On transient errors, retry with modified approach
7. Quote exact information when making claims
8. Must provide Final Answer when complete
`)

// ReActStateProvider tracks ReAct execution state.
type ReActStateProvider struct {
	currentStep  int
	observations []string
	mu           sync.RWMutex
}

// NewReActInstruction creates instruction with custom rules appended.
func NewReActInstruction(customRules string) agent.InstructionsStr {
	return DefaultReActInstruction + agent.InstructionsStr(customRules)
}

// NewReActStateProvider creates a new state provider.
func NewReActStateProvider() *ReActStateProvider {
	return &ReActStateProvider{
		currentStep:  1,
		observations: make([]string, 0),
	}
}

// AddObservation records a tool observation and advances step.
func (r *ReActStateProvider) AddObservation(obs string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.observations = append(r.observations, obs)
	r.currentStep++
}

// Reset clears state for a new execution.
func (r *ReActStateProvider) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentStep = 1
	r.observations = make([]string, 0)
}

// GetState implements agent.StateProvider interface.
// Returns the current ReAct execution state as a map.
func (r *ReActStateProvider) GetState(ctx context.Context) (map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	state := make(map[string]string)
	state["current_step"] = fmt.Sprintf("%d", r.currentStep)
	state["total_observations"] = fmt.Sprintf("%d", len(r.observations))

	if len(r.observations) > 0 {
		// 最近的观察
		state["last_observation"] = r.observations[len(r.observations)-1]
		// 所有观察的摘要
		state["observations_summary"] = strings.Join(r.observations, "\n")
	}

	return state, nil
}
