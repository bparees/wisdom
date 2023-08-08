package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/filters"
	"github.com/openshift/wisdom/pkg/model"
)

type Handler struct {
	UserId          string
	APIKey          string
	Filter          filters.Filter
	DefaultModel    string
	DefaultProvider string
	Models          map[string]api.Model
}

func (h *Handler) PromptRequestHandler(w http.ResponseWriter, r *http.Request) {
	var payload api.ModelInput
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Handle the prompt_request request here and generate the response based on the payload
	//response := fmt.Sprintf("Received prompt: %s\n", payload.Prompt)
	fmt.Printf("Running inference for prompt: %s\n", payload.Prompt)

	if payload.UserId == "" {
		payload.UserId = h.UserId
	}
	if payload.APIKey == "" {
		payload.APIKey = h.APIKey
	}
	if payload.ModelId == "" {
		payload.ModelId = h.DefaultModel
	}
	if payload.Provider == "" {
		payload.Provider = h.DefaultProvider
	}

	m, found := h.Models[payload.Provider+"|"+payload.ModelId]
	if !found {
		http.Error(w, fmt.Sprintf("Invalid provider/model: %s|%s", payload.Provider, payload.ModelId), http.StatusBadRequest)
		return
	}
	response, err := model.InvokeModel(payload, m, h.Filter)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(err.Error()))
		return
	}

	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (h *Handler) FeedbackHandler(w http.ResponseWriter, r *http.Request) {
	var payload api.FeedbackPayload
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
