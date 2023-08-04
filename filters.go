package main

import (
	"gopkg.in/yaml.v2"
)

type Filter struct {
	inputFilterChain    []InputFilter
	responseFilterChain []ResponseFilter
}
type InputFilter func(input *ModelInput) (*ModelInput, error)
type ResponseFilter func(response *ModelResponse) (*ModelResponse, error)

func NewFilter() Filter {
	filter := Filter{}
	//filter.responseFilterChain = append(filter.responseFilterChain, YamlLinter)
	return filter
}

func (f *Filter) FilterInput(input *ModelInput) (*ModelInput, error) {
	output := input
	var err error
	for _, filter := range f.inputFilterChain {
		output, err = filter(output)
		if err != nil {
			return output, err
		}
	}
	return output, err
}

func (f *Filter) FilterResponse(response *ModelResponse) (*ModelResponse, error) {
	output := response
	var err error
	for _, filter := range f.responseFilterChain {
		output, err = filter(output)
		if err != nil {
			return output, err
		}
	}
	return output, err
}

func YamlLinter(response *ModelResponse) (*ModelResponse, error) {
	return response, isValidYAML(response.Output)
}

func isValidYAML(yamlString string) error {
	var data interface{}
	err := yaml.Unmarshal([]byte(yamlString), &data)
	return err
}
