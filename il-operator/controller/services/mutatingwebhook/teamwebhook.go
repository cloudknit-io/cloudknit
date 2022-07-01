package mutatingwebhook

import (
	"context"
	"net/http"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/eventservice"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	kv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	kClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	teamCRDName = "teams.stable.compuzest.com"
)

// TeamMutatingWebhook implements admission.DecoderInjector.
// A decoder will be automatically injected.
type TeamMutatingWebhook struct {
	kc kClient.Client
	d  *admission.Decoder
	l  *logrus.Entry
	es eventservice.API
}

func NewTeamMutatingWebhook(kc kClient.Client, es eventservice.API, log *logrus.Entry) *TeamMutatingWebhook {
	return &TeamMutatingWebhook{kc: kc, es: es, l: log}
}

func (s *TeamMutatingWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	s.l.Infof("Mutating webhook received an admission request for Team object")
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
		EventType: string(eventservice.TeamSchemaValidationError),
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

	schema, err := s.getJSONSchema(ctx)
	if err != nil {
		s.l.Errorf(
			"error fetching JSON schema from Team Custom Resource Definition for company [%s] and team [%s]: %v",
			company, team, err,
		)
		event.Payload = []string{err.Error()}
		if err := s.es.Record(ctx, event, s.l); err != nil {
			s.l.Errorf("error recording [%s] event for company [%s] and team [%s]: %v", event.EventType, company, team, err)
		}
		return admission.Errored(http.StatusInternalServerError, err)
	}

	s.l.WithField("schema", schema).Info("Fetched JSON schema for Team CRD")
	result, err := validateJSONSchema(req.Object.Raw, schema)
	if err != nil {
		s.l.Errorf("error validating Team Custom Resource Definition JSON schema for company [%s] and team [%s]: %v", company, team, err)
		event.Payload = []string{err.Error()}
		if err := s.es.Record(ctx, event, s.l); err != nil {
			s.l.Errorf("error recording [%s] event for company [%s] and team [%s]: %v", event.EventType, company, team, err)
		}
		return admission.Errored(http.StatusInternalServerError, err)
	}

	if !result.Valid() {
		logErrors(result.Errors(), s.l)
		event.Payload = stringifyValidationErrors(result.Errors())
		if err := s.es.Record(ctx, event, s.l); err != nil {
			s.l.Errorf("error recording [%s] event for company [%s] and team [%s]: %v", event.EventType, company, team, err)
		}
		return admission.Errored(http.StatusBadRequest, buildValidationError(result.Errors()))
	}

	s.l.Infof("Mutating Admission Request passed validation for Team object [%s]", req.Name)

	event.EventType = string(eventservice.TeamSchemaValidationSuccess)
	if err := s.es.Record(ctx, event, s.l); err != nil {
		s.l.Errorf("error recording [%s] event for company [%s] and team [%s]: %v", event.EventType, company, team, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, req.Object.Raw)
}

func (s *TeamMutatingWebhook) getJSONSchema(ctx context.Context) (string, error) {
	crdList := kv1.CustomResourceDefinitionList{}
	if err := s.kc.List(ctx, &crdList); err != nil {
		return "", errors.Wrap(err, "error listing Team Custom Resource Definition")
	}

	for _, crd := range crdList.Items {
		if crd.Name == teamCRDName {
			return util.ToJSONString(crd.Spec.Versions[0].Schema.OpenAPIV3Schema), nil
		}
	}

	return "", errors.Errorf("could not find CRD with name: %s", teamCRDName)
}

// InjectDecoder injects the decoder.
func (s *TeamMutatingWebhook) InjectDecoder(d *admission.Decoder) error {
	s.d = d
	return nil
}
