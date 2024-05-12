package mistral

const (
	ModelMistralLargeLatest  = "mistral-large-latest"
	ModelMistralMediumLatest = "mistral-medium-latest"
	ModelMistralSmallLatest  = "mistral-small-latest"
	ModelOpenMixtral8x7b     = "open-mixtral-8x7b"
	ModelOpenMistral7b       = "open-mistral-7b"

	ModelMistralLarge2402  = "mistral-large-2402"
	ModelMistralMedium2312 = "mistral-medium-2312"
	ModelMistralSmall2402  = "mistral-small-2402"
	ModelMistralSmall2312  = "mistral-small-2312"
	ModelMistralTiny       = "mistral-tiny-2312"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
	RoleTool      = "tool"
)

// FinishReason the reason that a chat message was finished
type FinishReason string

const (
	FinishReasonStop   FinishReason = "stop"
	FinishReasonLength FinishReason = "length"
	FinishReasonError  FinishReason = "error"
)

// ResponseFormat the format that the response must adhere to
type ResponseFormat string

const (
	ResponseFormatText       ResponseFormat = "text"
	ResponseFormatJsonObject ResponseFormat = "json_object"
)

// ToolType type of tool defined for the llm
type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

const (
	ToolChoiceAny  = "any"
	ToolChoiceAuto = "auto"
	ToolChoiceNone = "none"
)

// Tool definition of a tool that the llm can call
type Tool struct {
	Type     ToolType `json:"type"`
	Function Function `json:"function"`
}

// Function definition of a function that the llm can call including its parameters
type Function struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

// FunctionCall represents a request to call an external tool by the llm
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ToolCall represents the call to a tool by the llm
type ToolCall struct {
	Id       string       `json:"id"`
	Type     ToolType     `json:"type"`
	Function FunctionCall `json:"function"`
}

// DeltaMessage represents the delta between the prior state of the message and the new state of the message when streaming responses.
type DeltaMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

// ChatMessage represents a single message in a chat.
type ChatMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}
