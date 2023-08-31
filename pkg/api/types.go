package api

import (
	"github.com/golang-jwt/jwt/v4"
)

type Filter struct {
	InputFilterChain    []InputFilter
	ResponseFilterChain []ResponseFilter
}

type InputFilter func(input ModelInput) (ModelInput, error)
type ResponseFilter func(response ModelResponse) (ModelResponse, error)

func NewFilter(inputFilters []InputFilter, responseFilters []ResponseFilter) Filter {
	filter := Filter{
		InputFilterChain:    inputFilters,
		ResponseFilterChain: responseFilters,
	}
	return filter
}

func (f Filter) FilterInput(input ModelInput) (ModelInput, error) {
	output := input
	var err error
	for _, filter := range f.InputFilterChain {
		output, err = filter(output)
		if err != nil {
			return output, err
		}
	}
	return output, err
}

func (f Filter) FilterResponse(response ModelResponse) (ModelResponse, error) {
	output := response
	var err error
	for _, filter := range f.ResponseFilterChain {
		output, err = filter(output)
		if err != nil {
			return output, err
		}
	}
	return output, err
}

type Model interface {
	Invoke(ModelInput) (ModelResponse, error)
	GetFilter() Filter
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
	Error          string `json:"error"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type APIToken struct {
	Token string `json:"token"`
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
	TLSCertFile          string          `yaml:"tlsCertFile"`
	TLSKeyFile           string          `yaml:"tlsKeyFile"`
	ClientID             string          `yaml:"clientId"`
	ClientSecret         string          `yaml:"clientSecret"`
	RedirectURL          string          `yaml:"redirectUrl"`
	SessionAuthKey       string          `yaml:"sessionAuthKey"`
	SessionEncryptionKey string          `yaml:"sessionEncryptionKey"`
	TokenEncryptionKey   string          `yaml:"tokenEncryptionKey"`
	AllowedUsers         map[string]bool `yaml:"allowedUsers"`
}

type Config struct {
	Models          []ModelConfig `yaml:"models"`
	ServerConfig    ServerConfig  `yaml:"serverConfig"`
	DefaultProvider string        `yaml:"defaultProvider"`
	DefaultModelId  string        `yaml:"defaultModelId"`
}
