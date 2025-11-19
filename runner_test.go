package nvgo

import (
	"testing"

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