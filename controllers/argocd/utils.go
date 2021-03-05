package argocd

import (
	"encoding/json"
	"errors"
	"github.com/go-logr/logr"
	"io"
	"io/ioutil"
	"os"
)

func logBody(log logr.Logger, body io.ReadCloser) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Error(err, "Error while deserializing body")
		return
	}
	bodyString := string(bodyBytes)
	log.Info(bodyString)
}

func toJson(log logr.Logger, data interface{}) ([]byte, error) {
	jsoned, err := json.Marshal(data)
	if err != nil {
		log.Error(err, "Failed to marshal data to json")
		return nil, err
	}

	return jsoned, nil
}

func fromJson(log logr.Logger, d interface{}, jsonData []byte) error {
	err := json.Unmarshal(jsonData, d)
	if err != nil {
		log.Error(err, "Failed to unmarshal data from json")
		return err
	}

	return nil
}

func readBody(log logr.Logger, stream io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Error(err, "Failed to read stream")
		return nil, err
	}

	return body, nil
}

func getArgocdCredentialsFromEnv(log logr.Logger) (*ArgocdCredentials, error) {
	username := os.Getenv("argocd_username")
	password := os.Getenv("argocd_password")
	if username == "" || password == "" {
		return nil, errors.New("missing 'argocd_username' or 'argocd_password' env variables")
	}

	creds := ArgocdCredentials{Username: username, Password: password}

	return &creds, nil
}
