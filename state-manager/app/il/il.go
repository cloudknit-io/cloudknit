package il

import (
	"fmt"

	"github.com/pkg/errors"
)

func BuildILComponentPath(meta *ComponentMeta, prefix string) (path string, err error) {
	if meta.Environment == "" || meta.Team == "" || meta.Component == "" {
		return "", errors.New("state must contain non-empty team, environment and component name")
	}

	path = fmt.Sprintf(
		"team/%s-team-environment/%s-environment-component/%s/terraform",
		meta.Team,
		meta.Environment,
		meta.Component,
	)
	if prefix != "" {
		path = fmt.Sprintf("%s/%s", prefix, path)
	}
	return path, nil
}
