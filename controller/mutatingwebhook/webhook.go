package mutatingwebhook

import (
	"context"
	"fmt"
	"net/http"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/eventservice"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
	kv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	kClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	environmentCRDName = "environments.stable.compuzest.com"
)

// Service implements admission.DecoderInjector.
// A decoder will be automatically injected.
type Service struct {
	kc kClient.Client
	d  *admission.Decoder
	l  *logrus.Entry
	es eventservice.API
}

func NewService(kc kClient.Client, es eventservice.API, log *logrus.Entry) *Service {
	return &Service{kc: kc, es: es, l: log}
}

func (s *Service) Handle(ctx context.Context, req admission.Request) admission.Response {
	s.l.Infof("Mutating webhook received an admission request for Environment object")
	s.l.Infof(string(req.Object.Raw))

	company := env.CompanyName()
	team := gjson.GetBytes(req.Object.Raw, "spec.teamName").String()
	environment := gjson.GetBytes(req.Object.Raw, "spec.envName").String()
	event := &eventservice.Event{
		Object:      req.Name,
		Company:     company,
		Team:        team,
		Environment: environment,
		EventType:   string(eventservice.SchemaValidationError),
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

	event.EventType = string(eventservice.SchemaValidationSuccess)
	if err := s.es.Record(ctx, event, s.l); err != nil {
		s.l.Errorf("error recording [%s] event for company [%s], team [%s] and environment [%s]: %v", event.EventType, company, team, environment, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, req.Object.Raw)
}

func logErrors(verrors []gojsonschema.ResultError, log *logrus.Entry) {
	log.Error("error validating Environment CRD JSON schema")
	for _, verr := range verrors {
		// Err implements the ResultError interface
		log.Errorf("- %s", verr)
	}
}

func stringifyValidationErrors(validationErrors []gojsonschema.ResultError) []string {
	stringErrors := make([]string, 0, len(validationErrors))
	for _, e := range validationErrors {
		stringErrors = append(stringErrors, e.String())
	}
	return stringErrors
}

func buildValidationError(verrors []gojsonschema.ResultError) error {
	var msg string
	for i, verr := range verrors {
		if i == 0 {
			msg += verr.String()
		} else {
			msg += fmt.Sprintf(": %s", verr.String())
		}
	}
	return errors.New(msg)
}

func (s *Service) getJSONSchema(ctx context.Context) (string, error) {
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

func validateJSONSchema(input []byte, schema string) (*gojsonschema.Result, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(input)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, errors.Wrap(err, "error running json schema validator")
	}

	return result, nil
}

// InjectDecoder injects the decoder.
func (s *Service) InjectDecoder(d *admission.Decoder) error {
	s.d = d
	return nil
}
