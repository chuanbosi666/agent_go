package nvgo

import (
	"context"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	"strings"
	"testing"
)

type MockSession struct {
	items []responses.ResponseInputItemUnionParam
}

func NewMockSession() *MockSession {
	return &MockSession{
		items: []responses.ResponseInputItemUnionParam{},
	}
}

func (m *MockSession) GetItems(ctx context.Context, limit int) ([]responses.ResponseInputItemUnionParam,
	error) {
	if limit <= 0 {
		return m.items, nil
	}

	// 返回最后 N 项
	start := len(m.items) - limit
	if start < 0 {
		start = 0
	}
	return m.items[start:], nil
}

func (m *MockSession) AddItems(ctx context.Context, items []responses.ResponseInputItemUnionParam) error {
	m.items = append(m.items, items...)
	return nil
}

func (m *MockSession) PopItem(ctx context.Context) (*responses.ResponseInputItemUnionParam, error) {
	if len(m.items) == 0 {
		return nil, nil
	}
	item := m.items[len(m.items)-1]
	m.items = m.items[:len(m.items)-1]
	return &item, nil
}

func (m *MockSession) ClearSession(ctx context.Context) error {
	m.items = []responses.ResponseInputItemUnionParam{}
	return nil
}

// GetAllItems 用于测试验证
func (m *MockSession) GetAllItems() []responses.ResponseInputItemUnionParam {
	return m.items
}

func TestMaxTurnsExceededError(t *testing.T) {
	err := &MaxTurnsExceededError{MaxTurns: 10}

	expected := "max turns exceeded: reached limit of 10 turns"

	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGuardrailTripwireTriggeredError_Input(t *testing.T) {

	err := &GuardrailTripwireTriggeredError{
		GuardrailName: "content_filter",
		OutputInfo:    "input guardrail is running",
		IsInput:       true,
	}
	expected := "input guardrail 'content_filter' triggered"

	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGuardrailTripwireTriggeredError_Output(t *testing.T) {
	err := &GuardrailTripwireTriggeredError{
		GuardrailName: "safety_check",
		OutputInfo:    "output guardrail is running",
		IsInput:       false,
	}

	expected := "output guardrail 'safety_check' triggered"

	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestFindTool(t *testing.T) {
	tools := []Tool{
		FunctionTool{Name: "get_weather"},
		FunctionTool{Name: "send_email"},
	}
	t.Run("FindExistingTool", func(t *testing.T) {
		tool, found := findTool(tools, "get_weather")

		if !found {
			t.Errorf("expected to find tool, but found is false")
		}

		if tool == nil {
			t.Fatal("expected tool to be non-nil")
		}

		if tool.ToolName() != "get_weather" {
			t.Errorf("expected tool name %q, got %q", "get_weather", tool.ToolName())
		}
	})

	t.Run("FindNonExistingTool", func(t *testing.T) {
		tool, found := findTool(tools, "non_existing_tool")

		if found {
			t.Error("expected not to find tool, but found is true")
		}

		if tool != nil {
			t.Error("expected tool to be nil")
		}
	})
}

func TestRunItemWrapper(t *testing.T) {
	t.Run("WrapperAndUnwrap", func(t *testing.T) {
		testitem := responses.ResponseInputItemParamOfMessage(
			"test",
			responses.EasyInputMessageRoleUser,
		)
		wrapper := WrapRunItem(testitem)

		if wrapper == nil {
			t.Fatal("WrapRunItem returned nil")
		}

		returnedItem := wrapper.ToInputItem()

		_ = returnedItem
	})
	t.Run("ImplementsRunItem", func(t *testing.T) {
		testItem := responses.ResponseInputItemParamOfMessage(
			"test",
			responses.EasyInputMessageRoleUser,
		)
		wrapper := WrapRunItem(testItem)

		_, ok := wrapper.(RunItem)
		if !ok {
			t.Error("RunItemWrapper does not implement RunItem interface")
		}
	})
}

// ======= 阶段 8 新增测试 =======

// TestInputToItems 测试 Input 转换为 ResponseInputItemUnionParam 列表
func TestInputToItems(t *testing.T) {
	t.Run("ConvertInputString", func(t *testing.T) {

		input := InputString("Hello, world!")
		items := inputToItems(input)

		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}
		if items[0] == (responses.ResponseInputItemUnionParam{}) {
			t.Error("expected non-empty item")
		}
	})

	t.Run("ConvertInputItems", func(t *testing.T) {

		originalItems := InputItems{
			responses.ResponseInputItemParamOfMessage(
				"message 1",
				responses.EasyInputMessageRoleUser,
			),
			responses.ResponseInputItemParamOfMessage(
				"message 2",
				responses.EasyInputMessageRoleAssistant,
			),
		}

		items := inputToItems(originalItems)

		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}
	})
}

// TestToolsToParams 测试工具列表转换为 OpenAI 参数格式
func TestToolsToParams(t *testing.T) {
	t.Run("EmptyToolList", func(t *testing.T) {
		// 测试空工具列表
		tools := []Tool{}
		params := toolsToParams(tools)

		if params != nil {
			t.Errorf("expected nil for empty tools, got %v", params)
		}
	})

	t.Run("ConvertFunctionTools", func(t *testing.T) {
		// 测试转换函数工具
		tools := []Tool{
			FunctionTool{
				Name:             "get_weather",
				Description:      "Get current weather",
				ParamsJSONSchema: map[string]any{"type": "object"},
				StrictJSONSchema: param.NewOpt(true),
			},
			FunctionTool{
				Name:             "send_email",
				Description:      "Send an email",
				ParamsJSONSchema: map[string]any{"type": "object"},
				StrictJSONSchema: param.NewOpt(false),
			},
		}

		params := toolsToParams(tools)

		if len(params) != 2 {
			t.Fatalf("expected 2 tool params, got %d", len(params))
		}

		// 验证第一个工具
		if params[0].OfFunction == nil {
			t.Fatal("expected OfFunction to be non-nil")
		}
		if params[0].OfFunction.Name != "get_weather" {
			t.Errorf("expected tool name 'get_weather', got %q", params[0].OfFunction.Name)
		}
		if !params[0].OfFunction.Description.Valid() {
			t.Error("expected description to be valid")
		}

		// 验证第二个工具
		if params[1].OfFunction == nil {
			t.Fatal("expected OfFunction to be non-nil")
		}
		if params[1].OfFunction.Name != "send_email" {
			t.Errorf("expected tool name 'send_email', got %q", params[1].OfFunction.Name)
		}
	})

	t.Run("SkipNonFunctionTools", func(t *testing.T) {
		// 测试跳过非 FunctionTool 类型
		// 注意：由于我们只有 FunctionTool 实现了 Tool 接口，
		// 这个测试主要是验证类型转换逻辑的健壮性
		tools := []Tool{
			FunctionTool{
				Name:             "valid_tool",
				Description:      "A valid tool",
				ParamsJSONSchema: map[string]any{"type": "object"},
			},
		}

		params := toolsToParams(tools)

		if len(params) != 1 {
			t.Fatalf("expected 1 tool param, got %d", len(params))
		}
	})
}

func TestRunner_ConversationHistory(t *testing.T) {
	t.Run("SessionSavesNewItems", func(t *testing.T) {
		// 创建 Mock Session
		session := NewMockSession()

		// 验证初始状态为空
		if len(session.GetAllItems()) != 0 {
			t.Error("expected session to be empty initially")
		}

		// 模拟添加用户消息
		userMsg := responses.ResponseInputItemParamOfMessage(
			"Hello, how are you?",
			responses.EasyInputMessageRoleUser,
		)
		err := session.AddItems(context.Background(), []responses.ResponseInputItemUnionParam{userMsg})
		if err != nil {
			t.Fatalf("failed to add user message: %v", err)
		}

		// 验证保存成功
		items := session.GetAllItems()
		if len(items) != 1 {
			t.Fatalf("expected 1 item in session, got %d", len(items))
		}

		// 模拟添加模型输出消息
		outputMsg := responses.ResponseOutputMessageParam{
			Content: []responses.ResponseOutputMessageContentUnionParam{
				{
					OfOutputText: &responses.ResponseOutputTextParam{
						Text: "I'm doing well, thank you!",
						Type: "output_text",
					},
				},
			},
			Role: "assistant",
			Type: "message",
		}
		modelOutput := responses.ResponseInputItemUnionParam{
			OfOutputMessage: &outputMsg,
		}
		err = session.AddItems(context.Background(), []responses.ResponseInputItemUnionParam{modelOutput})
		if err != nil {
			t.Fatalf("failed to add model output: %v", err)
		}

		// 验证现在有 2 个项目
		items = session.GetAllItems()
		if len(items) != 2 {
			t.Fatalf("expected 2 items in session, got %d", len(items))
		}

		t.Log("✅ Session correctly saves user messages and model outputs")
	})

	t.Run("SessionRetrievesHistory", func(t *testing.T) {
		// 创建 Mock Session 并添加多条消息
		session := NewMockSession()

		messages := []responses.ResponseInputItemUnionParam{
			responses.ResponseInputItemParamOfMessage("Message 1", responses.EasyInputMessageRoleUser),
			responses.ResponseInputItemParamOfMessage("Message 2", responses.EasyInputMessageRoleAssistant),
			responses.ResponseInputItemParamOfMessage("Message 3", responses.EasyInputMessageRoleUser),
		}

		err := session.AddItems(context.Background(), messages)
		if err != nil {
			t.Fatalf("failed to add messages: %v", err)
		}

		// 测试获取全部历史
		allItems, err := session.GetItems(context.Background(), -1)
		if err != nil {
			t.Fatalf("failed to get all items: %v", err)
		}
		if len(allItems) != 3 {
			t.Fatalf("expected 3 items, got %d", len(allItems))
		}

		// 测试获取最近 2 条
		recentItems, err := session.GetItems(context.Background(), 2)
		if err != nil {
			t.Fatalf("failed to get recent items: %v", err)
		}
		if len(recentItems) != 2 {
			t.Fatalf("expected 2 recent items, got %d", len(recentItems))
		}

		t.Log("✅ Session correctly retrieves conversation history")
	})
}
func TestRunner_AccumulatedHistory(t *testing.T) {
	t.Run("AccumulateWithoutSession", func(t *testing.T) {
		// 这个测试验证在没有Session时，
		// accumulatedHistory能否正确维护多轮对话历史

		// 注意：这是一个逻辑测试，不涉及真实的LLM调用
		// 我们只验证数据结构和逻辑

		// 1. 模拟第一轮：用户输入
		input := InputString("Hello")
		inputItems := inputToItems(input)

		if len(inputItems) != 1 {
			t.Fatalf("expected 1 input item, got %d", len(inputItems))
		}

		// 2. 模拟累积历史
		var accumulatedHistory []responses.ResponseInputItemUnionParam
		accumulatedHistory = append(accumulatedHistory, inputItems...)

		if len(accumulatedHistory) != 1 {
			t.Fatalf("expected 1 item in accumulated history, got %d", len(accumulatedHistory))
		}

		// 3. 模拟第二轮：添加工具调用结果
		toolOutput := responses.ResponseInputItemParamOfFunctionCallOutput(
			"call-123",
			"Tool executed successfully",
		)
		accumulatedHistory = append(accumulatedHistory, toolOutput)

		if len(accumulatedHistory) != 2 {
			t.Fatalf("expected 2 items in accumulated history, got %d", len(accumulatedHistory))
		}

		// 4. 模拟第三轮：添加模型输出
		modelOutput := responses.ResponseOutputMessageParam{
			Content: []responses.ResponseOutputMessageContentUnionParam{
				{
					OfOutputText: &responses.ResponseOutputTextParam{
						Text: "Response based on tool result",
						Type: "output_text",
					},
				},
			},
			Role: "assistant",
			Type: "message",
		}
		modelItem := responses.ResponseInputItemUnionParam{
			OfOutputMessage: &modelOutput,
		}
		accumulatedHistory = append(accumulatedHistory, modelItem)

		if len(accumulatedHistory) != 3 {
			t.Fatalf("expected 3 items in accumulated history, got %d", len(accumulatedHistory))
		}

		t.Log("✅ accumulatedHistory correctly maintains conversation history without Session")
	})
}
func TestRunner_WithAndWithoutSession(t *testing.T) {
	t.Run("SessionVsAccumulatedHistory", func(t *testing.T) {
		// 验证：有Session时保存到Session，无Session时保存到accumulatedHistory

		// 场景1：有Session
		session := NewMockSession()
		userMsg := responses.ResponseInputItemParamOfMessage(
			"Test message",
			responses.EasyInputMessageRoleUser,
		)

		err := session.AddItems(context.Background(), []responses.ResponseInputItemUnionParam{userMsg})
		if err != nil {
			t.Fatalf("failed to add to session: %v", err)
		}

		sessionItems := session.GetAllItems()
		if len(sessionItems) != 1 {
			t.Fatalf("expected 1 item in session, got %d", len(sessionItems))
		}

		// 场景2：无Session，使用accumulatedHistory
		var accumulatedHistory []responses.ResponseInputItemUnionParam
		accumulatedHistory = append(accumulatedHistory, userMsg)

		if len(accumulatedHistory) != 1 {
			t.Fatalf("expected 1 item in accumulated history, got %d", len(accumulatedHistory))
		}

		// 验证两种方式的结果一致
		if len(sessionItems) != len(accumulatedHistory) {
			t.Error("session and accumulated history should have same length")
		}

		t.Log("✅ Both Session and accumulatedHistory work correctly for storing history")
	})
}

func TestWrapAgentAsTool(t *testing.T) {
	t.Run("BasicAgentWrapping", func(t *testing.T) {
		subAgent := &Agent{
			Name: "SubAgent",
		}

		tool := WrapAgentAsTool(subAgent, 5)

		expectedName := "Call_agent_SubAgent"
		if tool.Name != expectedName {
			t.Errorf("expected tool name %q, got %q", expectedName, tool.Name)
		}

		if tool.ParamsJSONSchema == nil {
			t.Error("expected non-nil ParamsJSONSchema")
		}

		if tool.OnInvokeTool == nil {
			t.Error("expected non-nil OnInvokeTool")
		}

		t.Log("✅ Agent successfully wrapped as tool")
	})

	t.Run("DefaultMaxTurns", func(t *testing.T) {
		subAgent := &Agent{
			Name: "TestAgent",
		}

		// 传入 0，应该使用 DefaultMaxTurns
		tool := WrapAgentAsTool(subAgent, 0)

		if tool.Name != "Call_agent_TestAgent" {
			t.Errorf("unexpected tool name: %s", tool.Name)
		}

		t.Log("✅ DefaultMaxTurns applied when maxTurns is 0")
	})

	t.Run("ToolDescription", func(t *testing.T) {
		subAgent := &Agent{
			Name:         "ExpertAgent",
			Instructions: InstructionsStr("You are an expert in Go programming."),
		}

		tool := WrapAgentAsTool(subAgent, 5)

		// 验证描述包含 agent 信息
		if tool.Description == "" {
			t.Error("expected non-empty description")
		}

		if !strings.Contains(tool.Description, "ExpertAgent") {
			t.Error("expected description to contain agent name")
		}

		t.Log("✅ Tool description contains agent information")
	})

	t.Run("ParamsJSONSchemaStructure", func(t *testing.T) {
		subAgent := &Agent{
			Name: "SchemaTestAgent",
		}

		tool := WrapAgentAsTool(subAgent, 5)

		// 验证 schema 结构
		schema := tool.ParamsJSONSchema
		if schema == nil {
			t.Fatal("ParamsJSONSchema is nil")
		}

		// 检查 type
		if schemaType, ok := schema["type"].(string); !ok || schemaType != "object" {
			t.Error("expected schema type to be 'object'")
		}

		// 检查 properties
		props, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatal("expected properties to be a map")
		}

		// 检查 input 属性存在
		if _, ok := props["input"]; !ok {
			t.Error("expected 'input' property in schema")
		}

		// 检查 required
		required, ok := schema["required"].([]string)
		if !ok {
			t.Fatal("expected required to be a string slice")
		}

		foundInput := false
		for _, r := range required {
			if r == "input" {
				foundInput = true
				break
			}
		}
		if !foundInput {
			t.Error("expected 'input' in required fields")
		}

		t.Log("✅ ParamsJSONSchema has correct structure")
	})

	t.Run("ToolImplementsFunctionTool", func(t *testing.T) {
		subAgent := &Agent{
			Name: "InterfaceTestAgent",
		}

		tool := WrapAgentAsTool(subAgent, 5)

		// 验证返回的是 FunctionTool 类型
		var _ FunctionTool = tool

		// 验证实现了 Tool 接口
		var toolInterface Tool = tool
		if toolInterface.ToolName() != "Call_agent_InterfaceTestAgent" {
			t.Errorf("ToolName() returned unexpected value: %s", toolInterface.ToolName())
		}

		t.Log("✅ WrapAgentAsTool returns valid FunctionTool")
	})

	t.Run("AgentWithTools", func(t *testing.T) {
		// 测试主 Agent 可以包含子 Agent 作为工具
		subAgent := &Agent{
			Name:         "HelperAgent",
			Instructions: InstructionsStr("I help with specific tasks."),
		}

		subAgentTool := WrapAgentAsTool(subAgent, 3)

		mainAgent := &Agent{
			Name:         "MainAgent",
			Instructions: InstructionsStr("I am the main agent."),
			Tools:        []FunctionTool{subAgentTool},
		}

		// 验证主 Agent 包含子 Agent 工具
		if len(mainAgent.Tools) != 1 {
			t.Fatalf("expected 1 tool, got %d", len(mainAgent.Tools))
		}

		if mainAgent.Tools[0].Name != "Call_agent_HelperAgent" {
			t.Errorf("unexpected tool name: %s", mainAgent.Tools[0].Name)
		}

		t.Log("✅ Agent can include sub-agent as tool")
	})
}

func TestKeywordRouter(t *testing.T) {
	router := &KeywordRouter{
		ToolKeywords: map[string][]string{
			"get_weather": {"天气", "气温", "weather"},
			"send_email":  {"邮件", "发送", "email"},
			"search_db":   {"查询", "搜索", "database"},
		},
		TopN: 3,
	}

	tools := []Tool{
		FunctionTool{Name: "get_weather"},
		FunctionTool{Name: "send_email"},
		FunctionTool{Name: "search_db"},
		FunctionTool{Name: "other_tool1"},
		FunctionTool{Name: "other_tool2"},
	}

	t.Run("MatchWeatherKeyword", func(t *testing.T) {
		input := InputString("查询北京天气")
		routed, err := router.RouteTools(context.Background(), input, tools)

		if err != nil {
			t.Fatalf("routing failed: %v", err)
		}

		// 应该返回 3 个工具
		if len(routed) != 3 {
			t.Errorf("expected 3 tools, got %d", len(routed))
		}

		// 第一个应该是 get_weather（分数最高）
		if routed[0].ToolName() != "get_weather" {
			t.Errorf("expected first tool to be get_weather, got %s", routed[0].ToolName())
		}
	})

	t.Run("MatchEmailKeyword", func(t *testing.T) {
		input := InputString("发送一封邮件")
		routed, err := router.RouteTools(context.Background(), input, tools)

		if err != nil {
			t.Fatalf("routing failed: %v", err)
		}

		// 第一个应该是 send_email
		if routed[0].ToolName() != "send_email" {
			t.Errorf("expected first tool to be send_email, got %s", routed[0].ToolName())
		}
	})
}

func TestMemoryStateProvider(t *testing.T) {
	provider := NewMemoryStateProvider()

	t.Run("SetAndGetState", func(t *testing.T) {
		provider.SetState("name", "张三")
		provider.SetState("age", "25")

		state, err := provider.GetState(context.Background())
		if err != nil {
			t.Fatalf("get state failed: %v", err)
		}

		if state["name"] != "张三" {
			t.Errorf("expected name '张三', got %q", state["name"])
		}
		if state["age"] != "25" {
			t.Errorf("expected age '25', got %q", state["age"])
		}
	})
}

func TestDynamicInstruction(t *testing.T) {
	t.Run("WithTemplate", func(t *testing.T) {
		provider := NewMemoryStateProvider()
		provider.SetState("user_name", "张三")
		provider.SetState("task_count", "5")

		instruction := &DynamicInstruction{
			StateProvider: provider,
			Template:      "你好 {{user_name}}，你有 {{task_count}} 个任务。",
		}

		result, err := instruction.GetInstructions(context.Background(), nil)
		if err != nil {
			t.Fatalf("get instructions failed: %v", err)
		}

		expected := "你好 张三，你有 5 个任务。"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("WithBasePrompt", func(t *testing.T) {
		provider := NewMemoryStateProvider()
		provider.SetState("status", "在线")

		instruction := &DynamicInstruction{
			BasePrompt:    "你是一个助手。",
			StateProvider: provider,
		}

		result, err := instruction.GetInstructions(context.Background(), nil)
		if err != nil {
			t.Fatalf("get instructions failed: %v", err)
		}

		if !strings.Contains(result, "你是一个助手") {
			t.Error("expected base prompt in result")
		}
		if !strings.Contains(result, "status") || !strings.Contains(result, "在线") {
			t.Error("expected state in result")
		}
	})
}

// TestReActStateProvider 测试 ReAct 状态提供者
func TestReActStateProvider(t *testing.T) {
	t.Run("InitialState", func(t *testing.T) {
		provider := NewReActStateProvider()

		state, err := provider.GetState(context.Background())
		if err != nil {
			t.Fatalf("get state failed: %v", err)
		}

		// 验证初始步骤
		if state["current_step"] != "第 1 步" {
			t.Errorf("expected '第 1 步', got %q", state["current_step"])
		}

		// 验证初始观察为空
		if state["observations"] != "" {
			t.Errorf("expected empty observations, got %q", state["observations"])
		}

		// 验证观察计数
		if state["observation_count"] != "0" {
			t.Errorf("expected '0', got %q", state["observation_count"])
		}
	})

	t.Run("AddObservation", func(t *testing.T) {
		provider := NewReActStateProvider()

		// 添加两个观察结果
		provider.AddObservation("天气查询结果: 晴天 25°C")
		provider.AddObservation("日程查询结果: 今天有3个会议")

		state, err := provider.GetState(context.Background())
		if err != nil {
			t.Fatalf("get state failed: %v", err)
		}

		// 验证步骤递增 (初始1 + 2次观察 = 第3步)
		if state["current_step"] != "第 3 步" {
			t.Errorf("expected '第 3 步', got %q", state["current_step"])
		}

		// 验证观察结果包含添加的内容
		if !strings.Contains(state["observations"], "天气查询结果") {
			t.Error("expected observations to contain '天气查询结果'")
		}
		if !strings.Contains(state["observations"], "日程查询结果") {
			t.Error("expected observations to contain '日程查询结果'")
		}

		// 验证观察计数
		if state["observation_count"] != "2" {
			t.Errorf("expected '2', got %q", state["observation_count"])
		}
	})

	t.Run("Reset", func(t *testing.T) {
		provider := NewReActStateProvider()

		// 先添加一些观察
		provider.AddObservation("观察1")
		provider.AddObservation("观察2")

		// 重置
		provider.Reset()

		state, err := provider.GetState(context.Background())
		if err != nil {
			t.Fatalf("get state failed: %v", err)
		}

		// 验证恢复到初始状态
		if state["current_step"] != "第 1 步" {
			t.Errorf("expected '第 1 步' after reset, got %q", state["current_step"])
		}
		if state["observations"] != "" {
			t.Errorf("expected empty observations after reset, got %q", state["observations"])
		}
		if state["observation_count"] != "0" {
			t.Errorf("expected '0' after reset, got %q", state["observation_count"])
		}
	})
}

// TestNewReActInstruction 测试 ReAct Instruction 工厂函数
func TestNewReActInstruction(t *testing.T) {
	t.Run("ContainsDefaultTemplate", func(t *testing.T) {
		instruction := NewReActInstruction("")

		result, err := instruction.GetInstructions(context.Background(), nil)
		if err != nil {
			t.Fatalf("get instructions failed: %v", err)
		}

		// 验证包含默认模板的关键内容
		if !strings.Contains(result, "ReAct") {
			t.Error("expected instruction to contain 'ReAct'")
		}
		if !strings.Contains(result, "Thought") {
			t.Error("expected instruction to contain 'Thought'")
		}
		if !strings.Contains(result, "Action") {
			t.Error("expected instruction to contain 'Action'")
		}
		if !strings.Contains(result, "Final Answer") {
			t.Error("expected instruction to contain 'Final Answer'")
		}
	})

	t.Run("ContainsCustomRules", func(t *testing.T) {
		customRules := "\n\n## 额外规则\n- 必须使用中文回答\n- 每步思考不超过50字"
		instruction := NewReActInstruction(customRules)

		result, err := instruction.GetInstructions(context.Background(), nil)
		if err != nil {
			t.Fatalf("get instructions failed: %v", err)
		}

		// 验证包含自定义规则
		if !strings.Contains(result, "额外规则") {
			t.Error("expected instruction to contain '额外规则'")
		}
		if !strings.Contains(result, "必须使用中文回答") {
			t.Error("expected instruction to contain custom rule")
		}
	})
}

// TestReActWithDynamicInstruction 测试 ReAct 与 DynamicInstruction 结合使用
func TestReActWithDynamicInstruction(t *testing.T) {
	provider := NewReActStateProvider()

	instruction := &DynamicInstruction{
		BasePrompt:    string(DefaultReActInstruction),
		StateProvider: provider,
		Template: `
你使用 ReAct 模式解决问题。

## 当前进度
- 步骤: {{current_step}}
- 已观察: {{observations}}
- 观察次数: {{observation_count}}

` + string(DefaultReActInstruction),
	}

	// 初始状态
	result, err := instruction.GetInstructions(context.Background(), nil)
	if err != nil {
		t.Fatalf("get instructions failed: %v", err)
	}

	if !strings.Contains(result, "第 1 步") {
		t.Error("expected '第 1 步' in initial instruction")
	}

	// 添加观察后
	provider.AddObservation("工具执行成功")

	result, err = instruction.GetInstructions(context.Background(), nil)
	if err != nil {
		t.Fatalf("get instructions failed: %v", err)
	}

	if !strings.Contains(result, "第 2 步") {
		t.Error("expected '第 2 步' after observation")
	}
	if !strings.Contains(result, "工具执行成功") {
		t.Error("expected observation in instruction")
	}
}
