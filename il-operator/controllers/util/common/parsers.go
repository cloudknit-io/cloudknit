package common

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	y "gopkg.in/yaml.v2"
)

func ToJSON(data interface{}) ([]byte, error) {
	jsoned, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsoned, nil
}

func FromJSON(s interface{}, jsonData []byte) error {
	err := json.Unmarshal(jsonData, s)
	if err != nil {
		return err
	}

	return nil
}

func FromJSONMap(m map[string]interface{}, s interface{}) error {
	jsoned, err := ToJSON(m)
	if err != nil {
		return err
	}
	err = FromJSON(s, jsoned)
	if err != nil {
		return err
	}

	return nil
}

func FromYaml(yamlstring string, out interface{}) error {
	return y.Unmarshal([]byte(yamlstring), out)
}

func ToYaml(in interface{}) (ymlstring string, e error) {
	out, err := y.Marshal(in)
	if err != nil {
		return "", e
	}

	return string(out), nil
}

func ParseRepositoryName(url string) string {
	repoURI := url[strings.LastIndex(url, "/")+1:]
	return strings.TrimSuffix(repoURI, ".git")
}

func ParseRepositoryInfo(url string) (owner string, repo string, err error) {
	if url == "" {
		return "", "", errors.New("URL cannot be empty")
	}

	owner = url[strings.LastIndex(url, ":")+1 : strings.LastIndex(url, "/")]
	repoURI := url[strings.LastIndex(url, "/")+1:]
	repo = strings.TrimSuffix(repoURI, ".git")

	return owner, repo, nil
}

func RewriteGitURLToHTTPS(repoURL string) string {
	httpsRepo := repoURL
	httpsPrefix := "https://github.com/"
	if sshPrefix := "git@github.com:"; strings.HasPrefix(httpsRepo, sshPrefix) {
		httpsRepo = strings.ReplaceAll(httpsRepo, sshPrefix, httpsPrefix)
	}
	return httpsRepo
}
