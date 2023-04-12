package cloudknitservice

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PostError struct {
	EnvName      string   `json:"envName"`
	ErrorMessage []string `json:"errorMessage"`
	Status       string   `json:"status"`
}

type PatchEnvironment struct {
	IsReconcile bool `json:"isReconcile"`
}

func (s *Service) PostEnvironment(ctx context.Context, organizationName string,
	environment stablev1.Environment, log *logrus.Entry) error {
	endpoint := fmt.Sprintf("%s/%s/%s/%s/%s/%s", s.host,
		"v1/orgs", organizationName, "teams",
		environment.Spec.TeamName, "environments")

	log.
		Infof(
			"CloudKnitService Endpoint: %s",
			endpoint,
		)

	log.
		Infof(
			"Environment Post Call via CloudKnitService for Env: %s",
			environment.Spec.EnvName,
		)

	jsonBody, err := util.ToJSON(environment.Spec)
	if err != nil {
		return errors.Wrap(err, "error marshaling Environment Spec body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.Wrap(err, "error creating POST Environment request")
	}
	req.Header.Add("Content-Type", runtime.ContentTypeJSON)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error executing POST Environment request")
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return errors.Errorf("POST Environment returned a non-OK status code: [%d]", resp.StatusCode)
	}

	log.
		Infof(
			"Successful response from CloudKnitService for POST Environment request for Environment: %s",
			environment.Spec.EnvName,
		)

	return nil
}

func (s *Service) PatchEnvironment(ctx context.Context, organizationName string,
	environment stablev1.Environment, log *logrus.Entry) error {
	endpoint := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s", s.host,
		"v1/orgs", organizationName, "teams",
		environment.Spec.TeamName, "environments", environment.Name)

	log.
		Infof(
			"CloudKnitService Endpoint: %s",
			endpoint,
		)

	log.
		Infof(
			"Environment PAtch Call via CloudKnitService for Env: %s",
			environment.Spec.EnvName,
		)

	jsonBody, err := util.ToJSON(PatchEnvironment{
		IsReconcile: true,
	})
	if err != nil {
		return errors.Wrap(err, "error marshaling Patch Environment body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.Wrap(err, "error creating PATCH Environment request")
	}
	req.Header.Add("Content-Type", runtime.ContentTypeJSON)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error executing PATCH Environment request")
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return errors.Errorf("PATCH Environment returned a non-OK status code: [%d]", resp.StatusCode)
	}

	log.
		Infof(
			"Successful response from CloudKnitService for PATCH Environment request for Environment: %s",
			environment.Spec.EnvName,
		)

	return nil
}

func (s *Service) PostError(ctx context.Context, organizationName string,
	environment *stablev1.Environment, allErrs []string, log *logrus.Entry) error {
	endpoint := fmt.Sprintf("%s/%s/%s/%s/%s/%s", s.host,
		"v1/orgs", organizationName, "teams",
		environment.Spec.TeamName, "errors")

	log.
		Infof(
			"CloudKnitService Endpoint: %s",
			endpoint,
		)

	log.
		Infof(
			"Error Post Call via CloudKnitService for Env: %s",
			environment.Spec.EnvName,
		)

	jsonBody, err := util.ToJSON(PostError{
		EnvName:      environment.Spec.EnvName,
		ErrorMessage: allErrs,
	})

	log.Infof(
		"Error Post Call via CloudKnitService jsonBody: %s",
		jsonBody,
	)

	if err != nil {
		return errors.Wrap(err, "error marshaling Errors body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.Wrap(err, "error creating POST Error request")
	}
	req.Header.Add("Content-Type", runtime.ContentTypeJSON)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error executing POST Error request")
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return errors.Errorf("POST Error returned a non-OK status code: [%d]", resp.StatusCode)
	}

	log.
		Infof(
			"Successful response from CloudKnitService for POST Errors request for Environment: %s",
			environment.Spec.EnvName,
		)

	return nil
}
