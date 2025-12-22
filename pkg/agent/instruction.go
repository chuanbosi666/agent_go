package agent

import(
	"context"
	"fmt"
	"strings"
	"sync"
)

// InstructionsGetter provides instructions (system prompt) for an Agent.
type Instructions interface {
	GetInstructions(context.Context, *Agent) (string, error)
}

// InstructionsStr is a simple static string instruction.
type InstructionsStr string

// GetInstructions returns the string value.
func(s InstructionsStr) GetInstructions(context.Context, *Agent)(string, error){
	return s.String(), nil
}

// String returns the underlying string.
func (s InstructionsStr) String() string{
	return string(s)
}

//InstructionsFunc is a function that dynamically generates instructions.
type InstructionsFunc func(context.Context, *Agent)(string, error)

func(fn InstructionsFunc) GetInstructions(ctx context.Context, a *Agent) (string, error){
	return fn(ctx,a)
}

//
type StateProvider interface{
	GetState(ctx context.Context)(map[string]string, error)
}

// MemoryStateProvider is an in-memory thread-safe state provider.
type MemoryStateProvider struct{
	state 	map[string]string
	mu  	sync.RWMutex
}

// SetState sets a state value.
func (m *MemoryStateProvider) SetState(key, value string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.state[key] = value
}

// NewMemoryStateProvider creates a new in-memory state provider.
func NewMemoryStateProvider() *MemoryStateProvider {
    return &MemoryStateProvider{
            state: make(map[string]string),
    }
}

// GetState returns a copy of the current state.
func (m *MemoryStateProvider) GetState(ctx context.Context) (map[string]string, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    statecopy := make(map[string]string)
    for k, v := range m.state {
            statecopy[k] = v
    }
    return statecopy, nil
}

// DynamicInstruction generates instructions from a template and state provider.
type DynamicInstruction struct {
    // BasePrompt is the base instruction text.
    BasePrompt string
    // StateProvider provides dynamic state values.
	StateProvider StateProvider	
	// Template uses {{key}} placeholders replaced by state values.
    // If empty, state is appended as a list to BasePrompt.
    Template string
}
// GetInstructions generates the instruction by merging template with state.
func (d *DynamicInstruction) GetInstructions(ctx context.Context, agent *Agent) (string, error) {
        state, err := d.StateProvider.GetState(ctx)
        if err != nil {
                return "", fmt.Errorf("get state: %w", err)
        }

        var instruction string

        if d.Template != "" {
                instruction = replaceTemplate(d.Template, state)
        } else {
            instruction = d.BasePrompt + "\n\n## Current State\n"
            for k, v := range state {
                    instruction += fmt.Sprintf("- %s: %s\n", k, v)
            }
    }
    return instruction, nil
}

// replaceTemplate replaces {{key}} placeholders with state values.
func replaceTemplate(template string, state map[string]string) string {
    result := template
    for k, v := range state {
            placeholder := "{{" + k + "}}"
            result = strings.ReplaceAll(result, placeholder, v)
	}
    return result
}
  // InstructionsTemplate is a template-based instruction.
  type InstructionsTemplate string

  func (t InstructionsTemplate) GetInstructions(ctx context.Context, state map[string]string) (string, error) {  
        result := string(t)
        for k, v := range state {
                placeholder := "{{" + k + "}}"
                result = strings.ReplaceAll(result, placeholder, v)
        }
        return result, nil
  }

  // StateProviderFunc is a function adapter for StateProvider.
  type StateProviderFunc func(ctx context.Context) (map[string]string, error)

  func (f StateProviderFunc) GetState(ctx context.Context) (map[string]string, error) {
        return f(ctx)
  }