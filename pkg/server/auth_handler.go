package server

import (
	"net/http"
	"strings"
)

func (h *Handler) hasValidBearerToken(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	return h.BearerTokens[token]
}
