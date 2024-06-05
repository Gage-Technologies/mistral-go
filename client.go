package mistral

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	Endpoint          = "https://api.mistral.ai"
	CodestralEndpoint = "https://codestral.mistral.ai"
	DefaultMaxRetries = 5
	DefaultTimeout    = 120 * time.Second
)

var retryStatusCodes = map[int]bool{
	429: true,
	500: true,
	502: true,
	503: true,
	504: true,
}

type MistralClient struct {
	apiKey     string
	endpoint   string
	maxRetries int
	timeout    time.Duration
}

func NewMistralClient(apiKey string, endpoint string, maxRetries int, timeout time.Duration) *MistralClient {
	if apiKey == "" {
		apiKey = os.Getenv("MISTRAL_API_KEY")
	}
	if endpoint == "" {
		endpoint = Endpoint
	}
	if maxRetries == 0 {
		maxRetries = DefaultMaxRetries
	}
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	return &MistralClient{
		apiKey:     apiKey,
		endpoint:   endpoint,
		maxRetries: maxRetries,
		timeout:    timeout,
	}
}

// NewMistralClientDefault creates a new Mistral API client with the default endpoint and the given API key. Defaults to using MISTRAL_API_KEY from the environment.
func NewMistralClientDefault(apiKey string) *MistralClient {
	if apiKey == "" {
		apiKey = os.Getenv("MISTRAL_API_KEY")
	}

	return NewMistralClient(apiKey, Endpoint, DefaultMaxRetries, DefaultTimeout)
}

// NewCodestralClientDefault creates a new Codestral API client with the default endpoint and the given API key. Defaults to using CODESTRAL_API_KEY from the environment.
func NewCodestralClientDefault(apiKey string) *MistralClient {
	if apiKey == "" {
		apiKey = os.Getenv("CODESTRAL_API_KEY")
	}

	return NewMistralClient(apiKey, CodestralEndpoint, DefaultMaxRetries, DefaultTimeout)
}

func (c *MistralClient) request(method string, jsonData map[string]interface{}, path string, stream bool, params map[string]string) (interface{}, error) {
	uri, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}
	uri.Path = path
	jsonValue, _ := json.Marshal(jsonData)
	req, err := http.NewRequest(method, uri.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: c.timeout,
	}

	var resp *http.Response
	for i := 0; i < c.maxRetries; i++ {
		resp, err = client.Do(req)
		if err != nil {
			if i == c.maxRetries-1 {
				return nil, err
			}
			continue
		}
		if _, ok := retryStatusCodes[resp.StatusCode]; ok {
			time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
			continue
		}
		break
	}

	if resp.StatusCode >= 400 {
		responseBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("(HTTP Error %d) %s", resp.StatusCode, string(responseBytes))
	}

	if stream {
		return resp.Body, nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
