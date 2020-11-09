package main

import (
        "fmt"
        argocd "github.com/compuzest/environment-operator/controllers/argocd"
       	"k8s.io/apimachinery/pkg/util/json"
)

func main() {
        application := argocd.GenerateYaml()
        bytes, err := json.Marshal(application)
	if err != nil {
		panic(err)
	}
       fmt.Println(string(bytes))
}
