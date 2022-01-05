package zlstate

import "time"

type FetchZLStateRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
}

type FetchZLStateResponse struct {
	ZLState *ZLState `json:"zlstate"`
}

type FetchZLStateComponentRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
}

type FetchZLStateComponentResponse struct {
	Component *Component `json:"component"`
}

type UpdateZLStateComponentStatusRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
	Status      string `json:"status"`
}

type UpdateZLStateComponentStatusResponse struct {
	Message string `json:"message"`
}

type UpdateZLStateComponentRequest struct {
	Company     string `json:"company"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
	Status      string `json:"status"`
}

type UpdateZLStateComponentResponse struct{}

type ZLState struct {
	Company     string       `json:"company"`
	Team        string       `json:"team"`
	Environment string       `json:"environment"`
	Components  []*Component `json:"components"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type Component struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
