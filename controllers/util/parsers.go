package util

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	y "gopkg.in/yaml.v2"
)

const (
	httpsPrefix       = "https://github.com/"
	sshPrefix         = "git@github.com:"
	gitlabHttpsPrefix = "https://gitlab.com/"
	gitlabSshPrefix   = "git@gitlab.com:"
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

func ToJSONString(x interface{}) string {
	return string(ToJSONBytes(x, true))
}

func ToJSONCompact(x interface{}) string {
	return string(ToJSONBytes(x, false))
}

func ToJSONBytes(x interface{}, indent bool) []byte {
	if indent {
		b, _ := json.MarshalIndent(x, "", "  ")
		return b
	}
	b, _ := json.Marshal(x)
	return b
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

	httpsPrefix := "https://"
	if strings.HasPrefix(url, httpsPrefix) {
		trimmed := strings.TrimPrefix(strings.TrimSuffix(url, ".git"), httpsPrefix)
		splitted := strings.Split(trimmed, "/")
		owner = splitted[len(splitted)-2]
		repo = splitted[len(splitted)-1]
	} else {
		owner = url[strings.LastIndex(url, ":")+1 : strings.LastIndex(url, "/")]
		repoURI := url[strings.LastIndex(url, "/")+1:]
		repo = strings.TrimSuffix(repoURI, ".git")
	}
	return owner, repo, nil
}

func RewriteGitURLToHTTPS(repoURL string) string {
	transformed := repoURL
	if strings.HasPrefix(transformed, sshPrefix) {
		transformed = strings.ReplaceAll(transformed, sshPrefix, httpsPrefix)
	}
	if strings.HasPrefix(transformed, gitlabSshPrefix) {
		transformed = strings.ReplaceAll(transformed, gitlabSshPrefix, gitlabHttpsPrefix)
	}
	return transformed
}

func RewriteGitURLToSSH(repoURL string) string {
	transformed := repoURL
	if strings.HasPrefix(transformed, httpsPrefix) {
		transformed = strings.ReplaceAll(transformed, httpsPrefix, sshPrefix)
	}
	return transformed
}

func IsGitLabURL(repoURL string) bool {
	return strings.Contains(repoURL, "gitlab.com")
}
