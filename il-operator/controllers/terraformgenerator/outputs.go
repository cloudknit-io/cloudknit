package terraformgenerator

import (
	"fmt"
	"strings"

	"github.com/go-errors/errors"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
)

func standardizeVariables(vars []*v1.Variable) ([]*Variable, error) {
	standardized := make([]*Variable, 0, len(vars))
	for _, v := range vars {
		var value string
		if v.ValueFrom != "" {
			tokens := strings.Split(v.ValueFrom, ".")
			if len(tokens) != 2 {
				return nil, errors.Errorf("invalid output referenced: %s. Output should be in format <component>.<name>", v.Name)
			}
			value = referenceOutputFromRemoteState(tokens[0], tokens[1])
		} else {
			value = v.Value
		}
		standardized = append(standardized, &Variable{Name: v.Name, Value: value})

	}
	return standardized, nil
}

func referenceOutputFromRemoteState(component string, output string) string {
	return fmt.Sprintf("data.terraform_remote_state.%s.outputs.%s", component, output)
}
