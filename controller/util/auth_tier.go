package util

type (
	AuthTier = string
	AuthMode = string
)

const (
	AuthTierCompany        AuthTier = "company"
	AuthTierInternal       AuthTier = "internal"
	AuthTierServiceAccount AuthTier = "serviceAccount"
	AuthModeGitHubApp      AuthMode = "githubApp"
	AuthModeSSH            AuthMode = "ssh"
	AuthModeToken          AuthMode = "token"
)
