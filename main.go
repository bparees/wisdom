package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
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
	email  string
	apiKey string
	prompt string
}

func newStartServerCommand() *cobra.Command {
	o := options{}

	var cmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.email == "" {
				return fmt.Errorf("user email address is required")
			}
			if o.apiKey == "" {
				return fmt.Errorf("API key is required")
			}
			r := mux.NewRouter()

			h := Handler{
				email:  o.email,
				apiKey: o.apiKey,
				filter: NewFilter(),
			}
			r.HandleFunc("/prompt_request", h.PromptRequestHandler).Methods("POST")
			r.HandleFunc("/feedback", h.FeedbackHandler).Methods("POST")

			fmt.Println("Server listening on port 8080...")
			http.ListenAndServe(":8080", r)
			return nil
		},
	}

	cmd.Flags().StringVarP(&o.email, "email", "e", "", "User's email address")
	cmd.Flags().StringVarP(&o.apiKey, "apikey", "a", "", "User's API key")

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
			filter := NewFilter()
			//model := NewIBMModel("L3Byb2plY3RzL2czYmNfc3RhY2tfc3RnMl9lcG9jaDNfanVsXzMx", "https://wca.wisdomforocp-cf7808d3396a7c1915bd1818afbfb3c0-0000.us-south.containers.appdomain.cloud")
			model := NewOpenAIModel("gpt-3.5-turbo", "https://api.openai.com")

			input := ModelInput{
				UserId: o.email,
				APIKey: o.apiKey,
				Prompt: o.prompt,
			}
			response, err := invokeModel(input, model, filter)
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
	cmd.Flags().StringVarP(&o.prompt, "prompt", "p", "", "Model prompt")

	return cmd

}
