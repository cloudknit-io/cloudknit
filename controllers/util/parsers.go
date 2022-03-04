package util

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	y "gopkg.in/yaml.v2"
)

const (
	githubHTTPSPrefix = "https://github.com/"
	githubSSHPrefix   = "git@github.com:"
	gitlabHTTPSPrefix = "https://gitlab.com/"
	gitlabSSHPrefix   = "git@gitlab.com:"
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
		return "", "", errors.New("url cannot be empty")
	}

	if prefix := "https://"; strings.HasPrefix(url, prefix) {
		trimmed := strings.TrimPrefix(strings.TrimSuffix(url, ".git"), prefix)
		splitted := strings.Split(trimmed, "/")
		owner = splitted[len(splitted)-2]
		repo = splitted[len(splitted)-1]
	}
	if prefix := "git@"; strings.HasPrefix(url, prefix) {
		trimmed := strings.TrimPrefix(strings.TrimSuffix(url, ".git"), prefix)
		splitted := strings.Split(trimmed, ":")
		if len(splitted) != 2 {
			return "", "", errors.New("invalid url format")
		}
		info := strings.Split(splitted[1], "/")
		if len(info) != 2 {
			return "", "", errors.New("invalid url format")
		}
		owner = info[0]
		repo = info[1]
	}
	return owner, repo, nil
}

func RewriteGitURLToHTTPS(repoURL string) string {
	transformed := repoURL
	if strings.HasPrefix(transformed, githubSSHPrefix) {
		transformed = strings.ReplaceAll(transformed, githubSSHPrefix, githubHTTPSPrefix)
	} else if strings.HasPrefix(transformed, gitlabSSHPrefix) {
		transformed = strings.ReplaceAll(transformed, gitlabSSHPrefix, gitlabHTTPSPrefix)
	}
	return transformed
}

func RewriteGitURLToSSH(repoURL string) string {
	transformed := repoURL
	if strings.HasPrefix(transformed, githubHTTPSPrefix) {
		transformed = strings.ReplaceAll(transformed, githubHTTPSPrefix, githubSSHPrefix)
	} else if strings.HasPrefix(transformed, gitlabHTTPSPrefix) {
		transformed = strings.ReplaceAll(transformed, gitlabHTTPSPrefix, gitlabSSHPrefix)
	}
	return transformed
}

func IsGitLabURL(repoURL string) bool {
	return strings.Contains(repoURL, "gitlab.com")
}
