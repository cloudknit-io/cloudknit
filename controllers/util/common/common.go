package common

import (
	"bytes"
	"encoding/json"
	"github.com/go-logr/logr"
	"io"
	"io/ioutil"
	"net/http"
)

func LogBody(log logr.Logger, body io.ReadCloser) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Error(err, "Error while deserializing body")
		return
	}
	bodyString := string(bodyBytes)
	log.Info(bodyString)
}

func ToJson(log logr.Logger, data interface{}) ([]byte, error) {
	jsoned, err := json.Marshal(data)
	if err != nil {
		log.Error(err, "Failed to marshal data to json")
		return nil, err
	}

	return jsoned, nil
}

func FromJson(log logr.Logger, s interface{}, jsonData []byte) error {
	err := json.Unmarshal(jsonData, s)
	if err != nil {
		log.Error(err, "Failed to unmarshal data from json")
		return err
	}

	return nil
}

func FromJsonMap(log logr.Logger, m map[string]interface{}, s interface{}) error {
	jsoned, err := ToJson(log, m)
	if err != nil {
		return err
	}
	err = FromJson(log, s, jsoned)
	if err != nil {
		return err
	}

	return nil
}

func ReadBody(log logr.Logger, stream io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Error(err, "Failed to read stream")
		return nil, err
	}

	return body, nil
}

func CreateMockResponse(code int) *http.Response {
	r := http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte{})), StatusCode: code}
	return &r
}
