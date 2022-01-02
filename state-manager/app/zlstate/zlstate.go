package zlstate

import "time"

type Backend interface {
	Get(key string) (*ZLState, error)
	Put(key string, state *ZLState, force bool) error
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
