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
