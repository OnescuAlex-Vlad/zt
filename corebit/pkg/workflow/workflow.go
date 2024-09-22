package workflow

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Step struct {
	Name    string                 `yaml:"name"`
	Action  string                 `yaml:"action"`
	Input   map[string]interface{} `yaml:"input"`
	Output  map[string]interface{} `yaml:"output"`
}

type Workflow struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

func RunWorkflow(filePath string) error {
	var workflow Workflow
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, &workflow)
	if err != nil {
		return err
	}
	fmt.Printf("Executing workflow: %s\n", workflow.Name)
	for _, step := range workflow.Steps {
		fmt.Printf("Executing step: %s\n", step.Name)
		// Placeholder: Get and execute action by name
	}
	return nil
}
