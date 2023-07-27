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
	rootCmd.Execute()

}

type options struct {
	email  string
	apiKey string
}

func newStartServerCommand() *cobra.Command {
	o := options{}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.email == "" {
				return fmt.Errorf("User email address is required")
			}
			if o.apiKey == "" {
				return fmt.Errorf("API key is required")
			}
			r := mux.NewRouter()

			h := Handler{
				email:  o.email,
				apiKey: o.apiKey,
			}
			r.HandleFunc("/prompt_request", h.PromptRequestHandler).Methods("POST")
			r.HandleFunc("/feedback", h.FeedbackHandler).Methods("POST")

			fmt.Println("Server listening on port 8080...")
			http.ListenAndServe(":8080", r)
			return nil
		},
	}

	startCmd.Flags().StringVarP(&o.email, "email", "e", "", "User's email address")
	startCmd.Flags().StringVarP(&o.apiKey, "apikey", "a", "", "User's API key")

	return startCmd

}
