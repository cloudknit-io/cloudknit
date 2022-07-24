package api

import (
	"context"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
)

type EnvironmentValidator interface {
	ValidateEnvironmentCreate(context.Context, *v1.Environment) error
	ValidateEnvironmentUpdate(context.Context, *v1.Environment) error
}

type TeamValidator interface {
	ValidateTeamCreate(context.Context, *v1.Team) error
	ValidateTeamUpdate(context.Context, *v1.Team) error
}
