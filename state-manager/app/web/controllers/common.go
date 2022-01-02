package controllers

import (
	"fmt"
	"github.com/compuzest/zlifecycle-state-manager/app/env"
)

func BuildZLStateKey(team, environment string) string {
	return fmt.Sprintf("%s/%s.zlstate", team, environment)
}

func BuildZLStateBucketName(company string) string {
	return fmt.Sprintf("zlifecycle-%s-zlstate-%s", env.Config().Environment, company)
}
