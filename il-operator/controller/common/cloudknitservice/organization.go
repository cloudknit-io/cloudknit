package cloudknitservice

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (s *Service) GetOrganization(ctx context.Context, organizationName string, log *logrus.Entry) (*Organization, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", s.host, "v1/orgs", organizationName)

	log.
		Infof(
			"CloudKnitService Endpoint: %s",
			endpoint,
		)

	body := GetOrganizationBody{
		OrganizationName: organizationName,
	}

	log.
		Infof(
			"Fetching Organization via CloudKnitService for Org: %s",
			organizationName,
		)

	jsonBody, err := util.ToJSON(body)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling get organization body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating GET organization request")
	}
	req.Header.Add("Content-Type", runtime.ContentTypeJSON)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error executing GET organization request")
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("GET organization returned a non-OK status code: [%d]", resp.StatusCode)
	}

	respBody, err := util.ReadBody(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading GET organization response body")
	}

	var r Organization
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
