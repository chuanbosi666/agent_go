package nvgo

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var DefaultReActInstruction = InstructionsStr(`
你是一个使用ReAct模式解决问题的智能助手。
##思考模式
每次回复必须按照以下格式：

thought:[分析当前情况，思考下一步该做什么]
Action:[选择要执行的工具名称]
Action Input:[工具的输入参数， JSON格式]

当你观察到工具的结果后：

Observation: [工具返回的结果]
Thought: [根据结果思考，是否需要继续行动或已完成任务]

当任务完成时：

Thought: [总结分析过程]
Final Answer: [给用户的最终答案]

## 规则

1. 一步一步思考，不要跳过步骤
2. 每次只执行一个工具
3. 必须根据 Observation 调整后续计划
4. 如果工具执行失败，分析原因并尝试其他方法
5. 完成任务后必须给出 Final Answer
`)

type ReActStateProvider struct {
	currentStep  int
	observations []string
	mu           sync.RWMutex
}

func NewReActInstruction(customRules string) InstructionsStr {
	instruction := DefaultReActInstruction + InstructionsStr(customRules)
	return instruction
}

func NewReActStateProvider() *ReActStateProvider {
	return &ReActStateProvider{
		currentStep:  1,
		observations: make([]string, 0),
	}
}

func (r *ReActStateProvider) GetState(ctx context.Context) (map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return map[string]string{
		"current_step":      fmt.Sprintf("第 %d 步", r.currentStep),
		"observations":      strings.Join(r.observations, ";"),
		"observation_count": fmt.Sprintf("%d", len(r.observations)),
	}, nil
}

func (r *ReActStateProvider) AddObservation(obs string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.observations = append(r.observations, obs)

	r.currentStep++
}

func (r *ReActStateProvider) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentStep = 1
	r.observations = make([]string, 0)
}
