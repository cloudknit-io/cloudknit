package mutatingwebhook

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

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

func validateJSONSchema(input []byte, schema string) (*gojsonschema.Result, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(input)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, errors.Wrap(err, "error running json schema validator")
	}

	return result, nil
}
