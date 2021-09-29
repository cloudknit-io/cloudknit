package terraformgenerator

import (
	"fmt"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"strings"
)

func parseOutputs(outputs []*v1.Output) ([]string, error) {
	parsedOutputs := make([]string, 0, len(outputs))
	for _, o := range outputs {
		tokens := strings.Split(o.Name, ".")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("invalid output referenced: %s. Output should be in format <component>.<name>", o.Name)
		}
		parsed := fmt.Sprintf("data.terraform_remote_state.%s.outputs.%s", tokens[0], tokens[1])
		parsedOutputs = append(parsedOutputs, parsed)
	}
	return parsedOutputs, nil
}
