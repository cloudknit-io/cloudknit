package main

import (
	"encoding/json"
	"fmt"
	argoWorkflow "github.com/compuzest/zlifecycle-il-operator/controllers/argoworkflow"
	"github.com/compuzest/zlifecycle-il-operator/mocks"
	"log"
)

func main() {
	env := mocks.GetMockEnv(false)

	reverse := argoWorkflow.GenerateLegacyWorkflowOfWorkflows(env)

	pretty, err := json.Marshal(reverse)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(string(pretty))
}
