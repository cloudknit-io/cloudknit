package state_manager

import v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_state_manager.go -package=state_manager "github.com/compuzest/zlifecycle-il-operator/controller/common/state_manager" API
type API interface {
	Get(company, team, environment string) (*GetZLStateResponse, error)
	Put(company, team string, environment *v1.Environment) error
	PutComponent(company, team, environment string, component *Component) error
}
