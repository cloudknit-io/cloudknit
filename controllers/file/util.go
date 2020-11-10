package controllers

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
)

func SaveYamlFile(obj interface{}, fileName string) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	bytes, err2 := yaml.JSONToYAML(jsonBytes)
	if err2 != nil {
		panic(err2)
	}

	err3 := ioutil.WriteFile(fileName, bytes, 0644)
	if err3 != nil {
		panic(err3)
	}
}
