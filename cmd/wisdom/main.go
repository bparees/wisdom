package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/filters"
	"github.com/openshift/wisdom/pkg/model"
	"github.com/openshift/wisdom/pkg/model/ibm"
	"github.com/openshift/wisdom/pkg/model/openai"
	"github.com/openshift/wisdom/pkg/server"
)

func init() {
}

func main() {

	var rootCmd = &cobra.Command{Long: "Runs an inference router to proxy between user facing clients and LLMs"}

	rootCmd.AddCommand(newStartServerCommand())
	rootCmd.AddCommand(newInferCommand())
	rootCmd.Execute()

}

type options struct {
	email    string
	apiKey   string
	prompt   string
	provider string
	model    string
}

func newStartServerCommand() *cobra.Command {
	o := options{}

	var cmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := mux.NewRouter()

			h := server.Handler{
				UserId:          o.email,
				APIKey:          o.apiKey,
				Filter:          filters.NewFilter(),
				DefaultProvider: o.provider,
				DefaultModel:    o.model,
				Models:          createModels(),
			}

			r.HandleFunc("/prompt_request", h.PromptRequestHandler).Methods("POST")
			r.HandleFunc("/feedback", h.FeedbackHandler).Methods("POST")

			fmt.Println("Server listening on port 8080...")
			fmt.Printf("Default provider: %s\n", h.DefaultProvider)
			fmt.Printf("Default model: %s\n", h.DefaultModel)
			http.ListenAndServe(":8080", r)
			return nil
		},
	}

	cmd.Flags().StringVarP(&o.email, "email", "e", "", "Model email address used when not provided in the request")
	cmd.Flags().StringVarP(&o.apiKey, "apikey", "a", "", "Model API key used when not provided in the request")
	cmd.Flags().StringVarP(&o.provider, "provider", "p", "ibm", "Which LLM provider to use when not provided in the request.  Value values are: ibm, openai")
	cmd.Flags().StringVarP(&o.model, "model", "m", "L3Byb2plY3RzL2czYmNfc3RhY2tfc3RnMl9lcG9jaDNfanVsXzMx", "Which LLM model to use from the chosen provider when not provided in the request.  Valid values depend on the chosen provider.")

	return cmd

}

func newInferCommand() *cobra.Command {
	o := options{}

	var cmd = &cobra.Command{
		Use:   "infer",
		Short: "Do a single inference",
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.email == "" {
				return fmt.Errorf("user email address is required")
			}
			if o.apiKey == "" {
				return fmt.Errorf("API key is required")
			}
			if o.prompt == "" {
				return fmt.Errorf("model prompt is required")
			}
			filter := filters.NewFilter()
			m, err := getModel(o.provider, o.model)
			if err != nil {
				return err
			}

			input := api.ModelInput{
				UserId: o.email,
				APIKey: o.apiKey,
				Prompt: o.prompt,
			}
			response, err := model.InvokeModel(input, m, filter)
			if err != nil {
				if response != nil && response.Output != "" {
					fmt.Printf("Response(Error):\n%s\n", response.Output)
				}
				return fmt.Errorf("error invoking the LLM: %v", err)
			}

			fmt.Printf("Response:\n%s\n", response.Output)

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.email, "email", "e", "", "User's email address")
	cmd.Flags().StringVarP(&o.apiKey, "apikey", "a", "", "User's API key")
	cmd.Flags().StringVarP(&o.provider, "provider", "p", "ibm", "Which backend LLM provider to use.  Valid values are: ibm, openai")
	cmd.Flags().StringVarP(&o.model, "model", "m", "L3Byb2plY3RzL2czYmNfc3RhY2tfc3RnMl9lcG9jaDNfanVsXzMx", "Which LLM model to use from the provider.  Valid values depend on the chosen provider.")
	cmd.Flags().StringVarP(&o.prompt, "inference", "i", "", "Model prompt to be inferred")

	return cmd

}

func createModels() map[string]api.Model {
	models := make(map[string]api.Model)
	models[ibm.PROVIDER_ID+"|"+ibm.MODEL_ID] = ibm.NewIBMModel(ibm.MODEL_ID, "https://wca.wisdomforocp-cf7808d3396a7c1915bd1818afbfb3c0-0000.us-south.containers.appdomain.cloud")
	models[openai.PROVIDER_ID+"|"+openai.MODEL_ID] = openai.NewOpenAIModel(openai.MODEL_ID, "https://api.openai.com")
	return models
}

func getModel(provider, modelId string) (api.Model, error) {
	var model api.Model
	switch strings.ToLower(provider) {
	case ibm.PROVIDER_ID:
		model = ibm.NewIBMModel(modelId, "https://wca.wisdomforocp-cf7808d3396a7c1915bd1818afbfb3c0-0000.us-south.containers.appdomain.cloud")
	case openai.PROVIDER_ID:
		model = openai.NewOpenAIModel(modelId, "https://api.openai.com")
	default:
		return nil, fmt.Errorf("invalid provider specified: %s\nValid values are [ibm,openai]", provider)
	}
	return model, nil
}
