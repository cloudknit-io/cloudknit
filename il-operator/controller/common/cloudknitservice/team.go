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

func (s *Service) PostTeam(ctx context.Context, organizationName string,
	team stablev1.Team, log *logrus.Entry) error {
	endpoint := fmt.Sprintf("%s/%s/%s/%s", s.host, "v1/orgs", organizationName, "teams")

	log.
		Infof(
			"CloudKnitService Endpoint: %s",
			endpoint,
		)

	log.
		Infof(
			"Team Post Call via CloudKnitService for Team: %s",
			team.Spec.TeamName,
		)

	jsonBody, err := util.ToJSON(team.Spec)
	if err != nil {
		return errors.Wrap(err, "error marshaling Team Spec body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.Wrap(err, "error creating POST Team request")
	}
	req.Header.Add("Content-Type", runtime.ContentTypeJSON)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error executing POST Team request")
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return errors.Errorf("POST team returned a non-OK status code: [%d]", resp.StatusCode)
	}

	log.
		Infof(
			"Successful response from CloudKnitService for POST Team request for TeamName: %s",
			team.Spec.TeamName,
		)

	return nil
}

func (s *Service) GetTeam(ctx context.Context, organizationName string,
	teamName string, log *logrus.Entry) (*Team, error) {
	endpoint := fmt.Sprintf("%s/%s/%s/%s/%s", s.host, "v1/orgs", organizationName, "teams", teamName)

	log.
		Infof(
			"CloudKnitService Endpoint: %s",
			endpoint,
		)

	log.
		Infof(
			"Team Get Call via CloudKnitService for Team: %s",
			teamName,
		)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating GET Team request")
	}
	req.Header.Add("Content-Type", runtime.ContentTypeJSON)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error executing GET Team request")
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, errors.Errorf("POST team returned a non-OK status code: [%d]", resp.StatusCode)
	}

	log.
		Infof(
			"Successful response from CloudKnitService for GET Team request for TeamName: %s",
			teamName,
		)

	respBody, err := util.ReadBody(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading GET organization response body")
	}

	var r Team
	if err := util.FromJSON(&r, respBody); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling GET organization response body")
	}

	log.
		Infof(
			"Successful response from CloudKnitService for getting organization for OrgName: %s. respBody: %s",
			organizationName, respBody,
		)

	return &r, nil
}
