package filters

import (
	"github.com/openshift/wisdom/pkg/filters/yaml"

	"github.com/openshift/wisdom/pkg/api"
)

type Filter struct {
	inputFilterChain    []InputFilter
	responseFilterChain []ResponseFilter
}
type InputFilter func(input *api.ModelInput) (*api.ModelInput, error)
type ResponseFilter func(response *api.ModelResponse) (*api.ModelResponse, error)

func NewFilter() Filter {
	filter := Filter{}
	//filter.responseFilterChain = append(filter.responseFilterChain, markdown.MarkdownStripper, yaml.YamlLinter)
	filter.responseFilterChain = append(filter.responseFilterChain, yaml.YamlLinter)
	return filter
}

func (f *Filter) FilterInput(input *api.ModelInput) (*api.ModelInput, error) {
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

func (f *Filter) FilterResponse(response *api.ModelResponse) (*api.ModelResponse, error) {
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
