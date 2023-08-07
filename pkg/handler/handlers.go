package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/filters"
	"github.com/openshift/wisdom/pkg/models"
)

type Handler struct {
	Email     string
	APIKey    string
	Filter    filters.Filter
	Providers map[string]api.Model
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

	input := api.ModelInput{
		UserId: h.Email,
		APIKey: h.APIKey,
		Prompt: payload.Prompt,
	}

	model, found := h.Providers[payload.ModelId]
	if !found {
		http.Error(w, fmt.Sprintf("Invalid provider/model: %s", payload.ModelId), http.StatusBadRequest)
		return
	}
	response, err := models.InvokeModel(input, model, h.Filter)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.Output))
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
