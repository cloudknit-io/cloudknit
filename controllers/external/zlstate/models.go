package zlstate

import (
	"time"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
)

type PutZLStateBody struct {
	Company     string   `json:"company"`
	Team        string   `json:"team"`
	Environment string   `json:"environment"`
	ZLState     *ZLState `json:"zlstate"`
}

type PutZLStateResponse struct {
	Message string `json:"message"`
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
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	Status        string            `json:"status"`
	DependsOn     []string          `json:"dependsOn"`
	Module        *v1.Module        `json:"module,omitempty"`
	Tags          []*v1.Tags        `json:"tags,omitempty"`
	VariablesFile *v1.VariablesFile `json:"variablesFile,omitempty"`
	OverlayFiles  []*v1.OverlayFile `json:"overlayFiles,omitempty"`
	OverlayData   []*v1.OverlayData `json:"overlayData,omitempty"`
	Variables     []*v1.Variable    `json:"variables,omitempty"`
	Secrets       []*v1.Secret      `json:"secrets,omitempty"`
	Outputs       []*v1.Output      `json:"outputs,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}
