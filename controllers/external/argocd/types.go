/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package argocd

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

//go:generate mockgen --build_flags=--mod=mod -destination=./mock_argocd_api.go -package=argocd "github.com/compuzest/zlifecycle-il-operator/controllers/external/argocd" API

type API interface {
	GetAuthToken() (*GetTokenResponse, error)
	ListRepositories(bearerToken string) (*RepositoryList, *http.Response, error)
	CreateRepository(body interface{}, bearerToken string) (*http.Response, error)
	CreateApplication(application *appv1.Application, bearerToken string) (*http.Response, error)
	DeleteApplication(name string, bearerToken string) error
	DoesApplicationExist(name string, bearerToken string) (exists bool, err error)
	CreateProject(project *CreateProjectBody, bearerToken string) (*http.Response, error)
	DoesProjectExist(name string, bearerToken string) (exists bool, response *http.Response, err error)
	ListClusters(name *string, bearerToken string) (*ListClustersResponse, error)
	RegisterCluster(body *RegisterClusterBody, bearerToken string) (*http.Response, error)
	UpdateCluster(clusterURL string, body *UpdateClusterBody, updatedFields []string, bearerToken string) (*http.Response, error)
}

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

/****************************/
/*           RBAC           */
/**************************.*.*/
type (
	EntryIdentifier = string
	Permission      = string
)

const (
	Policy EntryIdentifier = "p"
	Group  EntryIdentifier = "g"
)

const (
	Allow Permission = "allow"
	Deny  Permission = "deny"
)

type RbacPolicy struct {
	Identifier EntryIdentifier
	Subject    string
	Resource   string
	Action     string
	Object     string
	Permission Permission
}

type RbacGroup struct {
	Identifier EntryIdentifier
	Group      string
	Role       string
}

type RbacMap struct {
	Policies map[string][]*RbacPolicy
	Groups   map[string][]*RbacGroup
}
