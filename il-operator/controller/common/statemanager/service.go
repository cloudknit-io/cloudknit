package statemanager

import (
	"context"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_state_manager.go -package=statemanager "github.com/compuzest/zlifecycle-il-operator/controller/common/statemanager" API
type API interface {
	Get(ctx context.Context, company, team, environment string, log *logrus.Entry) (*GetZLStateResponse, error)
	Put(ctx context.Context, company, team string, environment *v1.Environment, log *logrus.Entry) error
	PutComponent(ctx context.Context, company, team, environment string, component *Component, log *logrus.Entry) error
}
