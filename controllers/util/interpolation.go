package util

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	substituteRegex = "(\\${[a-zA-Z1-9.]*})"
	paramRegex      = "[a-zA-Z1-9]*\\.[a-zA-Z1-9]*"
)

type Variables map[string]string

func Interpolate(template string, vars Variables) (string, error) {
	params := extractParams(template)

	for _, p := range params {
		pname, err := extractParameterName(p)
		if err != nil {
			return "", errors.Wrap(err, "error extracting parameter name")
		}
		pvalue := vars[pname]
		if pvalue == "" {
			continue
		}
		template = strings.ReplaceAll(template, p, pvalue)
	}

	return template, nil
}

func extractParameterName(p string) (string, error) {
	param := extractParamFromSubstitute(p)
	if len(param) != 1 {
		return "", errors.Errorf("invalid param name: parameter substitution format should be ${<scope>.<name>}: %s", p)
	}
	return param[0], nil
}

func extractParamFromSubstitute(p string) []string {
	r := regexp.MustCompile(paramRegex)
	return r.FindAllString(p, -1)
}

func extractParams(key string) []string {
	r := regexp.MustCompile(substituteRegex)
	return r.FindAllString(key, -1)
}
