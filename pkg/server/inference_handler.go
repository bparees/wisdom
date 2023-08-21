package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/model"
)

func (h *Handler) InferHandler(w http.ResponseWriter, r *http.Request) {
	if !h.hasValidBearerToken(r) {
		http.Error(w, "No valid bearer token found", http.StatusUnauthorized)
		return
	}
	var payload api.ModelInput
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Printf("Running inference for prompt: %s\n", payload.Prompt)

	if payload.Provider == "" {
		payload.Provider = h.DefaultProvider
	}
	if payload.ModelId == "" {
		payload.ModelId = h.DefaultModel
	}
	m, found := h.Models[payload.Provider+"/"+payload.ModelId]
	if !found {
		http.Error(w, fmt.Sprintf("Invalid provider/model: %s|%s", payload.Provider, payload.ModelId), http.StatusBadRequest)
		return
	}
	response, err := model.InvokeModel(payload, m, h.Filter)

	buf := bytes.Buffer{}
	if response != nil {
		err = json.NewEncoder(&buf).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte(err.Error()))
			return
		}
	}

	w.Header().Set("Content-Type", "text/json")
	if err != nil || (response != nil && response.ErrorMessage != "") {
		w.WriteHeader(http.StatusExpectationFailed)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Write(buf.Bytes())
}

/*
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
*/
