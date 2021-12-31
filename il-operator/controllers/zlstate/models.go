package zlstate

import "time"

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
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
