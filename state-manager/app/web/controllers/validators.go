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

func validateGetZLStateRequest(req *GetZLStateRequest) error {
	if req.Company == "" {
		return errors.New(`request body is missing field "company"`)
	}
	if req.Team == "" {
		return errors.New(`request body is missing field "team"`)
	}
	if req.Environment == "" {
		return errors.New(`request body is missing field "environment"`)
	}

	return nil
}

func validatePostZLStateRequest(req *PostZLStateRequest) error {
	if req.Company == "" {
		return errors.New(`request body is missing field "company"`)
	}
	if req.Team == "" {
		return errors.New(`request body is missing field "team"`)
	}
	if req.Environment == "" {
		return errors.New(`request body is missing field "environment"`)
	}
	if req.ZLState == nil {
		return errors.New(`request body is missing field "zlstate"`)
	}

	return nil
}

func validatePatchZLStateRequest(req *PatchZLStateRequest) error {
	if req.Company == "" {
		return errors.New(`request body is missing field "company"`)
	}
	if req.Team == "" {
		return errors.New(`request body is missing field "team"`)
	}
	if req.Environment == "" {
		return errors.New(`request body is missing field "environment"`)
	}
	if req.Component == "" {
		return errors.New(`request body is missing field "component"`)
	}
	if req.Status == "" {
		return errors.New(`request body is missing field "status"`)
	}

	return nil
}
