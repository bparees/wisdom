package ibm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openshift/wisdom/pkg/api"
)

const (
	PROVIDER_ID = "ibm"
	MODEL_ID    = "L3Byb2plY3RzL2czYmNfc3RhY2tfc3RnMl9lcG9jaDNfanVsXzMx"
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
}

func NewIBMModel(modelId, url, userId, apiKey string) *IBMModel {
	return &IBMModel{
		modelId: modelId,
		url:     url,
		apiKey:  apiKey,
		userId:  userId,
	}
}

func (m *IBMModel) Invoke(input api.ModelInput) (*api.ModelResponse, error) {

	if input.UserId == "" && m.userId == "" {
		return nil, fmt.Errorf("user email address is required, none provided")
	}
	if input.APIKey == "" && m.apiKey == "" {
		return nil, fmt.Errorf("api key is required, none provided")
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
		TaskID:  "yaml-only-output",
		Mode:    "synchronous",
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		//fmt.Println("Error encoding JSON:", err)
		return nil, err
	}

	apiURL := m.url + "/api/v1/jobs"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		//fmt.Println("Error creating HTTP request:", err)
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		//fmt.Println("API request failed with status:", resp.Status)
		return nil, err
	}

	// Parse the JSON response into the APIResponse struct
	var apiResp IBMModelResponsePayload
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		fmt.Println("Error decoding API response:", err)
		return nil, err
	}
	response := api.ModelResponse{}
	response.Input = input.Prompt
	response.Output = apiResp.TaskOutput
	response.RawOutput = apiResp.AllTokens
	//output := apiResp.AllTokens[len(apiResp.InputTokens):]
	response.RequestID = apiResp.JobID

	return &response, err
}
