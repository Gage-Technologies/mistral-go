package mistral

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFIM(t *testing.T) {
	client := NewMistralClientDefault("")
	params := FIMRequestParams{
		Model:       ModelCodestralLatest,
		Prompt:      "def f(",
		Suffix:      "return a + b",
		MaxTokens:   64,
		Temperature: 0,
		Stop:        []string{"\n"},
	}
	res, err := client.FIM(&params)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Greater(t, len(res.Choices), 0)
	assert.Equal(t, res.Choices[0].Message.Content, "a, b):")
	assert.Equal(t, res.Choices[0].FinishReason, FinishReasonStop)
}

func TestFIMWithStop(t *testing.T) {
	client := NewMistralClientDefault("")
	params := FIMRequestParams{
		Model:       ModelCodestralLatest,
		Prompt:      "def is_odd(n): \n return n % 2 == 1 \n def test_is_odd():",
		Suffix:      "test_is_odd()",
		MaxTokens:   64,
		Temperature: 0,
		Stop:        []string{"False"},
	}
	res, err := client.FIM(&params)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Greater(t, len(res.Choices), 0)
	assert.Equal(t, res.Choices[0].Message.Content, "\n assert is_odd(1) == True\n assert is_odd(2) == ")
	assert.Equal(t, res.Choices[0].FinishReason, FinishReasonStop)
}

func TestFIMInvalidModel(t *testing.T) {
	client := NewMistralClientDefault("")
	params := FIMRequestParams{
		Model:       "invalid-model",
		Prompt:      "This is a test prompt",
		Suffix:      "This is a test suffix",
		MaxTokens:   10,
		Temperature: 0.5,
	}
	res, err := client.FIM(&params)
	assert.Error(t, err)
	assert.Nil(t, res)
}
