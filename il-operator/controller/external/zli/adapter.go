package zli

import "github.com/compuzest/zlifecycle-il-operator/controller/util"

type AuthMode string

const (
	authModeGitHubApp         = "github-app"
	authModeGitHubAppPublic   = "github-app-public"
	authModeGitHubAppInternal = "github-app-internal"
	authModeGitToken          = "token"
	authModeGitSSH            = "ssh"
)

func AuthModeToZLIAuthMode(authMode util.AuthMode, isPublic bool) AuthMode {
	switch authMode {
	case util.AuthModeSSH:
		return authModeGitSSH
	case util.AuthModeGitHubApp:
		if isPublic {
			return authModeGitHubAppPublic
		}
		return authModeGitHubAppInternal
	case util.AuthModeToken:
		return authModeGitToken
	default:
		return authModeGitSSH
	}
}
