package main

func invokeModel(input ModelInput, model Model, filter Filter) (*ModelResponse, error) {

	response, err := model.Invoke(input)
	if err != nil {
		return response, err
	}
	output, err := filter.FilterResponse(response)
	return output, err
}
