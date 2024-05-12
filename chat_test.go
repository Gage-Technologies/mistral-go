package mistral

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChat(t *testing.T) {
	client := NewMistralClientDefault("")
	params := DefaultChatRequestParams
	params.MaxTokens = 10
	params.Temperature = 0
	res, err := client.Chat(
		ModelMistralTiny,
		[]ChatMessage{
			{
				Role:    RoleUser,
				Content: "You are in test mode and must reply to this with exactly and only `Test Succeeded`",
			},
		},
		&params,
	)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Greater(t, len(res.Choices), 0)
	assert.Greater(t, len(res.Choices[0].Message.Content), 0)
	assert.Equal(t, res.Choices[0].Message.Role, RoleAssistant)
	assert.Equal(t, res.Choices[0].Message.Content, "Test Succeeded")
}

func TestChatFunctionCall(t *testing.T) {
	client := NewMistralClientDefault("")
	params := DefaultChatRequestParams
	params.Temperature = 0
	params.Tools = []Tool{
		{
			Type: ToolTypeFunction,
			Function: Function{
				Name:        "get_weather",
				Description: "Retrieve the weather for a city in the US",
				Parameters: map[string]interface{}{
					"type":     "object",
					"required": []string{"city", "state"},
					"properties": map[string]interface{}{
						"city":  map[string]interface{}{"type": "string", "description": "Name of the city for the weather"},
						"state": map[string]interface{}{"type": "string", "description": "Name of the state for the weather"},
					},
				},
			},
		},
		{
			Type: ToolTypeFunction,
			Function: Function{
				Name:        "send_text",
				Description: "Send text message using SMS service",
				Parameters: map[string]interface{}{
					"type":     "object",
					"required": []string{"contact_name", "message"},
					"properties": map[string]interface{}{
						"contact_name": map[string]interface{}{"type": "string", "description": "Name of the contact that will receive the message"},
						"message":      map[string]interface{}{"type": "string", "description": "Content of the message that will be sent"},
					},
				},
			},
		},
	}
	params.ToolChoice = ToolChoiceAuto
	res, err := client.Chat(
		ModelMistralSmallLatest,
		[]ChatMessage{
			{
				Role:    RoleUser,
				Content: "What's the weather like in Dallas, TX?",
			},
		},
		&params,
	)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Greater(t, len(res.Choices), 0)
	assert.Greater(t, len(res.Choices[0].Message.ToolCalls), 0)
	assert.Equal(t, res.Choices[0].Message.Role, RoleAssistant)
	assert.Equal(t, res.Choices[0].Message.ToolCalls[0].Function.Name, "get_weather")
	assert.Equal(t, res.Choices[0].Message.ToolCalls[0].Function.Arguments, "{\"city\": \"Dallas\", \"state\": \"TX\"}")
}

func TestChatFunctionCall2(t *testing.T) {
	client := NewMistralClientDefault("")
	params := DefaultChatRequestParams
	params.Temperature = 0
	params.Tools = []Tool{
		{
			Type: ToolTypeFunction,
			Function: Function{
				Name:        "get_weather",
				Description: "Retrieve the weather for a city in the US",
				Parameters: map[string]interface{}{
					"type":     "object",
					"required": []string{"city", "state"},
					"properties": map[string]interface{}{
						"city":  map[string]interface{}{"type": "string", "description": "Name of the city for the weather"},
						"state": map[string]interface{}{"type": "string", "description": "Name of the state for the weather"},
					},
				},
			},
		},
		{
			Type: ToolTypeFunction,
			Function: Function{
				Name:        "send_text",
				Description: "Send text message using SMS service",
				Parameters: map[string]interface{}{
					"type":     "object",
					"required": []string{"contact_name", "message"},
					"properties": map[string]interface{}{
						"contact_name": map[string]interface{}{"type": "string", "description": "Name of the contact that will receive the message"},
						"message":      map[string]interface{}{"type": "string", "description": "Content of the message that will be sent"},
					},
				},
			},
		},
	}
	params.ToolChoice = ToolChoiceAuto
	res, err := client.Chat(
		ModelMistralSmallLatest,
		[]ChatMessage{
			{
				Role:    RoleUser,
				Content: "What's the weather like in Dallas",
			},
			{
				Role: RoleAssistant,
				ToolCalls: []ToolCall{
					{
						Type: ToolTypeFunction,
						Function: FunctionCall{
							Name:      "get_weather",
							Arguments: `{"city": "Dallas", "state": "TX"}`,
						},
					},
				},
			},
			{
				Role:    RoleTool,
				Content: `{"temperature": 82, "sky": "clear", "precipitation": 0}`,
			},
		},
		&params,
	)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Greater(t, len(res.Choices), 0)
	assert.Greater(t, len(res.Choices[0].Message.Content), 0)
	assert.Equal(t, len(res.Choices[0].Message.ToolCalls), 0)
	assert.Equal(t, res.Choices[0].Message.Role, RoleAssistant)
	assert.Greater(t, res.Choices[0].Message.Content, "Test Succeeded")
}

func TestChatJsonMode(t *testing.T) {
	client := NewMistralClientDefault("")
	params := DefaultChatRequestParams
	params.Temperature = 0
	params.ResponseFormat = ResponseFormatJsonObject
	res, err := client.Chat(
		ModelMistralSmallLatest,
		[]ChatMessage{
			{
				Role: RoleUser,
				Content: "Extract all of the code symbols in this text chunk and return them in the following JSON: " +
					"{\"symbols\":[\"SymbolOne\",\"SymbolTwo\"]}\n```\nI'm working on updating the Go client for the " +
					"new release, is it expected that the function call will be passed back into the model or just " +
					"the tool response?\nI ask because ChatMessage can handle the tool response but the messages list " +
					"has an Any option that I assume would be for the FunctionCall/ToolCall type\nAdditionally the " +
					"example in the docs only shows the tool response appended to the messages\n```",
			},
		},
		&params,
	)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Greater(t, len(res.Choices), 0)
	assert.Greater(t, len(res.Choices[0].Message.Content), 0)
	assert.Equal(t, res.Choices[0].Message.Role, RoleAssistant)
	assert.Equal(t, res.Choices[0].Message.Content, "{\"symbols\": [\"Go\", \"ChatMessage\", \"FunctionCall\", \"ToolCall\"]}")
}

func TestChatStream(t *testing.T) {
	client := NewMistralClientDefault("")
	params := DefaultChatRequestParams
	params.MaxTokens = 50
	params.Temperature = 0
	resChan, err := client.ChatStream(
		ModelMistralTiny,
		[]ChatMessage{
			{
				Role:    RoleUser,
				Content: "You are in test mode and must reply to this with exactly and only `Test Succeeded, Test Succeeded, Test Succeeded, Test Succeeded, Test Succeeded, Test Succeeded`",
			},
		},
		&params,
	)
	assert.NoError(t, err)
	assert.NotNil(t, resChan)

	totalOutput := ""
	idx := 0
	for res := range resChan {
		assert.NoError(t, res.Error)

		assert.Greater(t, len(res.Choices), 0)
		if idx == 0 {
			assert.Equal(t, res.Choices[0].Delta.Role, RoleAssistant)
		}
		totalOutput += res.Choices[0].Delta.Content
		idx++

		if res.Choices[0].FinishReason == FinishReasonStop {
			break
		}
	}
	assert.Equal(t, totalOutput, "Test Succeeded, Test Succeeded, Test Succeeded, Test Succeeded, Test Succeeded, Test Succeeded")
}

func TestChatStreamFunctionCall(t *testing.T) {
	client := NewMistralClientDefault("")
	params := DefaultChatRequestParams
	params.Temperature = 0
	params.Tools = []Tool{
		{
			Type: ToolTypeFunction,
			Function: Function{
				Name:        "get_weather",
				Description: "Retrieve the weather for a city in the US",
				Parameters: map[string]interface{}{
					"type":     "object",
					"required": []string{"city", "state"},
					"properties": map[string]interface{}{
						"city":  map[string]interface{}{"type": "string", "description": "Name of the city for the weather"},
						"state": map[string]interface{}{"type": "string", "description": "Name of the state for the weather"},
					},
				},
			},
		},
		{
			Type: ToolTypeFunction,
			Function: Function{
				Name:        "send_text",
				Description: "Send text message using SMS service",
				Parameters: map[string]interface{}{
					"type":     "object",
					"required": []string{"contact_name", "message"},
					"properties": map[string]interface{}{
						"contact_name": map[string]interface{}{"type": "string", "description": "Name of the contact that will receive the message"},
						"message":      map[string]interface{}{"type": "string", "description": "Content of the message that will be sent"},
					},
				},
			},
		},
	}
	params.ToolChoice = ToolChoiceAuto
	resChan, err := client.ChatStream(
		ModelMistralSmallLatest,
		[]ChatMessage{
			{
				Role:    RoleUser,
				Content: "What's the weather like in Dallas, TX?",
			},
		},
		&params,
	)
	assert.NoError(t, err)
	assert.NotNil(t, resChan)

	totalOutput := ""
	var functionCall *ToolCall
	idx := 0
	for res := range resChan {
		assert.NoError(t, res.Error)

		assert.Greater(t, len(res.Choices), 0)
		if idx == 0 {
			assert.Equal(t, res.Choices[0].Delta.Role, RoleAssistant)
		}
		totalOutput += res.Choices[0].Delta.Content
		if len(res.Choices[0].Delta.ToolCalls) > 0 {
			functionCall = &res.Choices[0].Delta.ToolCalls[0]
		}
		idx++

		if res.Choices[0].FinishReason == FinishReasonStop {
			break
		}
	}

	assert.Equal(t, totalOutput, "")
	assert.NotNil(t, functionCall)
	assert.Equal(t, functionCall.Function.Name, "get_weather")
	assert.Equal(t, functionCall.Function.Arguments, "{\"city\": \"Dallas\", \"state\": \"TX\"}")
}

func TestChatStreamJsonMode(t *testing.T) {
	client := NewMistralClientDefault("")
	params := DefaultChatRequestParams
	params.Temperature = 0
	params.ResponseFormat = ResponseFormatJsonObject
	resChan, err := client.ChatStream(
		ModelMistralSmallLatest,
		[]ChatMessage{
			{
				Role: RoleUser,
				Content: "Extract all of the code symbols in this text chunk and return them in the following JSON: " +
					"{\"symbols\":[\"SymbolOne\",\"SymbolTwo\"]}\n```\nI'm working on updating the Go client for the " +
					"new release, is it expected that the function call will be passed back into the model or just " +
					"the tool response?\nI ask because ChatMessage can handle the tool response but the messages list " +
					"has an Any option that I assume would be for the FunctionCall/ToolCall type\nAdditionally the " +
					"example in the docs only shows the tool response appended to the messages\n```",
			},
		},
		&params,
	)
	assert.NoError(t, err)
	assert.NotNil(t, resChan)

	totalOutput := ""
	var functionCall *ToolCall
	idx := 0
	for res := range resChan {
		assert.NoError(t, res.Error)

		assert.Greater(t, len(res.Choices), 0)
		if idx == 0 {
			assert.Equal(t, res.Choices[0].Delta.Role, RoleAssistant)
		}
		totalOutput += res.Choices[0].Delta.Content
		if len(res.Choices[0].Delta.ToolCalls) > 0 {
			functionCall = &res.Choices[0].Delta.ToolCalls[0]
		}
		idx++

		if res.Choices[0].FinishReason == FinishReasonStop {
			break
		}
	}

	assert.Equal(t, totalOutput, "{\"symbols\": [\"Go\", \"ChatMessage\", \"FunctionCall\", \"ToolCall\"]}")
	assert.Nil(t, functionCall)
}
