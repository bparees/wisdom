package main

// ResponseType is an enum for valid response types.
type ResponseType string

const (
	YAML    ResponseType = "yaml"
	Natural ResponseType = "natural"
	Command ResponseType = "command"
)

type IBMModelRequestPayload struct {
	Prompt  string `json:"prompt"`
	ModelID string `json:"model_id"`
	TaskID  string `json:"task_id"`
	Mode    string `json:"mode"`
}

type IBMModelResponsePayload struct {
	AllTokens    string `json:"all_tokens"`
	InputTokens  string `json:"input_tokens"`
	JobID        string `json:"job_id"`
	Model        string `json:"model"`
	Status       string `json:"status"`
	TaskID       string `json:"task_id"`
	TaskOutput   string `json:"task_output"`
	OutputTokens string `json:"output_tokens"`
}

// PromptInputPayload represents the payload for the prompt_request endpoint.
type PromptInputPayload struct {
	Prompt         string       `json:"prompt"`
	Context        string       `json:"context"`
	ConversationID string       `json:"conversationId"`
	ResponseType   ResponseType `json:"responseType"`
}

// FeedbackPayload represents the payload for the feedback endpoint.
type FeedbackPayload struct {
	RequestID         string `json:"requestId"`
	Response          string `json:"response"`
	ConversationID    string `json:"conversationId"`
	ResponseAccepted  bool   `json:"responseAccepted"`
	CorrectedResponse string `json:"correctedResponse"`
	UserComments      string `json:"userComments"`
}
