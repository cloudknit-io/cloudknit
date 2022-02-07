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

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/go-logr/logr"
	_ "github.com/golang/mock/mockgen/model" // workaround for mockgen failing
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=../../mocks/mock_argocd_api.go -package=mocks "github.com/compuzest/zlifecycle-il-operator/controllers/argocd" API

const (
	ModeGitHubApp = "githubApp"
	ModeSSH       = "ssh"
)

type API interface {
	GetAuthToken() (*GetTokenResponse, error)
	ListRepositories(bearerToken string) (*RepositoryList, *http.Response, error)
	CreateRepository(body interface{}, bearerToken string) (*http.Response, error)
	CreateApplication(application *appv1.Application, bearerToken string) (*http.Response, error)
	DeleteApplication(name string, bearerToken string) error
	DoesApplicationExist(name string, bearerToken string) (exists bool, err error)
	CreateProject(project *CreateProjectBody, bearerToken string) (*http.Response, error)
	DoesProjectExist(name string, bearerToken string) (exists bool, response *http.Response, err error)
	UpdateCluster(clusterURL string, body *UpdateClusterBody, updatedFields []string, bearerToken string) (*http.Response, error)
}

type HTTPAPI struct {
	ctx       context.Context
	serverURL string
	log       logr.Logger
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
	GitHubAppInstallationID string
	GitHubAppID             string
}

type CreateRepoViaGitHubAppBody struct {
	Repo                    string `json:"repo"`
	Name                    string `json:"name"`
	GitHubAppPrivateKey     string `json:"githubAppPrivateKey"`
	GitHubAppInstallationID string `json:"githubAppInstallationID"`
	GitHubAppID             string `json:"githubAppID"`
}

type CreateRepoViaSSHBody struct {
	Repo          string `json:"repo"`
	Name          string `json:"name"`
	SSHPrivateKey string `json:"sshPrivateKey"`
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
