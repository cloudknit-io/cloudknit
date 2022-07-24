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
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// TeamValidatingWebhook implements admission.DecoderInjector.
// A decoder will be automatically injected.
type TeamValidatingWebhook struct {
	v  api.TeamValidator
	kc kClient.Client
	d  *admission.Decoder
	l  *logrus.Entry
	es eventservice.API
}

func NewTeamValidatingWebhook(v api.TeamValidator, es eventservice.API, log *logrus.Entry) *TeamValidatingWebhook {
	return &TeamValidatingWebhook{v: v, es: es, l: log}
}

func (s *TeamValidatingWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	s.l.Infof("Validating webhook received an admission request for Team object")
	s.l.Infof(string(req.Object.Raw))

	company := env.CompanyName()
	team := gjson.GetBytes(req.Object.Raw, "spec.teamName").String()
	event := &eventservice.Event{
		Scope:  string(eventservice.ScopeTeam),
		Object: req.Name,
		Meta: &eventservice.Meta{
			Company: company,
			Team:    team,
		},
		EventType: string(eventservice.TeamValidationError),
	}

	teamObject := v1.Team{}
	if err := s.d.Decode(req, &teamObject); err != nil {
		s.l.Errorf("error decoding raw object from admission request for company [%s] and team [%s]: %v", company, team, err)
		event.Payload = []string{err.Error()}
		if err := s.es.Record(ctx, event, s.l); err != nil {
			s.l.Errorf("error recording [%s] event for company [%s] and team [%s]: %v", event.EventType, company, team, err)
		}
		return admission.Errored(http.StatusBadRequest, err)
	}

	var err error
	switch {
	case req.Operation == admission2.Create:
		err = s.v.ValidateTeamCreate(ctx, &teamObject)
	case req.Operation == admission2.Update:
		err = s.v.ValidateTeamUpdate(ctx, &teamObject)
	}

	if err != nil {
		return admission.Denied(err.Error())
	}

	s.l.Infof("Validating Admission Request passed validation for Team object [%s]", req.Name)
	return admission.Allowed("")
}

// InjectDecoder injects the decoder.
func (s *TeamValidatingWebhook) InjectDecoder(d *admission.Decoder) error {
	s.d = d
	return nil
}
