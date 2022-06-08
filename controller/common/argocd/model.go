package argocd

import (
	"context"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/sirupsen/logrus"
)

type HTTPClient struct {
	ctx       context.Context
	serverURL string
	log       *logrus.Entry
}

type Credentials struct {
	Username string
	Password string
}

type GetTokenBody struct {
	Username string
	Password string
}

type GetTokenResponse struct {
	Token string
}

type RepoOpts struct {
	RepoURL                 string
	SSHPrivateKey           string
	Mode                    string
	GitHubAppPrivateKey     []byte
	GitHubAppInstallationID int64
	GitHubAppID             int64
}

type CreateRepoViaGitHubAppBody struct {
	Repo                    string `json:"repo"`
	Name                    string `json:"name"`
	GitHubAppPrivateKey     string `json:"githubAppPrivateKey"`
	GitHubAppInstallationID int64  `json:"githubAppInstallationID"`
	GitHubAppID             int64  `json:"githubAppID"`
}

type CreateRepoViaSSHBody struct {
	Repo          string `json:"repo"`
	Name          string `json:"name"`
	SSHPrivateKey string `json:"sshPrivateKey"`
}

type RegisterClusterBody struct {
	Name          string         `json:"name"`
	Config        *ClusterConfig `json:"config"`
	Namespaces    []string       `json:"namespaces"`
	Server        string         `json:"server"`
	ServerVersion string         `json:"serverVersion"`
}

type ClusterConfig struct {
	TLSClientConfig *TLSClientConfig `json:"tlsClientConfig"`
	BearerToken     string           `json:"bearerToken"`
}

type TLSClientConfig struct {
	CAData     string `json:"caData"`
	ServerName string `json:"serverName"`
}

type ListClustersResponse struct {
	Items []*struct {
		ClusterResources bool `json:"clusterResources"`
		Config           struct {
			AwsAuthConfig struct {
				ClusterName string `json:"clusterName"`
				RoleARN     string `json:"roleARN"`
			} `json:"awsAuthConfig"`
			BearerToken        string `json:"bearerToken"`
			ExecProviderConfig struct {
				APIVersion  string            `json:"apiVersion"`
				Args        []string          `json:"args"`
				Command     string            `json:"command"`
				Env         map[string]string `json:"env"`
				InstallHint string            `json:"installHint"`
			} `json:"execProviderConfig"`
			Password        string `json:"password"`
			TLSClientConfig struct {
				CaData     string `json:"caData"`
				CertData   string `json:"certData"`
				Insecure   bool   `json:"insecure"`
				KeyData    string `json:"keyData"`
				ServerName string `json:"serverName"`
			} `json:"tlsClientConfig"`
			Username string `json:"username"`
		} `json:"config"`
		ConnectionState struct {
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"connectionState"`
		Info struct {
			APIVersions     []string `json:"apiVersions"`
			ConnectionState struct {
				Message string `json:"message"`
				Status  string `json:"status"`
			} `json:"connectionState"`
			ServerVersion string `json:"serverVersion"`
		} `json:"info"`
		Name          string   `json:"name"`
		Namespaces    []string `json:"namespaces"`
		Server        string   `json:"server"`
		ServerVersion string   `json:"serverVersion"`
		Shard         string   `json:"shard"`
	} `json:"items"`
	Metadata struct {
		ResourceVersion string `json:"resourceVersion"`
		SelfLink        string `json:"selfLink"`
	} `json:"metadata"`
}

type CreateProjectBody struct {
	Project *appv1.AppProject `json:"project"`
}

type RepositoryList struct {
	Items []Repository `json:"items"`
}

type Repository struct {
	Repo string `json:"repo"`
	Name string `json:"name"`
}

type UpdateClusterBody map[string]interface{}
