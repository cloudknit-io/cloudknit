package validator

import (
	"strconv"
	"strings"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/pkg/errors"
)

func GetOutputFromComponent(name string, ec *v1.EnvironmentComponent) *v1.Output {
	for _, op := range ec.Outputs {
		if op.Name == name {
			return op
		}
	}

	return nil
}

func SplitValueFrom(vf string) (compName string, outputVarName string, err error) {
	output := strings.Split(vf, ".")

	if len(output) != 2 {
		return compName, outputVarName, errors.Errorf("valueFrom string is not properly formatted: %s", vf)
	}

	compName = output[0]
	outputVarName = output[1]

	if left := strings.Index(outputVarName, "["); left >= 0 {
		right := strings.Index(outputVarName, "]")

		if right == -1 {
			return "", "", errors.Errorf("valueFrom is not formatted properly: %s", outputVarName)
		}

		if _, err := strconv.ParseInt(outputVarName[left+1:right], 10, 64); err != nil {
			return "", "", errors.Errorf("valueFrom index is not an integer: %s", outputVarName)
		}

		return compName, outputVarName[:left], nil
	}

	// component name, output variable name, error
	return compName, outputVarName, nil
}
