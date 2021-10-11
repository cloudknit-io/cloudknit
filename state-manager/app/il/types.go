package il

type ComponentMeta struct {
	IL          string `json:"il"`
	Team        string `json:"team"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
}

type ZState struct {
	RepoURL string         `json:"repoUrl"`
	Meta    *ComponentMeta `json:"meta"`
}
