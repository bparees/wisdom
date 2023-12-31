package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/model"
)

func (h *Handler) CORSHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "Access-Control-Allow-Origin", "*")
	http.Header.Add(w.Header(), "Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	http.Header.Add(w.Header(), "Access-Control-Allow-Headers", "Content-Type, Authorization")

}

func (h *Handler) InferHandler(w http.ResponseWriter, r *http.Request) {

	http.Header.Add(w.Header(), "Access-Control-Allow-Origin", "*")

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

	log.Debugf("Using provider/model %s/%s for prompt:\n%s\n", payload.Provider, payload.ModelId, payload.Prompt)

	response, err := model.InvokeModel(payload, m)
	if err != nil {
		log.Errorf("failed to invoke model: %v", err)
		http.Error(w, "Failed to invoke model", http.StatusInternalServerError)
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
	if err != nil || (response.Error != "") {
		log.Debugf("model invocation returning error: %v", err)
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
