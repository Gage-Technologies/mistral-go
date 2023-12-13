package mistral

import (
	"fmt"
	"net/http"
)

// ModelPermission represents the permissions of a model.
type ModelPermission struct {
	ID                 string `json:"id"`
	Object             string `json:"object"`
	Created            int    `json:"created"`
	AllowCreateEngine  bool   `json:"allow_create_engine"`
	AllowSampling      bool   `json:"allow_sampling"`
	AllowLogprobs      bool   `json:"allow_logprobs"`
	AllowSearchIndices bool   `json:"allow_search_indices"`
	AllowView          bool   `json:"allow_view"`
	AllowFineTuning    bool   `json:"allow_fine_tuning"`
	Organization       string `json:"organization"`
	Group              string `json:"group,omitempty"`
	IsBlocking         bool   `json:"is_blocking"`
}

// ModelCard represents a model card.
type ModelCard struct {
	ID         string            `json:"id"`
	Object     string            `json:"object"`
	Created    int               `json:"created"`
	OwnedBy    string            `json:"owned_by"`
	Root       string            `json:"root,omitempty"`
	Parent     string            `json:"parent,omitempty"`
	Permission []ModelPermission `json:"permission"`
}

// ModelList represents a list of models.
type ModelList struct {
	Object string      `json:"object"`
	Data   []ModelCard `json:"data"`
}

func (c *MistralClient) ListModels() (*ModelList, error) {
	response, err := c.request(http.MethodGet, nil, "v1/models", false, nil)
	if err != nil {
		return nil, err
	}

	respData, ok := response.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response type: %T", response)
	}

	var modelList ModelList
	err = mapToStruct(respData, &modelList)
	if err != nil {
		return nil, err
	}

	return &modelList, nil
}
