package model

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/openshift/wisdom/pkg/api"
)

func InvokeModel(input api.ModelInput, model api.Model) (api.ModelResponse, error) {
	log.Debugf("model input:\n%#v", input)
	input, err := model.GetFilter().FilterInput(input)
	if err != nil {
		return api.ModelResponse{}, fmt.Errorf("error filtering input: %s", err)
	}
	log.Debugf("model filtered input:\n%#v", input)
	response, err := model.Invoke(input)
	log.Debugf("model response:\n%#v\nerror: %v", response, err)
	if err != nil {
		response.Error = err.Error()
		return response, err
	}

	output, err := model.GetFilter().FilterResponse(response)
	if err != nil {
		response.Error = err.Error()
	}
	log.Debugf("model filtered output:\n%#v", output)
	return output, err
}
