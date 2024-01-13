package mistral

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

type FinishReason string

const (
	FinishReasonStop   FinishReason = "stop"
	FinishReasonLength FinishReason = "length"
)

// ChatRequestParams represents the parameters for the Chat/ChatStream method of MistralClient.
type ChatRequestParams struct {
	Temperature float64 `json:"temperature"` // The temperature to use for sampling. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic. We generally recommend altering this or TopP but not both.
	TopP        float64 `json:"top_p"`       // An alternative to sampling with temperature, called nucleus sampling, where the model considers the results of the tokens with top_p probability mass. So 0.1 means only the tokens comprising the top 10% probability mass are considered. We generally recommend altering this or Temperature but not both.
	RandomSeed  int     `json:"random_seed"`
	MaxTokens   int     `json:"max_tokens"`
}

var DefaultChatRequestParams = ChatRequestParams{
	Temperature: 1,
	TopP:        1,
	RandomSeed:  42069,
	MaxTokens:   4000,
}

// ChatMessage represents a single message in a chat.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponseChoice represents a choice in the chat completion response.
type ChatCompletionResponseChoice struct {
	Index        int          `json:"index"`
	Message      ChatMessage  `json:"message"`
	FinishReason FinishReason `json:"finish_reason,omitempty"`
}

// ChatCompletionResponseChoice represents a choice in the chat completion response.
type ChatCompletionResponseChoiceStream struct {
	Index        int          `json:"index"`
	Delta        ChatMessage  `json:"delta"`
	FinishReason FinishReason `json:"finish_reason,omitempty"`
}

// ChatCompletionResponse represents the response from the chat completion endpoint.
type ChatCompletionResponse struct {
	ID      string                         `json:"id"`
	Object  string                         `json:"object"`
	Created int                            `json:"created"`
	Model   string                         `json:"model"`
	Choices []ChatCompletionResponseChoice `json:"choices"`
	Usage   UsageInfo                      `json:"usage"`
}

// ChatCompletionStreamResponse represents the streamed response from the chat completion endpoint.
type ChatCompletionStreamResponse struct {
	ID      string                               `json:"id"`
	Model   string                               `json:"model"`
	Choices []ChatCompletionResponseChoiceStream `json:"choices"`
	Created int                                  `json:"created,omitempty"`
	Object  string                               `json:"object,omitempty"`
	Usage   UsageInfo                            `json:"usage,omitempty"`
	Error   error                                `json:"error,omitempty"`
}

// UsageInfo represents the usage information of a response.
type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
}

func (c *MistralClient) Chat(model string, messages []ChatMessage, params *ChatRequestParams) (*ChatCompletionResponse, error) {
	if params == nil {
		params = &DefaultChatRequestParams
	}

	requestData := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"temperature": params.Temperature,
		"max_tokens":  params.MaxTokens,
		"top_p":       params.TopP,
		"random_seed": params.RandomSeed,
	}

	response, err := c.request(http.MethodPost, requestData, "v1/chat/completions", false, nil)
	if err != nil {
		return nil, err
	}

	respData, ok := response.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response type: %T", response)
	}

	var chatResponse ChatCompletionResponse
	err = mapToStruct(respData, &chatResponse)
	if err != nil {
		return nil, err
	}

	return &chatResponse, nil
}

// ChatStream sends a chat message and returns a channel to receive streaming responses.
func (c *MistralClient) ChatStream(model string, messages []ChatMessage, params *ChatRequestParams) (<-chan ChatCompletionStreamResponse, error) {
	if params == nil {
		params = &DefaultChatRequestParams
	}

	responseChannel := make(chan ChatCompletionStreamResponse)

	requestData := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"temperature": params.Temperature,
		"max_tokens":  params.MaxTokens,
		"top_p":       params.TopP,
		"random_seed": params.RandomSeed,
		"stream":      true,
	}

	response, err := c.request(http.MethodPost, requestData, "v1/chat/completions", true, nil)
	if err != nil {
		return nil, err
	}

	respBody, ok := response.(io.ReadCloser)
	if !ok {
		return nil, fmt.Errorf("invalid response type: %T", response)
	}

	// Execute the HTTP request in a separate goroutine.
	go func() {
		defer close(responseChannel)
		defer respBody.Close()

		// Assuming ChatCompletionStreamResponse is already defined in your Go code.
		// Assuming responseChannel is a channel of ChatCompletionStreamResponse.

		// Create a buffered reader to read the stream line by line.
		reader := bufio.NewReader(respBody)

		for {
			// Read a line from the buffered reader.
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break // End of stream.
			} else if err != nil {
				responseChannel <- ChatCompletionStreamResponse{Error: fmt.Errorf("error reading stream response: %w", err)}
				return
			}

			// Skip empty lines.
			if bytes.Equal(line, []byte("\n")) {
				continue
			}

			// Check if the line starts with "data: ".
			if bytes.HasPrefix(line, []byte("data: ")) {
				// Trim the prefix and any leading or trailing whitespace.
				jsonLine := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data: ")))

				// Check for the special "[DONE]" message.
				if bytes.Equal(jsonLine, []byte("[DONE]")) {
					break
				}

				// Decode the JSON object from the line.
				var streamResponse ChatCompletionStreamResponse
				if err := json.Unmarshal(jsonLine, &streamResponse); err != nil {
					responseChannel <- ChatCompletionStreamResponse{Error: fmt.Errorf("error decoding stream response: %w", err)}
					continue
				}

				// Send the decoded response to the channel.
				responseChannel <- streamResponse
			}
		}
	}()

	// Return the response channel.
	return responseChannel, nil
}

// mapToStruct is a helper function to convert a map to a struct.
func mapToStruct(m map[string]interface{}, s interface{}) error {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, s)
}
