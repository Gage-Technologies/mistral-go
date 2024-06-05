package mistral

import (
	"fmt"
	"net/http"
)

// FIMRequestParams represents the parameters for the FIM method of MistralClient.
type FIMRequestParams struct {
	Model       string   `json:"model"`
	Prompt      string   `json:"prompt"`
	Suffix      string   `json:"suffix"`
	MaxTokens   int      `json:"max_tokens"`
	Temperature float64  `json:"temperature"`
	Stop        []string `json:"stop,omitempty"`
}

// FIMCompletionResponse represents the response from the FIM completion endpoint.
type FIMCompletionResponse struct {
	ID      string                        `json:"id"`
	Object  string                        `json:"object"`
	Created int                           `json:"created"`
	Model   string                        `json:"model"`
	Choices []FIMCompletionResponseChoice `json:"choices"`
	Usage   UsageInfo                     `json:"usage"`
}

// FIMCompletionResponseChoice represents a choice in the FIM completion response.
type FIMCompletionResponseChoice struct {
	Index        int          `json:"index"`
	Message      ChatMessage  `json:"message"`
	FinishReason FinishReason `json:"finish_reason,omitempty"`
}

// FIM sends a FIM request and returns the completion response.
func (c *MistralClient) FIM(params *FIMRequestParams) (*FIMCompletionResponse, error) {
	requestData := map[string]interface{}{
		"model":       params.Model,
		"prompt":      params.Prompt,
		"suffix":      params.Suffix,
		"max_tokens":  params.MaxTokens,
		"temperature": params.Temperature,
	}

	if params.Stop != nil {
		requestData["stop"] = params.Stop
	}

	response, err := c.request(http.MethodPost, requestData, "v1/fim/completions", false, nil)
	if err != nil {
		return nil, err
	}

	respData, ok := response.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response type: %T", response)
	}

	var fimResponse FIMCompletionResponse
	err = mapToStruct(respData, &fimResponse)
	if err != nil {
		return nil, err
	}

	return &fimResponse, nil
}
