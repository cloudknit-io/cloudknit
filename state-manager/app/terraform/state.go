package terraform

import tfjson "github.com/hashicorp/terraform-json"

func NewStateWrapper(state *tfjson.State) *StateWrapper {
	return &StateWrapper{state: state}
}

func (s *StateWrapper) ParseResources() []string {
	return extractResources(s.state.Values.RootModule)
}

func extractResources(m *tfjson.StateModule) []string {
	if m == nil {
		return []string{}
	}
	var resources []string
	for _, r := range m.Resources {
		resources = append(resources, r.Address)
	}
	for _, m := range m.ChildModules {
		resources = append(resources, extractResources(m)...)
	}
	return resources
}

func (s *StateWrapper) GetRawState() *tfjson.State {
	return s.state
}
