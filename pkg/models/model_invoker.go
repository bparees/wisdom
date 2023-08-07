package models

import (
	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/filters"
)

func InvokeModel(input api.ModelInput, model api.Model, filter filters.Filter) (*api.ModelResponse, error) {

	response, err := model.Invoke(input)
	if err != nil {
		return response, err
	}
	output, err := filter.FilterResponse(response)
	return output, err
}
