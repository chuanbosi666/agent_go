package nvgo

import (
	"testing"

	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

func TestMaxTurnsExceededError(t *testing.T){
	err := &MaxTurnsExceededError{MaxTurns:10}

	expected := "max turns exceeded: reached limit of 10 turns"

	if err.Error() != expected{
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGuardrailTripwireTriggeredError_Input(t *testing.T){

	err := &GuardrailTripwireTriggeredError{
		GuardrailName: "content_filter",
		OutputInfo: "input guardrail is running",
		IsInput: true,
	}
	expected := "input guardrail 'content_filter' triggered"

	if err.Error() != expected{
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGuardrailTripwireTriggeredError_Output(t *testing.T){
	err := &GuardrailTripwireTriggeredError{
		GuardrailName: "safety_check",
		OutputInfo: "output guardrail is running",
		IsInput: false,
	}

	expected := "output guardrail 'safety_check' triggered"

	if err.Error() != expected{
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestFindTool(t *testing.T){
	tools := []Tool{
		FunctionTool{Name: "get_weather"},
		FunctionTool{Name: "send_email"},
	}
	t.Run("FindExistingTool", func(t *testing.T) {
		tool, found := findTool(tools, "get_weather")

		if !found {
			t.Errorf("expected to find tool, but found is false")
		}

		if tool == nil{
			t.Fatal("expected tool to be non-nil")
		}

		if tool.ToolName() != "get_weather"{
			t.Errorf("expected tool name %q, got %q", "get_weather", tool.ToolName())
		}
	})

	t.Run("FindNonExistingTool", func(t *testing.T){
		tool, found := findTool(tools, "non_existing_tool")

		if found {
			t.Error("expected not to find tool, but found is true")
		}

		if tool != nil{
			t.Error("expected tool to be nil")
		}
	})
}

func TestRunItemWrapper(t *testing.T){
	t.Run("WrapperAndUnwrap", func(t *testing.T){
		testitem := responses.ResponseInputItemParamOfMessage(
			"test",
			responses.EasyInputMessageRoleUser,
		)
		wrapper := WrapRunItem(testitem)

		if wrapper == nil{
			t.Fatal("WrapRunItem returned nil")
		}

		returnedItem := wrapper.ToInputItem()

		_=returnedItem
	})
	t.Run("ImplementsRunItem", func (t *testing.T)  {
		testItem := responses.ResponseInputItemParamOfMessage(
			"test",
			responses.EasyInputMessageRoleUser,
		)
		wrapper := WrapRunItem(testItem)

		_, ok :=wrapper.(RunItem)
		if !ok{
			t.Error("RunItemWrapper does not implement RunItem interface")
		}
	})
}

// ======= 阶段 8 新增测试 =======

// TestInputToItems 测试 Input 转换为 ResponseInputItemUnionParam 列表
func TestInputToItems(t *testing.T) {
	t.Run("ConvertInputString", func(t *testing.T) {
		// 测试字符串输入
		input := InputString("Hello, world!")
		items := inputToItems(input)

		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}

		// 验证转换后的项不为空
		if items[0] == (responses.ResponseInputItemUnionParam{}) {
			t.Error("expected non-empty item")
		}
	})

	t.Run("ConvertInputItems", func(t *testing.T) {
		// 测试列表输入
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
