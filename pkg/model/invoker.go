package model

import (
	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/filters"
)

func InvokeModel(input api.ModelInput, model api.Model, filter filters.Filter) (*api.ModelResponse, error) {

	response, err := model.Invoke(input)
	if response == nil {
		response = &api.ModelResponse{}
	}
	if err != nil {
		response.ErrorMessage = err.Error()
		return response, err
	}
	output, err := filter.FilterResponse(response)
	if err != nil {
		response.ErrorMessage = err.Error()
	}
	return output, err
}
