package mistral

import (
	"fmt"
	"net/http"
)

// EmbeddingObject represents an embedding object in the response.
type EmbeddingObject struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingResponse represents the response from the embeddings endpoint.
type EmbeddingResponse struct {
	ID     string            `json:"id"`
	Object string            `json:"object"`
	Data   []EmbeddingObject `json:"data"`
	Model  string            `json:"model"`
	Usage  UsageInfo         `json:"usage"`
}

func (c *MistralClient) Embeddings(model string, input []string) (*EmbeddingResponse, error) {
	requestData := map[string]interface{}{
		"model": model,
		"input": input,
	}

	response, err := c.request(http.MethodPost, requestData, "v1/embeddings", false, nil)
	if err != nil {
		return nil, err
	}

	respData, ok := response.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response type: %T", response)
	}

	var embeddingResponse EmbeddingResponse
	err = mapToStruct(respData, &embeddingResponse)
	if err != nil {
		return nil, err
	}

	return &embeddingResponse, nil
}
