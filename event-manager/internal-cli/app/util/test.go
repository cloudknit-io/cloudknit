package util

import (
	"fmt"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/google/uuid"
)

func NewTestDirName() string {
	return fmt.Sprintf("%s-%s", env.TestDir, uuid.New().String())
}
