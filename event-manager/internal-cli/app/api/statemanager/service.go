package statemanager

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Manager interface {
	GetState(*FetchZLStateRequest) (*FetchZLStateResponse, error)
	GetComponent(request *FetchZLStateComponentRequest) (*FetchZLStateComponentResponse, error)
	PatchEnvironmentComponentStatus(request *UpdateZLStateComponentStatusRequest) (*UpdateZLStateComponentStatusResponse, error)
}

type HTTPStateManager struct {
	ctx        context.Context
	log        *logrus.Entry
	host       string
	httpClient *http.Client
}

func NewHTTPStateManager(ctx context.Context, log *logrus.Entry) *HTTPStateManager {
	return &HTTPStateManager{
		ctx:        ctx,
		log:        log,
		host:       env.StateManagerURL,
		httpClient: common.NewHTTPClient(),
	}
}

func (s *HTTPStateManager) GetState(request *FetchZLStateRequest) (*FetchZLStateResponse, error) {
	endpoint := "zl/state"
	url := fmt.Sprintf("%s/%s", s.host, endpoint)

	s.log.WithFields(logrus.Fields{
		"stateManagerURL": s.host,
		"endpoint":        endpoint,
		"company":         request.Company,
		"team":            request.Team,
		"environment":     request.Environment,
	}).Info("Fetching environment zLstate through zLifecycle State Manager")

	jsonBody, err := common.ToJSON(request)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling fetch zLstate request body")
	}

	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrapf(err, "error creating POST %s request", endpoint)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing POST %s request", endpoint)
	}
	defer common.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("POST %s returned a non-OK status code: [%d]", endpoint, resp.StatusCode)
	}

	respBody, err := common.ReadBody(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading POST %s response body", endpoint)
	}

	r := FetchZLStateResponse{}
	if err := common.FromJSON(&r, respBody); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling POST %s response body", endpoint)
	}

	s.log.WithFields(logrus.Fields{
		"method":     "POST",
		"statusCode": resp.StatusCode,
	}).Info("Successful response for environment zLstate from zLifecycle State Manager")

	return &r, nil
}

func (s *HTTPStateManager) GetComponent(request *FetchZLStateComponentRequest) (*FetchZLStateComponentResponse, error) {
	endpoint := "zl/state/component"
	url := fmt.Sprintf("%s/%s", s.host, endpoint)

	s.log.WithFields(logrus.Fields{
		"stateManagerURL": s.host,
		"endpoint":        endpoint,
		"company":         request.Company,
		"team":            request.Team,
		"environment":     request.Environment,
		"component":       request.Component,
	}).Info("Fetching environment component state from zLstate through zLifecycle State Manager")

	jsonBody, err := common.ToJSON(request)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling fetch zLstate component request body")
	}

	req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrapf(err, "error creating POST %s request", endpoint)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing POST %s request", endpoint)
	}
	defer common.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("POST %s returned a non-OK status code: [%d]", endpoint, resp.StatusCode)
	}

	respBody, err := common.ReadBody(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading POST %s response body", endpoint)
	}

	r := FetchZLStateComponentResponse{}
	if err := common.FromJSON(&r, respBody); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling POST %s response body", endpoint)
	}

	s.log.WithFields(logrus.Fields{
		"method":     "POST",
		"statusCode": resp.StatusCode,
	}).Info("Successful response for environment component state from zLifecycle State Manager")

	return &r, nil
}

func (s *HTTPStateManager) PatchEnvironmentComponentStatus(request *UpdateZLStateComponentStatusRequest) (*UpdateZLStateComponentStatusResponse, error) {
	endpoint := "zl/state/component"
	url := fmt.Sprintf("%s/%s", s.host, endpoint)

	s.log.WithFields(logrus.Fields{
		"stateManagerURL": s.host,
		"endpoint":        endpoint,
		"company":         request.Company,
		"team":            request.Team,
		"environment":     request.Environment,
		"component":       request.Component,
		"status":          request.Status,
	}).Info("Patching zLstate environment component status through zLifecycle State Manager")

	jsonBody, err := common.ToJSON(request)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling fetch zLstate component request body")
	}

	req, err := http.NewRequestWithContext(s.ctx, "PATCH", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrapf(err, "error creating PATCH %s request", endpoint)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing PATCH %s request", endpoint)
	}
	defer common.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("PATCH %s returned a non-OK status code: [%d]", endpoint, resp.StatusCode)
	}

	respBody, err := common.ReadBody(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading PATCH %s response body", endpoint)
	}

	r := UpdateZLStateComponentStatusResponse{}
	if err := common.FromJSON(&r, respBody); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling PATCH %s response body", endpoint)
	}

	s.log.WithFields(logrus.Fields{
		"method":     "PATCH",
		"statusCode": resp.StatusCode,
	}).Info("Successful response for environment component status update from zLifecycle State Manager")

	return &r, nil
}

var _ Manager = (*HTTPStateManager)(nil)
