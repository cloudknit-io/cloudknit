package validator

import (
	"fmt"
	"regexp"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	git2 "github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/operations/git"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	errInitEnvironmentValidator = "error initializing environment validator"
	errInitTeamValidator        = "error initializing team validator"
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#rfc-1035-label-names
	// starts with alpha
	// ends with alphanumeric
	// cannot contain connective hyphens.
	nameRegex      = `^[a-zA-Z]+[a-zA-Z0-9]*(-[a-zA-Z0-9]+)*$`
	maxFieldLength = 63
)

var r = regexp.MustCompile(nameRegex)

func validateRFC1035String(str string) error {
	if !r.MatchString(str) {
		return errors.New("RFC1035 string must contain only lowercase alphanumeric characters " +
			"or '-', start with an alphabetic character, and end with an alphanumeric character")
	}
	return nil
}

func validateStringLength(str string) error {
	if len(str) > maxFieldLength {
		return errors.Errorf("string must not exceed %d characters in length", maxFieldLength)
	}
	return nil
}

func checkPaths(fs file.API, gitClient git2.API, source string, paths []string, fld *field.Path, l *logrus.Entry) field.ErrorList {
	var allErrs field.ErrorList

	dir, cleanup, err := git.CloneTemp(gitClient, source, l)
	if err != nil {
		l.Errorf("error temp cloning repo [%s]: %v", source, err)
		fe := field.InternalError(fld, errors.Errorf("error validating access to source repository [%s]", source))
		allErrs = append(allErrs, fe)
		return allErrs
	}

	for _, path := range paths {
		if exists, _ := fs.FileExistsInDir(dir, path); !exists {
			fe := field.Invalid(fld, path, fmt.Sprintf("file does not exist on given path in source repository [%s]", source))
			allErrs = append(allErrs, fe)
		}
	}
	defer cleanup()

	return allErrs
}
