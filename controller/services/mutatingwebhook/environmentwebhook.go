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
	environmentCRDName = "environments.stable.compuzest.com"
)

// EnvironmentMutatingWebhook implements admission.DecoderInjector.
// A decoder will be automatically injected.
type EnvironmentMutatingWebhook struct {
	kc kClient.Client
	d  *admission.Decoder
	l  *logrus.Entry
	es eventservice.API
}

func NewEnvironmentMutatingWebhook(kc kClient.Client, es eventservice.API, log *logrus.Entry) *EnvironmentMutatingWebhook {
	return &EnvironmentMutatingWebhook{kc: kc, es: es, l: log}
}

func (s *EnvironmentMutatingWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	s.l.Infof("Mutating webhook received an admission request for Environment object")
	s.l.Infof(string(req.Object.Raw))

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
		EventType: string(eventservice.EnvironmentSchemaValidationError),
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

	schema, err := s.getJSONSchema(ctx)
	if err != nil {
		s.l.Errorf(
			"error fetching JSON schema from Environment Custom Resource Definition for company [%s], team [%s] and environment [%s]: %v",
			company, team, environment, err,
		)
		event.Payload = []string{err.Error()}
		if err := s.es.Record(ctx, event, s.l); err != nil {
			s.l.Errorf("error recording [%s] event for company [%s], team [%s] and environment [%s]: %v", event.EventType, company, team, environment, err)
		}
		return admission.Errored(http.StatusInternalServerError, err)
	}

	s.l.WithField("schema", schema).Info("Fetched JSON schema for Environment CRD")
	result, err := validateJSONSchema(req.Object.Raw, schema)
	if err != nil {
		s.l.Errorf("error validating Environment Custom Resource Definition JSON schema for company [%s], team [%s] and environment [%s]: %v", company, team, environment, err)
		event.Payload = []string{err.Error()}
		if err := s.es.Record(ctx, event, s.l); err != nil {
			s.l.Errorf("error recording [%s] event for company [%s], team [%s] and environment [%s]: %v", event.EventType, company, team, environment, err)
		}
		return admission.Errored(http.StatusInternalServerError, err)
	}

	if !result.Valid() {
		logErrors(result.Errors(), s.l)
		event.Payload = stringifyValidationErrors(result.Errors())
		if err := s.es.Record(ctx, event, s.l); err != nil {
			s.l.Errorf("error recording [%s] event for company [%s], team [%s] and environment [%s]: %v", event.EventType, company, team, environment, err)
		}
		return admission.Errored(http.StatusBadRequest, buildValidationError(result.Errors()))
	}

	s.l.Infof("Mutating Admission Request passed validation for Environment object [%s]", req.Name)

	event.EventType = string(eventservice.EnvironmentSchemaValidationSuccess)
	if err := s.es.Record(ctx, event, s.l); err != nil {
		s.l.Errorf("error recording [%s] event for company [%s], team [%s] and environment [%s]: %v", event.EventType, company, team, environment, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, req.Object.Raw)
}

func (s *EnvironmentMutatingWebhook) getJSONSchema(ctx context.Context) (string, error) {
	crdList := kv1.CustomResourceDefinitionList{}
	if err := s.kc.List(ctx, &crdList); err != nil {
		return "", errors.Wrap(err, "error listing Environment Custom Resource Definition")
	}

	for _, crd := range crdList.Items {
		if crd.Name == environmentCRDName {
			return util.ToJSONString(crd.Spec.Versions[0].Schema.OpenAPIV3Schema), nil
		}
	}

	return "", errors.Errorf("could not find CRD with name: %s", environmentCRDName)
}

// InjectDecoder injects the decoder.
func (s *EnvironmentMutatingWebhook) InjectDecoder(d *admission.Decoder) error {
	s.d = d
	return nil
}
