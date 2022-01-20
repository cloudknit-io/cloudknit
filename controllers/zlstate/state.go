package zlstate

import (
	"bytes"
	"context"
	"fmt"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=../../mocks/mock_state_manager.go -package=mocks "github.com/compuzest/zlifecycle-il-operator/controllers/zlstate" StateManager
type StateManager interface {
	Put(company string, team string, environment string, state *ZLState) error
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
		host:       env.Config.ZLifecycleStateManagerURL,
		httpClient: common.GetHTTPClient(),
	}
}

func (s *HTTPStateManager) Put(company string, environment *v1.Environment) error {
	endpoint := fmt.Sprintf("%s/%s", s.host, "zl/state")

	zlstate := newZLState(company, environment)
	body := PutZLStateBody{
		Company:     company,
		Team:        environment.Spec.TeamName,
		Environment: environment.Spec.EnvName,
		ZLState:     zlstate,
	}

	s.log.WithField("state", zlstate).Info("Persisting zLstate through State Manager")

	jsonBody, err := common.ToJSON(body)
	if err != nil {
		return errors.Wrap(err, "error marshaling put zLstate body")
	}

	req, err := http.NewRequestWithContext(s.ctx, "PUT", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.Wrap(err, "error creating PUT zLstate request")
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error executing PUT zLstate request")
	}
	defer common.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return errors.Errorf("PUT zLstate returned a non-OK status code: [%d]", resp.StatusCode)
	}

	respBody, err := common.ReadBody(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading PUT zLstate response body")
	}

	r := PutZLStateResponse{}
	if err := common.FromJSON(&r, respBody); err != nil {
		return errors.Wrap(err, "error unmarshaling PUT zLstate response body")
	}

	s.log.WithField("message", r.Message).Info("Successful response from State Manager")

	return nil
}

func newZLState(company string, environment *v1.Environment) *ZLState {
	components := make([]*Component, 0, len(environment.Spec.Components))
	for _, ec := range environment.Spec.Components {
		component := &Component{
			Name:   ec.Name,
			Status: "not_provisioned",
		}
		components = append(components, component)
	}
	return &ZLState{
		Company:     company,
		Team:        environment.Spec.TeamName,
		Environment: environment.Spec.EnvName,
		Components:  components,
	}
}
