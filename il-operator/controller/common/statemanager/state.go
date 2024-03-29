package statemanager

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	statusNotProvisioned = "not_provisioned"
)

type Service struct {
	host       string
	httpClient *http.Client
}

func NewService(host string) *Service {
	return &Service{
		host:       host,
		httpClient: util.GetHTTPClient(),
	}
}

func (s *Service) Put(ctx context.Context, company, team string, environment *v1.Environment, log *logrus.Entry) error {
	endpoint := fmt.Sprintf("%s/%s", s.host, "zl/state")

	zlstate := newZLState(company, environment)
	body := PutZLStateBody{
		Company:     company,
		Team:        team,
		Environment: environment.Spec.EnvName,
		ZLState:     zlstate,
	}

	log.
		WithField("state", util.ToJSONString(zlstate)).
		Infof(
			"Persisting zLstate for company [%s], team [%s] and environment [%s] via State Manager",
			company, environment.Spec.TeamName, environment.Spec.EnvName,
		)

	jsonBody, err := util.ToJSON(body)
	if err != nil {
		return errors.Wrap(err, "error marshaling put zLstate body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.Wrap(err, "error creating PUT zLstate request")
	}
	req.Header.Add("Content-Type", runtime.ContentTypeJSON)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error executing PUT zLstate request")
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return errors.Errorf("PUT zLstate returned a non-OK status code: [%d]", resp.StatusCode)
	}

	respBody, err := util.ReadBody(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading PUT zLstate response body")
	}

	var r PutZLStateResponse
	if err := util.FromJSON(&r, respBody); err != nil {
		return errors.Wrap(err, "error unmarshaling PUT zLstate response body")
	}

	log.
		WithField("message", r.Message).
		Infof(
			"Successful response from State Manager for adding zL state for company [%s], team [%s] and environment [%s]",
			company, team, environment.Spec.EnvName,
		)

	return nil
}

func newZLState(company string, environment *v1.Environment) *ZLState {
	components := make([]*Component, 0, len(environment.Spec.Components))
	for _, ec := range environment.Spec.Components {
		components = append(components, ToZLStateComponent(ec))
	}
	return &ZLState{
		Company:     company,
		Team:        environment.Spec.TeamName,
		Environment: environment.Spec.EnvName,
		Components:  components,
	}
}

func ToZLStateComponent(ec *v1.EnvironmentComponent) *Component {
	return &Component{
		Name:          ec.Name,
		Status:        statusNotProvisioned,
		Type:          ec.Type,
		DependsOn:     ec.DependsOn,
		Module:        ec.Module,
		Tags:          ec.Tags,
		VariablesFile: ec.VariablesFile,
		OverlayFiles:  ec.OverlayFiles,
		OverlayData:   ec.OverlayData,
		Variables:     ec.Variables,
		Secrets:       ec.Secrets,
		Outputs:       ec.Outputs,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}
}

func (s *Service) Get(ctx context.Context, company, team, environment string, log *logrus.Entry) (*GetZLStateResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.host, "zl/state")

	body := GetZLStateBody{
		Company:     company,
		Team:        team,
		Environment: environment,
	}

	log.
		Infof(
			"Fetching zLstate for company [%s], team [%s] and environment [%s] via State Manager",
			company, team, environment,
		)

	jsonBody, err := util.ToJSON(body)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling get zLstate body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "error creating GET zLstate request")
	}
	req.Header.Add("Content-Type", runtime.ContentTypeJSON)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error executing GET zLstate request")
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("GET zLstate returned a non-OK status code: [%d]", resp.StatusCode)
	}

	respBody, err := util.ReadBody(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading GET zLstate response body")
	}

	var r GetZLStateResponse
	if err := util.FromJSON(&r, respBody); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling GET zLstate response body")
	}

	log.
		Infof(
			"Successful response from State Manager for getting zL state for company [%s], team [%s] and environment [%s]",
			company, team, environment,
		)

	return &r, nil
}

func (s *Service) PutComponent(ctx context.Context, company, team, environment string, component *Component, log *logrus.Entry) error {
	endpoint := fmt.Sprintf("%s/%s", s.host, "zl/state/component")

	body := PutZLStateComponentBody{
		Company:     company,
		Team:        team,
		Environment: environment,
		Component:   component,
	}

	log.
		Infof(
			"Adding zLstate component [%s] for company [%s], team [%s] and environment [%s] via State Manager",
			component.Name, company, team, environment,
		)

	jsonBody, err := util.ToJSON(body)
	if err != nil {
		return errors.Wrap(err, "error marshaling put zLstate component body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.Wrap(err, "error creating PUT zLstate component request")
	}
	req.Header.Add("Content-Type", runtime.ContentTypeJSON)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error executing PUT zLstate component request")
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return errors.Errorf("PUT zLstate component returned a non-OK status code: [%d]", resp.StatusCode)
	}

	log.
		Infof(
			"Successful response from State Manager for adding new component to zL state for company [%s], team [%s] and environment [%s]",
			company, team, environment,
		)

	return nil
}
