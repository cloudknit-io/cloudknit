package common

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func ToJson(data interface{}) ([]byte, error) {
	jsoned, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsoned, nil
}

func FromJson(s interface{}, jsonData []byte) error {
	err := json.Unmarshal(jsonData, s)
	if err != nil {
		return err
	}

	return nil
}

func FromJsonMap(m map[string]interface{}, s interface{}) error {
	jsoned, err := ToJson(m)
	if err != nil {
		return err
	}
	err = FromJson(s, jsoned)
	if err != nil {
		return err
	}

	return nil
}

func ReadBody(stream io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func CreateMockResponse(code int) *http.Response {
	r := http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte{})), StatusCode: code}
	return &r
}

// Helper functions to check and remove string from a slice of strings.
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func Stringify(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Find(s []string, e string) *string {
	for _, a := range s {
		if a == e {
			return &e
		}
	}
	return nil
}
