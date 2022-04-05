package zlstate

import (
	"time"
	// fix for gomock
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=./mock_backend.go -package=zlstate "github.com/compuzest/zlifecycle-state-manager/app/zlstate" Backend
type Backend interface {
	Get(key string) (*ZLState, error)
	Put(key string, state *ZLState, force bool) error
	UpsertComponent(key string, component *Component) (*ZLState, error)
}

type ZLState struct {
	Company     string       `json:"company"`
	Team        string       `json:"team"`
	Environment string       `json:"environment"`
	Components  []*Component `json:"components"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type Component struct {
	Name          string         `json:"name"`
	Type          string         `json:"type"`
	Status        string         `json:"status"`
	DependsOn     []string       `json:"dependsOn"`
	Module        *Module        `json:"module,omitempty"`
	Tags          []*Tags        `json:"tags,omitempty"`
	VariablesFile *VariablesFile `json:"variablesFile,omitempty"`
	OverlayFiles  []*OverlayFile `json:"overlayFiles,omitempty"`
	OverlayData   []*OverlayData `json:"overlayData,omitempty"`
	Variables     []*Variable    `json:"variables,omitempty"`
	Secrets       []*Secret      `json:"secrets,omitempty"`
	Outputs       []*Output      `json:"outputs,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

type Tags struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Module struct {
	Source  string `json:"source"`
	Path    string `json:"path,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type VariablesFile struct {
	Source string `json:"source"`
	Path   string `json:"path"`
	Ref    string `json:"ref,omitempty"`
}

type OverlayFile struct {
	Source string   `json:"source"`
	Ref    string   `json:"ref,omitempty"`
	Paths  []string `json:"paths"`
}

type OverlayData struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Variable struct {
	Name      string `json:"name"`
	Value     string `json:"value,omitempty"`
	ValueFrom string `json:"valueFrom,omitempty"`
}

type Secret struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
	Key   string `json:"key"`
}

type Output struct {
	Name      string `json:"name"`
	Sensitive bool   `json:"sensitive,omitempty"`
}
