package ibm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/filters/markdown"
	"github.com/openshift/wisdom/pkg/filters/yaml"
)

type IBMModelRequestPayload struct {
	Prompt  string `json:"prompt"`
	ModelID string `json:"model_id"`
	TaskID  string `json:"task_id"`
	Mode    string `json:"mode"`
}

type IBMModelResponsePayload struct {
	AllTokens   string `json:"all_tokens"`
	InputTokens string `json:"input_tokens"`
	JobID       string `json:"job_id"`
	Model       string `json:"model"`
	Status      string `json:"status"`
	TaskID      string `json:"task_id"`
	TaskOutput  string `json:"task_output"`
}

type IBMModel struct {
	modelId string
	url     string
	apiKey  string
	userId  string
	filter  api.Filter
}

func NewIBMModel(modelId, url, userId, apiKey string) *IBMModel {
	filter := api.NewFilter(nil, []api.ResponseFilter{markdown.MarkdownStripper, yaml.YamlLinter})
	return &IBMModel{
		modelId: modelId,
		url:     url,
		apiKey:  apiKey,
		userId:  userId,
		filter:  filter,
	}
}

func (m *IBMModel) GetFilter() api.Filter {
	return m.filter
}

func (m *IBMModel) Invoke(input api.ModelInput) (api.ModelResponse, error) {

	if input.UserId == "" && m.userId == "" {
		return api.ModelResponse{}, fmt.Errorf("user email address is required, none provided")
	}
	if input.APIKey == "" && m.apiKey == "" {
		return api.ModelResponse{}, fmt.Errorf("api key is required, none provided")
	}

	apiKey, userId := m.apiKey, m.userId
	if input.APIKey != "" {
		apiKey = input.APIKey
	}
	if input.UserId != "" {
		userId = input.UserId
	}

	payload := IBMModelRequestPayload{
		Prompt:  input.Prompt,
		ModelID: m.modelId,
		TaskID:  "yaml-only-raw-output",
		Mode:    "synchronous",
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		//fmt.Println("Error encoding JSON:", err)
		return api.ModelResponse{}, err
	}

	apiURL := m.url + "/api/v1/jobs"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		//fmt.Println("Error creating HTTP request:", err)
		return api.ModelResponse{}, err
	}

	// Set the "Content-Type" header to "application/json"
	req.Header.Set("Content-Type", "application/json")

	// Set the "Authorization" header with the bearer token
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Email", userId)

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
		return api.ModelResponse{}, fmt.Errorf("API request failed with status: %v", resp.Status)
	}

	// Parse the JSON response into the APIResponse struct
	var apiResp IBMModelResponsePayload
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		//fmt.Println("Error decoding API response:", err)
		return api.ModelResponse{}, err
	}
	response := api.ModelResponse{}
	response.Input = input.Prompt
	response.Output = apiResp.TaskOutput
	response.RawOutput = apiResp.AllTokens
	//output := apiResp.AllTokens[len(apiResp.InputTokens):]
	response.RequestID = apiResp.JobID

	return response, err
}

func (m *IBMModel) FilterInput(input api.ModelInput) (api.ModelInput, error) {
	return m.filter.FilterInput(input)
}
