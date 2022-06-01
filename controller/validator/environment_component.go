package validator

import (
	"fmt"
	"strconv"
	"strings"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
)

func GetOutputFromComponent(name string, ec *v1.EnvironmentComponent) *v1.Output {
	for _, op := range ec.Outputs {
		if op.Name == name {
			return op
		}
	}

	return nil
}

func SplitValueFrom(vf string) (string, string, error) {
	output := strings.Split(vf, ".")

	if len(output) != 2 {
		return "", "", fmt.Errorf("valueFrom string is not properly formatted: %s", vf)
	}

	outputName := output[1]

	if left := strings.Index(outputName, "["); left >= 0 {
		right := strings.Index(outputName, "]")

		if _, err := strconv.ParseInt(outputName[left+1:right], 10, 64); err != nil {
			return "", "", fmt.Errorf("valueFrom index is not an integer: %s", outputName)
		}

		return output[0], outputName[:left], nil
	}

	// component name, output variable name, error
	return output[0], output[1], nil
}
