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

func FromYaml(in string, out interface{}) error {
	return y.Unmarshal([]byte(in), out)
}

func ToYaml(in interface{}) (ymlstring string, e error) {
	out, err := y.Marshal(in)
	if err != nil {
		return "", e
	}

	return string(out), nil
}

func ReverseParseGitURL(url string) (org, repo string, err error) {
	if url == "" {
		return "", "", errors.New("url param cannot be empty")
	}
	url = strings.TrimSuffix(url, ".git")
	repoParsed := false
	end := len(url) - 1
	isDelimiter := func(c int32) bool { return c == '/' || c == ':' }
	for i := range url {
		c := int32(url[end-i])
		if isDelimiter(c) {
			if !repoParsed {
				repoParsed = true
				continue
			}
			break
		}
		if !repoParsed {
			repo = string(c) + repo
			continue
		}
		org = string(c) + org
	}

	if org == "" || repo == "" {
		return "", "", errors.Errorf("invalid git url: %s", url)
	}

	return
}

func RewriteGitURLToHTTPS(repoURL string) string {
	switch {
	case strings.Contains(repoURL, githubSSHPrefix):
		return strings.ReplaceAll(repoURL, githubSSHPrefix, githubHTTPSPrefix)
	case strings.Contains(repoURL, gitlabSSHPrefix):
		return strings.ReplaceAll(repoURL, gitlabSSHPrefix, gitlabHTTPSPrefix)
	default:
		return repoURL
	}
}

func RewriteGitURLToSSH(repoURL string) string {
	switch {
	case strings.Contains(repoURL, githubHTTPSPrefix):
		return strings.ReplaceAll(repoURL, githubHTTPSPrefix, githubSSHPrefix)
	case strings.Contains(repoURL, gitlabHTTPSPrefix):
		return strings.ReplaceAll(repoURL, gitlabHTTPSPrefix, gitlabSSHPrefix)
	default:
		return repoURL
	}
}
