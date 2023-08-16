package api

type Model interface {
	Invoke(ModelInput) (*ModelResponse, error)
}

// ModelInput represents the payload for the prompt_request endpoint.
type ModelInput struct {
	UserId         string `json:"userid"`
	APIKey         string `json:"apikey"`
	ModelId        string `json:"modelId"`
	Provider       string `json:"provider"`
	Prompt         string `json:"prompt"`
	Context        string `json:"context"`
	ConversationID string `json:"conversationId"`
}

type ModelResponse struct {
	Input          string `json:"input_tokens"`
	Status         string `json:"status"`
	RequestID      string `json:"requestId"`
	ConversationID string `json:"conversationId"`
	Output         string `json:"output"`
	RawOutput      string `json:"raw_output"`
	ErrorMessage   string `json:"error"`
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

type ModelConfig struct {
	UserId string `yaml:"userId"`
	APIKey string `yaml:"apiKey"`

	Provider string `yaml:"provider"`
	ModelId  string `yaml:"modelId"`
	URL      string `yaml:"url"`
}

type ServerConfig struct {
	TLSCertFile  string   `yaml:"tlsCertFile"`
	TLSKeyFile   string   `yaml:"tlsKeyFile"`
	BearerTokens []string `yaml:"bearerTokens"`
}

type Config struct {
	Models          []ModelConfig `yaml:"models"`
	ServerConfig    ServerConfig  `yaml:"serverConfig"`
	DefaultProvider string        `yaml:"defaultProvider"`
	DefaultModelId  string        `yaml:"defaultModelId"`
}
