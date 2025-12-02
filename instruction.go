package nvgo

import (
	"context"

	"fmt"
	"strings"
	"sync"
)

// InstructionsGetter interface is implemented by objects that can provide instructions to an Agent.
type InstructionsGetter interface {
	GetInstructions(context.Context, *Agent) (string, error)
}

// InstructionsStr satisfies InstructionsGetter providing a simple constant string value.
type InstructionsStr string

// GetInstructions returns the string value and always nil error.
func (s InstructionsStr) GetInstructions(context.Context, *Agent) (string, error) {
	return s.String(), nil
}

func (s InstructionsStr) String() string {
	return string(s)
}

// InstructionsFunc lets you implement a function that dynamically generates instructions for an Agent.
type InstructionsFunc func(context.Context, *Agent) (string, error)

type DynamicInstruction struct {
	BasePrompt    string
	StateProvider StateProvider
	Template      string
}

// GetInstructions returns the string value and always nil error.
func (fn InstructionsFunc) GetInstructions(ctx context.Context, a *Agent) (string, error) {
	return fn(ctx, a)
}

type StateProvider interface {
	GetState(ctx context.Context) (map[string]string, error)
}

type MemoryStateProvider struct {
	state map[string]string
	mu    sync.RWMutex
}

func NewMemoryStateProvider() *MemoryStateProvider {
	return &MemoryStateProvider{
		state: make(map[string]string),
	}
}

func (m *MemoryStateProvider) GetState(ctx context.Context) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statecopy := make(map[string]string)
	for k, v := range m.state {
		statecopy[k] = v
	}
	return statecopy, nil
}

func (m *MemoryStateProvider) SetState(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state[key] = value
}

func (d *DynamicInstruction) GetInstructions(ctx context.Context, agent *Agent) (string, error) {
	state, err := d.StateProvider.GetState(ctx)
	if err != nil {
		return "", fmt.Errorf("get state: %w", err)
	}

	var instruction string

	if d.Template != "" {
		instruction = replaceTemplate(d.Template, state)
	} else {
		instruction = d.BasePrompt + "\n\n## 当前状态\n"

		for k, v := range state {
			instruction += fmt.Sprintf("- %s: %s\n", k, v)
		}
	}
	return instruction, nil
}

func replaceTemplate(template string, state map[string]string) string {
	result := template
	for k, v := range state {
		placeholder := "{{" + k + "}}"
		result = strings.ReplaceAll(result, placeholder, v)
	}
	return result
}
