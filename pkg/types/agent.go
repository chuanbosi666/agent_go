package types

// AgentLike 是跨包使用的最小 Agent 接口
// 只包含核心标识方法，避免包之间的循环依赖
type AgentLike interface {
	// GetName 返回 Agent 的名称
	GetName() string
	// GetModel 返回使用的 LLM 模型名称
	GetModel() string
}
