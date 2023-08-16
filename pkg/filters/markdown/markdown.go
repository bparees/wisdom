package markdown

import (
	"fmt"
	"regexp"

	//gomarkdown "github.com/gomarkdown/markdown"
	log "github.com/sirupsen/logrus"

	"github.com/openshift/wisdom/pkg/api"
)

var (
	markdownRegex = regexp.MustCompile("(?s)`{3}.*?\n(.*)`{3}")
)

func MarkdownStripper(response *api.ModelResponse) (*api.ModelResponse, error) {

	if response.Output == "" {
		return response, fmt.Errorf("response output is empty")
	}
	log.Debugf("Stripping markdown from response:\n %s\n", response.Output)

	//response.Output = markdownRegex.ReplaceAllString(response.Output, "")
	matches := markdownRegex.FindStringSubmatch(response.Output)
	response.Output = matches[1]
	/*
		node := gomarkdown.Parse([]byte(response.Output), nil)

		fmt.Printf("%#v\n", node.GetChildren()[0].GetChildren()[0])
		//fmt.Printf("%#v\n", node.GetChildren()[0].GetChildren()[1])
		//fmt.Printf("%s\n", node.GetChildren()[0].GetChildren()[1].AsLeaf().Literal)

		response.Output = string(node.GetChildren()[0].GetChildren()[0].AsLeaf().Literal)
	*/

	//response.Output = stripmd.Strip(response.Output)

	log.Debugf("Stripped markdown from response:\n %s\n", response.Output)
	return response, nil
}
