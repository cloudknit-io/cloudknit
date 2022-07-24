package validating

import (
	"context"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/webhooks/api"
	admission2 "k8s.io/api/admission/v1"
	"net/http"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/eventservice"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// EnvironmentValidatingWebhook implements admission.DecoderInjector.
// A decoder will be automatically injected.
type EnvironmentValidatingWebhook struct {
	v  api.EnvironmentValidator
	d  *admission.Decoder
	l  *logrus.Entry
	es eventservice.API
}

func NewEnvironmentValidatingWebhook(v api.EnvironmentValidator, es eventservice.API, log *logrus.Entry) *EnvironmentValidatingWebhook {
	return &EnvironmentValidatingWebhook{v: v, es: es, l: log}
}

func (s *EnvironmentValidatingWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	s.l.Infof("Validating webhook received an admission request for Environment object")

	company := env.CompanyName()
	team := gjson.GetBytes(req.Object.Raw, "spec.teamName").String()
	environment := gjson.GetBytes(req.Object.Raw, "spec.envName").String()
	event := &eventservice.Event{
		Scope:  string(eventservice.ScopeEnvironment),
		Object: req.Name,
		Meta: &eventservice.Meta{
			Company:     company,
			Team:        team,
			Environment: environment,
		},
		EventType: string(eventservice.EnvironmentValidationError),
	}

	environmentObject := v1.Environment{}
	if err := s.d.Decode(req, &environmentObject); err != nil {
		s.l.Errorf("error decoding raw object from admission request for company [%s], team [%s] and environment [%s]: %v", company, team, environment, err)
		event.Payload = []string{err.Error()}
		if err := s.es.Record(ctx, event, s.l); err != nil {
			s.l.Errorf("error recording [%s] event for company [%s], team [%s] and environment [%s]: %v", event.EventType, company, team, environment, err)
		}
		return admission.Errored(http.StatusBadRequest, err)
	}

	var err error
	switch {
	case req.Operation == admission2.Create:
		err = s.v.ValidateEnvironmentCreate(ctx, &environmentObject)
	case req.Operation == admission2.Update:
		err = s.v.ValidateEnvironmentUpdate(ctx, &environmentObject)
	}

	if err != nil {
		return admission.Denied(err.Error())
	}

	s.l.Infof("Validating Admission Request passed validation for Environment object [%s]", req.Name)
	return admission.Allowed("")
}

// InjectDecoder injects the decoder.
func (s *EnvironmentValidatingWebhook) InjectDecoder(d *admission.Decoder) error {
	s.d = d
	return nil
}
