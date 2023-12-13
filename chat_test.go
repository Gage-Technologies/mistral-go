package mistral

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChat(t *testing.T) {
	client := NewMistralClient("", "", 3, time.Second*10)
	params := DefaultChatRequestParams
	params.MaxTokens = 10
	params.Temperature = 0
	res, err := client.Chat(
		"mistral-tiny",
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

func TestChatStream(t *testing.T) {
	client := NewMistralClient("", "", 3, time.Second*10)
	params := DefaultChatRequestParams
	params.MaxTokens = 50
	params.Temperature = 0
	resChan, err := client.ChatStream(
		"mistral-tiny",
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

		if res.Choices[0].FinishReason == FinishReasonStop {
			break
		}

		assert.Greater(t, len(res.Choices), 0)
		if idx == 0 {
			assert.Equal(t, res.Choices[0].Delta.Role, RoleAssistant)
		}
		totalOutput += res.Choices[0].Delta.Content
		idx++
	}
	assert.Equal(t, totalOutput, "Test Succeeded, Test Succeeded, Test Succeeded, Test Succeeded, Test Succeeded, Test Succeeded")
}
