package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"

	"github.com/openshift/wisdom/pkg/api"
	"github.com/openshift/wisdom/pkg/model"
	hf "github.com/openshift/wisdom/pkg/model/huggingface"
	"github.com/openshift/wisdom/pkg/model/ibm"
	"github.com/openshift/wisdom/pkg/model/openai"
	"github.com/openshift/wisdom/pkg/server"
)

var (
	models map[string]api.Model
)

func main() {

	var rootCmd = &cobra.Command{
		Long:         "Runs an inference router to proxy between user facing clients and LLMs",
		SilenceUsage: true,
	}

	rootCmd.AddCommand(newStartServerCommand())
	rootCmd.AddCommand(newInferCommand())
	rootCmd.Execute()

}

type options struct {
	configFile string
	verbosity  string
}

type inferOptions struct {
	options
	provider string
	modelId  string
	prompt   string
}

func loadConfig(filename string) (api.Config, error) {
	var config api.Config
	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	yamlParser := yaml.NewDecoder(configFile)
	err = yamlParser.Decode(&config)
	log.Debugf("Loaded config: %#v", config)
	return config, err
}

func newStartServerCommand() *cobra.Command {
	o := options{}

	var cmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			level, err := log.ParseLevel(o.verbosity)
			if err != nil {
				log.WithError(err).Fatal("Cannot parse log-level")
			}
			log.SetLevel(level)

			if o.configFile == "" {
				return fmt.Errorf("config file is required")
			}
			config, err := loadConfig(o.configFile)
			if err != nil {
				return fmt.Errorf("error loading configfile %s: %v", o.configFile, err)
			}
			r := mux.NewRouter()

			models = initModels(config)

			h := server.Handler{
				DefaultProvider: config.DefaultProvider,
				DefaultModel:    config.DefaultModelId,
				Models:          models,
				ClientID:        config.ServerConfig.ClientID,
				ClientSecret:    config.ServerConfig.ClientSecret,
				AllowedUsers:    config.ServerConfig.AllowedUsers,
			}
			tokenKey, err := base64.StdEncoding.DecodeString(config.ServerConfig.TokenEncryptionKey)
			if err != nil {
				return err
			}
			h.TokenEncryptionKey = tokenKey

			h.AuthConfig = oauth2.Config{
				ClientID:     h.ClientID,
				ClientSecret: h.ClientSecret,
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://github.com/login/oauth/authorize",
					TokenURL: "https://github.com/login/oauth/access_token",
				},
				//RedirectURL: "https://localhost:8443/githubcallback",
				RedirectURL: config.ServerConfig.RedirectURL,
				Scopes:      []string{"user:email"},
			}

			authKey, err := base64.StdEncoding.DecodeString(config.ServerConfig.SessionAuthKey)
			if err != nil {
				return err
			}
			encKey, err := base64.StdEncoding.DecodeString(config.ServerConfig.SessionEncryptionKey)
			if err != nil {
				return err
			}

			h.CookieStore = sessions.NewCookieStore(authKey, encKey)
			h.CookieStore.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   0,
				HttpOnly: false,
			}

			r.HandleFunc("/infer", h.InferHandler).Methods("POST")
			r.HandleFunc("/infer", h.CORSHandler).Methods("OPTIONS")
			//r.HandleFunc("/feedback", h.FeedbackHandler).Methods("POST")
			r.HandleFunc("/login", h.HandleLogin)
			r.HandleFunc("/githubcallback", h.HandleGithubCallback)
			r.HandleFunc("/apitoken", h.HandleApiToken)

			log.Infof("Default model provider: %s", h.DefaultProvider)
			log.Infof("Default model: %s", h.DefaultModel)

			if config.ServerConfig.TLSCertFile != "" && config.ServerConfig.TLSKeyFile != "" {

				log.Info("Server listening on TLS port 8443...")
				err = http.ListenAndServeTLS(":8443", config.ServerConfig.TLSCertFile, config.ServerConfig.TLSKeyFile, r)
			} else {
				log.Warn("No TLS configuration provided, traffic will not be encryption")
				log.Info("Server listening on port 8080...")
				err = http.ListenAndServe(":8080", r)
			}
			return err
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&o.configFile, "config", "c", "", "Config file to use")
	flags.StringVarP(&o.verbosity, "verbosity", "v", "info", "Log verbosity level (trace,debug,info,warn,error) (default info)")

	/*
		cmd.Flags().StringVarP(&o.email, "email", "e", "", "Model email address used when not provided in the request")
		cmd.Flags().StringVarP(&o.apiKey, "apikey", "a", "", "Model API key used when not provided in the request")
		cmd.Flags().StringVarP(&o.provider, "provider", "p", "ibm", "Which LLM provider to use when not provided in the request.  Value values are: ibm, openai")
		cmd.Flags().StringVarP(&o.model, "model", "m", "L3Byb2plY3RzL2czYmNfc3RhY2tfc3RnMl9lcG9jaDNfanVsXzMx", "Which LLM model to use from the chosen provider when not provided in the request.  Valid values depend on the chosen provider.")
	*/
	return cmd

}

func newInferCommand() *cobra.Command {
	o := inferOptions{}

	var cmd = &cobra.Command{
		Use:   "infer",
		Short: "Do a single inference",
		RunE: func(cmd *cobra.Command, args []string) error {
			level, err := log.ParseLevel(o.verbosity)
			if err != nil {
				log.WithError(err).Fatal("Cannot parse log-level")
			}
			log.SetLevel(level)

			if o.configFile == "" {
				return fmt.Errorf("config file is required")
			}
			config, err := loadConfig(o.configFile)
			if err != nil {
				return fmt.Errorf("error loading configfile %s: %v", o.configFile, err)
			}

			models = initModels(config)

			if o.prompt == "" {
				return fmt.Errorf("model prompt is required")
			}

			// If the user didn't specify a provider or model, use the defaults from the config file
			if o.provider == "" {
				o.provider = config.DefaultProvider
			}
			if o.modelId == "" {
				o.modelId = config.DefaultModelId
			}

			m, err := getModel(o.provider, o.modelId)
			if err != nil {
				return err
			}

			input := api.ModelInput{
				Prompt: o.prompt,
			}
			log.Debugf("Using provider/model %s/%s for prompt:\n%s\n", o.provider, o.modelId, o.prompt)
			response, err := model.InvokeModel(input, m)
			if err != nil {
				if response.Error != "" {
					log.Debugf("Response(Error):\n%s", response.Error)
				}
				return fmt.Errorf("error invoking the LLM: %v", err)
			}

			fmt.Printf("Response:\n%s\n", response.Output)

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&o.configFile, "config", "c", "", "Config file to use")
	flags.StringVarP(&o.prompt, "inference", "i", "", "Model prompt to be inferred")
	flags.StringVarP(&o.modelId, "model", "m", "", "Which LLM model to use from the provider.")
	flags.StringVarP(&o.provider, "provider", "p", "", "Which backend LLM provider to use.")
	flags.StringVarP(&o.verbosity, "verbosity", "v", "info", "Log verbosity level (trace,debug,info,warn,error) (default info)")

	return cmd

}

func initModels(config api.Config) map[string]api.Model {
	models := make(map[string]api.Model)
	for _, m := range config.Models {
		log.Debugf("Initializing model: %v", m)
		switch m.Provider {
		case "ibm":
			models[m.Provider+"/"+m.ModelId] = ibm.NewIBMModel(m.ModelId, m.URL, m.UserId, m.APIKey)
		case "openai":
			models[m.Provider+"/"+m.ModelId] = openai.NewOpenAIModel(m.ModelId, m.URL, m.APIKey)
		case "huggingface":
			models[m.Provider+"/"+m.ModelId] = hf.NewHFModel(m.ModelId, m.URL, m.APIKey)

		default:
			log.Errorf("unknown provider: %s", m.Provider)
		}
	}
	return models
}

func getModel(provider, modelId string) (api.Model, error) {
	if model, found := models[provider+"/"+modelId]; !found {
		return nil, fmt.Errorf("Provider/Model not found, valid provider/models: %q", reflect.ValueOf(models).MapKeys())
	} else {
		return model, nil
	}
}
