package yaml

import (
	"fmt"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"

	"github.com/openshift/wisdom/pkg/api"
)

func YamlLinter(response *api.ModelResponse) (*api.ModelResponse, error) {
	if err := isValidYAML(response.Output); err != nil {
		return response, fmt.Errorf("response output is not valid YAML: %s", err)
	}
	return response, nil
}

func isValidYAML(yamlString string) error {
	var data interface{}

	log.Debugf("Validating YAML:\n%s", yamlString)
	err := yaml.Unmarshal([]byte(yamlString), &data)

	return err
}
