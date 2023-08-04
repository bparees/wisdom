package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	email  string
	apiKey string
	filter Filter
	models map[string]Model
}

func (h *Handler) PromptRequestHandler(w http.ResponseWriter, r *http.Request) {
	var payload ModelInput
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Handle the prompt_request request here and generate the response based on the payload
	//response := fmt.Sprintf("Received prompt: %s\n", payload.Prompt)
	fmt.Printf("Running inference for prompt: %s\n", payload.Prompt)

	input := ModelInput{
		UserId: h.email,
		APIKey: h.apiKey,
		Prompt: payload.Prompt,
	}

	model, found := h.models[payload.ModelId]
	if !found {
		// error model not found.
	}
	response, err := invokeModel(input, model, h.filter)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.Output))
}

func (h *Handler) FeedbackHandler(w http.ResponseWriter, r *http.Request) {
	var payload FeedbackPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Handle the feedback request here

	response := "Feedback received."

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
