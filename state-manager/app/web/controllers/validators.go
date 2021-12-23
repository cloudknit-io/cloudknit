package controllers

import (
	"github.com/compuzest/zlifecycle-state-manager/app/il"
	"github.com/pkg/errors"
)

func validateZState(state *il.ZState) error {
	if state == nil {
		return errors.New("invalid zstate")
	}
	if state.Meta == nil {
		return errors.New(`state is missing field "meta"`)
	}
	if state.RepoURL == "" {
		return errors.New(`state is missing field "repoUrl"`)
	}
	if state.Meta.IL == "" {
		return errors.New(`state.meta is missing field "il"`)
	}
	if state.Meta.Team == "" {
		return errors.New(`state.meta is missing field "team"`)
	}
	if state.Meta.Environment == "" {
		return errors.New(`state.meta is missing field "environment"`)
	}
	if state.Meta.Component == "" {
		return errors.New(`state.meta is missing field "component"`)
	}

	return nil
}
