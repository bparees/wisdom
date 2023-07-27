package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func invokeModel(email, apiKey, prompt string) (*IBMModelResponsePayload, error) {
	// Create the JSON payload
	payload := IBMModelRequestPayload{
		Prompt:  prompt,
		ModelID: "L3Byb2plY3RzL2dyYW5pdGUvZzJiX2xyNWVuMDZfbWNsMTAyNF9jNDBr",
		TaskID:  "yaml-only-output",
		Mode:    "synchronous",
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		//fmt.Println("Error encoding JSON:", err)
		return nil, err
	}

	apiURL := "https://wca.wisdomforocp-cf7808d3396a7c1915bd1818afbfb3c0-0000.us-south.containers.appdomain.cloud/api/v1/jobs"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		//fmt.Println("Error creating HTTP request:", err)
		return nil, err
	}

	// Set the "Content-Type" header to "application/json"
	req.Header.Set("Content-Type", "application/json")

	// Set the "Authorization" header with the bearer token
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Email", email)

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
	return &apiResp, nil
}
