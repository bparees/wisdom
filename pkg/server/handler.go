package server

import (
	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/filters"
)

type Handler struct {
	Filter          filters.Filter
	DefaultModel    string
	DefaultProvider string
	Models          map[string]api.Model
	BearerTokens    map[string]bool
}
