package util

import (
	"fmt"
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

	return output[0], output[1], nil
}
