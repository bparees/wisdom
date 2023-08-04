package main

type Model interface {
	Invoke(ModelInput) (*ModelResponse, error)
}

// ModelInput represents the payload for the prompt_request endpoint.
type ModelInput struct {
	UserId         string       `json:"userid"`
	APIKey         string       `json:"apikey"`
	ModelId        string       `json:"modelid"`
	Prompt         string       `json:"prompt"`
	Context        string       `json:"context"`
	ConversationID string       `json:"conversationId"`
	ResponseType   ResponseType `json:"responseType"`
}

type ModelResponse struct {
	Input          string `json:"input_tokens"`
	Status         string `json:"status"`
	RequestID      string `json:"request_id"`
	ConversationID string `json:"conversation_id"`
	Output         string `json:"output"`
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
