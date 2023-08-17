package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/openshift/wisdom/pkg/api"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.AuthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) HandleGithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := h.AuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Errorf("failed to exchange token: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := h.AuthConfig.Client(context.Background(), token)
	// Make an authenticated API request to get the user's information
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Errorf("failed to get user information: %v", err)
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Errorf("failed to decode response: %v", err)
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	userID := user["login"] // Extract the user's GitHub ID

	// Create a session for the user
	session, err := h.CookieStore.Get(r, "wisdom-session")
	if err != nil {
		log.Errorf("failed to create session: %v", err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Store the OAuth token in the session
	session.Values["username"] = userID
	err = session.Save(r, w)

	if err != nil {
		log.Errorf("failed to save session: %v", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	log.Debugf("stored session for user %s", userID)

	http.Redirect(w, r, "/apitoken", http.StatusSeeOther)
}

func (h *Handler) HandleApiToken(w http.ResponseWriter, r *http.Request) {
	username, authorized := h.isAuthorized(w, r)
	if !authorized {
		log.Debugf("user %q not authorized to get api token", username)
		url := r.URL
		url.Path = "/login"
		http.Redirect(w, r, url.String(), http.StatusFound)
		return
	}

	claims := &api.Claims{
		Username:         username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			//ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(h.TokenEncryptionKey)
	if err != nil {
		log.Errorf("failed to sign token: %v", err)
		url := r.URL
		url.Path = "/login"
		http.Redirect(w, r, url.String(), http.StatusFound)
		return
	}
	apiToken := api.APIToken{
		Token: tokenString,
	}
	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(apiToken)

	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (h *Handler) isAuthorized(w http.ResponseWriter, r *http.Request) (string, bool) {
	// Retrieve the session for the user
	session, err := h.CookieStore.Get(r, "wisdom-session")
	if err != nil {
		log.Errorf("failed to get session: %v", err)
		url := r.URL
		url.Path = "/login"
		http.Redirect(w, r, url.String(), http.StatusFound)
		return "", false
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		log.Errorf("session does not contain a valid username")
		url := r.URL
		url.Path = "/login"
		http.Redirect(w, r, url.String(), http.StatusFound)
		return "", false
	}

	// Check if the username is in the allow list
	if _, found := h.AllowedUsers[username]; !found {
		log.Debugf("user %s not in allow list", username)
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return username, false
	}
	return username, true
}

func (h *Handler) hasValidBearerToken(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Debug("no authorization header")
		return false
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Debug("authorization header does not start with Bearer")
		return false
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &api.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return h.TokenEncryptionKey, nil
	})
	if err != nil {
		log.Errorf("failed to parse jwt token: %v", err)
		if err == jwt.ErrSignatureInvalid {
			return false
		}
		return false
	}
	if !tkn.Valid {
		log.Debug("token is not valid")
		return false
	}
	if _, found := h.AllowedUsers[claims.Username]; !found {
		log.Debugf("user %s not in allowed users", claims.Username)
		return false
	}
	return true
}
