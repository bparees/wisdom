package server

import (
	"github.com/gorilla/sessions"
	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/filters"
	"golang.org/x/oauth2"
)

type Handler struct {
	Filter          filters.Filter
	DefaultModel    string
	DefaultProvider string
	Models          map[string]api.Model
	ClientID        string
	ClientSecret    string
	//SessionAuthKey       string
	//SessionEncryptionKey string
	AuthConfig         oauth2.Config
	CookieStore        *sessions.CookieStore
	TokenEncryptionKey []byte
	AllowedUsers       map[string]bool
}
