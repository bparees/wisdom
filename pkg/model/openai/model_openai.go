package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openshift/wisdom/pkg/api"
)

// OpenAI
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIModelRequestPayload struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

type OpenAIModelResponsePayload struct {
	ID      string `json:"id"`
	Choices []struct {
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
}

type OpenAIModel struct {
	modelId string
	url     string
	apiKey  string
	filter  api.Filter
}

func NewOpenAIModel(modelId, url, apiKey string) *OpenAIModel {
	filter := api.NewFilter(nil, nil)

	return &OpenAIModel{
		modelId: modelId,
		url:     url,
		apiKey:  apiKey,
		filter:  filter,
	}
}

func (m *OpenAIModel) GetFilter() api.Filter {
	return m.filter
}

func (m *OpenAIModel) Invoke(input api.ModelInput) (api.ModelResponse, error) {

	if input.APIKey == "" && m.apiKey == "" {
		return api.ModelResponse{}, fmt.Errorf("api key is required, none provided")
	}

	apiKey := m.apiKey
	if input.APIKey != "" {
		apiKey = input.APIKey
	}

	payload := OpenAIModelRequestPayload{
		Model: m.modelId,
	}
	payload.Messages = append(payload.Messages, OpenAIMessage{Role: "user", Content: input.Prompt})

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		//fmt.Println("Error encoding JSON:", err)
		return api.ModelResponse{}, err
	}

	apiURL := m.url + "/v1/chat/completions"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		//fmt.Println("Error creating HTTP request:", err)
		return api.ModelResponse{}, err
	}

	// Set the "Content-Type" header to "application/json"
	req.Header.Set("Content-Type", "application/json")

	// Set the "Authorization" header with the bearer token
	req.Header.Set("Authorization", "Bearer "+apiKey)
	//req.Header.Set("Email", input.UserId)

	// Make the API call
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println("Error making API request:", err)
		return api.ModelResponse{}, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return api.ModelResponse{}, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Parse the JSON response into the APIResponse struct
	var apiResp OpenAIModelResponsePayload
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		fmt.Println("Error decoding API response:", err)
		return api.ModelResponse{}, err
	}
	if len(apiResp.Choices) == 0 {
		return api.ModelResponse{}, fmt.Errorf("model returned no valid responses: %v", apiResp)
	}
	response := api.ModelResponse{}
	response.Input = input.Prompt
	response.Output = apiResp.Choices[0].Message.Content
	response.RawOutput = apiResp.Choices[0].Message.Content
	return response, err
}
