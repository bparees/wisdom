package model

import (
	log "github.com/sirupsen/logrus"

	"github.com/openshift/wisdom/pkg/api"
)

func InvokeModel(input api.ModelInput, model api.Model) (*api.ModelResponse, error) {
	response, err := model.Invoke(input)
	if response == nil {
		response = &api.ModelResponse{}
	}
	log.Debugf("model response: %#v", response)
	if err != nil {
		response.Error = err.Error()
		return response, err
	}

	output, err := model.GetFilter().FilterResponse(response)
	if err != nil {
		response.Error = err.Error()
	}
	return output, err
}
