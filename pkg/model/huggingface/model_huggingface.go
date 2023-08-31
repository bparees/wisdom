package ibm

import (
	"context"
	"fmt"

	huggingface "github.com/hupe1980/go-huggingface"

	"github.com/openshift/wisdom/pkg/api"
)

type HFModel struct {
	modelId string
	url     string
	apiKey  string
	filter  api.Filter
}

func NewHFModel(modelId, url, apiKey string) *HFModel {
	//filter := api.NewFilter(nil, []api.ResponseFilter{markdown.MarkdownStripper, yaml.YamlLinter})
	filter := api.Filter{}

	return &HFModel{
		modelId: modelId,
		url:     url,
		apiKey:  apiKey,
		filter:  filter,
	}
}

func (m *HFModel) GetFilter() api.Filter {
	return m.filter
}

func (m *HFModel) Invoke(input api.ModelInput) (api.ModelResponse, error) {

	if input.APIKey == "" && m.apiKey == "" {
		return api.ModelResponse{}, fmt.Errorf("api key is required, none provided")
	}

	apiKey := m.apiKey
	if input.APIKey != "" {
		apiKey = input.APIKey
	}
	client := huggingface.NewInferenceClient(apiKey)
	client.SetModel(m.modelId)

	req := &huggingface.TextGenerationRequest{
		Inputs: input.Prompt,
		Model:  input.ModelId,
	}

	a := 100
	req.Parameters.MaxNewTokens = &a
	b := 30.0
	req.Parameters.MaxTime = &b
	c := 2
	req.Parameters.NumReturnSequences = &c

	resp, err := client.TextGeneration(context.Background(), req)
	if err != nil {
		return api.ModelResponse{}, fmt.Errorf("error making api request: %v", err)
	}

	response := api.ModelResponse{}
	response.Input = input.Prompt
	response.Output = resp.GeneratedText

	return response, err
}
