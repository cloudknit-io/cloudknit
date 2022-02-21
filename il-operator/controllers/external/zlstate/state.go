package zlstate

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_state_manager.go -package=zlstate "github.com/compuzest/zlifecycle-il-operator/controllers/external/zlstate" API
type API interface {
	Put(company string, environment *v1.Environment) error
}

type HTTPClient struct {
	ctx        context.Context
	log        *logrus.Entry
	host       string
	httpClient *http.Client
}

func NewHTTPStateManager(ctx context.Context, log *logrus.Entry) *HTTPClient {
	return &HTTPClient{
		ctx:        ctx,
		log:        log,
		host:       env.Config.ZLifecycleStateManagerURL,
		httpClient: util.GetHTTPClient(),
	}
}

func (s *HTTPClient) Put(company string, environment *v1.Environment) error {
	endpoint := fmt.Sprintf("%s/%s", s.host, "zl/state")

	zlstate := newZLState(company, environment)
	body := PutZLStateBody{
		Company:     company,
		Team:        environment.Spec.TeamName,
		Environment: environment.Spec.EnvName,
		ZLState:     zlstate,
	}

	s.log.WithField("state", util.ToJSONString(zlstate)).Info("Persisting zLstate via State Manager")

	jsonBody, err := util.ToJSON(body)
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
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 200 {
		return errors.Errorf("PUT zLstate returned a non-OK status code: [%d]", resp.StatusCode)
	}

	respBody, err := util.ReadBody(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading PUT zLstate response body")
	}

	r := PutZLStateResponse{}
	if err := util.FromJSON(&r, respBody); err != nil {
		return errors.Wrap(err, "error unmarshaling PUT zLstate response body")
	}

	s.log.WithField("message", r.Message).Info("Successful response from State Manager")

	return nil
}

func newZLState(company string, environment *v1.Environment) *ZLState {
	components := make([]*Component, 0, len(environment.Spec.Components))
	for _, ec := range environment.Spec.Components {
		component := &Component{
			Name:          ec.Name,
			Status:        "not_provisioned",
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
		components = append(components, component)
	}
	return &ZLState{
		Company:     company,
		Team:        environment.Spec.TeamName,
		Environment: environment.Spec.EnvName,
		Components:  components,
	}
}
