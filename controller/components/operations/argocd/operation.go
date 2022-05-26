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
	"fmt"
	argocdapi "github.com/compuzest/zlifecycle-il-operator/controller/common/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/aws/awseks"
	"strings"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/go-logr/logr"
)

func GenerateNewRbacConfig(log *logrus.Entry, oldPolicyCsv string, oidcGroup string, role string, additionalRoles []string) (newPolicyCsv string, err error) {
	rbacMap, err := parsePolicyCsv(oldPolicyCsv)
	if err != nil {
		return "", errors.Wrap(err, "error parsing policy CSV")
	}
	subject := fmt.Sprintf("role:%s", role)
	var projects []string
	projects = append(projects, role)
	projects = append(projects, additionalRoles...)
	log.WithFields(logrus.Fields{
		"subject":         subject,
		"additionalRoles": additionalRoles,
		"projects":        projects,
		"oidcGroup":       oidcGroup,
	}).Info(
		"Generating new RBAC configuration",
	)
	rbacMap.updateRbac(subject, projects, oidcGroup)

	return rbacMap.generatePolicyCsv(), nil
}

func GenerateAdminRbacConfig(log *logrus.Entry, oldPolicyCsv string, oidcGroup string, admin string) (newPolicyCsv string, err error) {
	rbacMap, err := parsePolicyCsv(oldPolicyCsv)
	if err != nil {
		return "", errors.Wrap(err, "error parsing policy CSV")
	}
	adminSubject := fmt.Sprintf("role:%s", admin)
	log.WithFields(logrus.Fields{
		"subject":   adminSubject,
		"oidcGroup": oidcGroup,
	}).Info("Generating admin RBAC configuration")

	return rbacMap.generateAdminRbac(adminSubject, oidcGroup), nil
}

func DeleteApplication(log *logrus.Entry, api argocdapi.API, app string) error {
	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		return err
	}
	bearer := toBearerToken(tokenResponse.Token)

	exists, err := api.DoesApplicationExist(app, bearer)
	if err != nil {
		return errors.Wrapf(err, "error checking does application [%s] exist", app)
	}
	if !exists {
		log.WithField("argocdApp", app).Info("Application does not exist, probably it has been already deleted")
		return nil
	}

	return api.DeleteApplication(app, bearer)
}

func RegisterRepo(log *logrus.Entry, api argocdapi.API, repoOpts *argocdapi.RepoOpts) (bool, error) {
	repoURI := repoOpts.RepoURL[strings.LastIndex(repoOpts.RepoURL, "/")+1:]
	repoName := strings.TrimSuffix(repoURI, ".git")
	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		return false, errors.Wrap(err, "error getting auth token")
	}
	bearer := toBearerToken(tokenResponse.Token)

	repositories, resp1, err := api.ListRepositories(bearer)
	if err != nil {
		return false, errors.Wrap(err, "error listing repositories")
	}
	defer util.CloseBody(resp1.Body)

	l := log.WithFields(logrus.Fields{
		"repoName": repoName,
		"repoUrl":  repoOpts.RepoURL,
	})
	if isRepoRegistered(repositories, repoOpts.RepoURL) {
		l.Infof("Repository %s already registered on ArgoCD", repoOpts.RepoURL)
		return false, nil
	}
	l.Infof("Repository %s is not registered on ArgoCD", repoOpts.RepoURL)

	var body interface{}
	if repoOpts.Mode == util.AuthModeGitHubApp {
		l.WithFields(logrus.Fields{
			"installationId": repoOpts.GitHubAppInstallationID,
			"appId":          repoOpts.GitHubAppID,
		}).Infof("Registering git repo %s in ArgoCD using GitHub App mode", repoOpts.RepoURL)
		body = argocdapi.CreateRepoViaGitHubAppBody{
			Repo:                    util.RewriteGitHubURLToHTTPS(repoOpts.RepoURL, false),
			Name:                    repoName,
			GitHubAppPrivateKey:     string(repoOpts.GitHubAppPrivateKey),
			GitHubAppInstallationID: repoOpts.GitHubAppInstallationID,
			GitHubAppID:             repoOpts.GitHubAppID,
		}
	} else {
		l.Infof("Registering git repo %s in ArgoCD using SSH mode", repoOpts.RepoURL)
		body = argocdapi.CreateRepoViaSSHBody{Repo: repoOpts.RepoURL, Name: repoName, SSHPrivateKey: repoOpts.SSHPrivateKey}
	}
	resp2, err := api.CreateRepository(body, bearer)
	if err != nil {
		return false, errors.Wrapf(err, "error registering repository [%s]", repoName)
	}
	defer util.CloseBody(resp2.Body)

	l.Infof("Successfully registered repository %s in ArgoCD", repoOpts.RepoURL)
	return true, nil
}

func isRepoRegistered(repos *argocdapi.RepositoryList, repoURL string) bool {
	for _, r := range repos.Items {
		if r.Repo == repoURL {
			return true
		}
	}
	return false
}

func TryCreateBootstrapApps(ctx context.Context, api argocdapi.API, log logr.Logger) error {
	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		return errors.Wrap(err, "error getting auth token")
	}
	bearer := toBearerToken(tokenResponse.Token)

	companyBootstrapApp := "company-bootstrap"
	exists, err := api.DoesApplicationExist(companyBootstrapApp, bearer)
	if err != nil {
		return errors.Wrapf(err, "error checking does application [%s] exists", companyBootstrapApp)
	}
	if exists {
		log.Info("Application already registered on ArgoCD",
			"application", companyBootstrapApp,
		)
	} else {
		companyResp, companyErr := api.CreateApplication(GenerateCompanyBootstrapApp(), bearer)
		if companyErr != nil {
			return errors.Wrap(companyErr, "error creating company bootstrap application")
		}
		defer util.CloseBody(companyResp.Body)
		log.Info("Successfully registered application on ArgoCD",
			"application", "company-bootstrap",
		)
	}

	configWatcherBootstrapApp := "config-watcher-bootstrap"
	exists2, err2 := api.DoesApplicationExist(configWatcherBootstrapApp, bearer)
	if err2 != nil {
		return errors.Wrapf(err2, "error checking does application [%s] exist", configWatcherBootstrapApp)
	}
	if exists2 {
		log.Info("Application already registered on ArgoCD",
			"application", configWatcherBootstrapApp,
		)
	} else {
		companyResp2, companyErr2 := api.CreateApplication(GenerateConfigWatcherBootstrapApp(), bearer)
		if companyErr2 != nil {
			return errors.Wrap(companyErr2, "error creating config watcher bootstrap app")
		}
		defer util.CloseBody(companyResp2.Body)
		log.Info("Successfully registered application on ArgoCD",
			"application", "config-watcher-bootstrap",
		)
	}

	return nil
}

func UpdateDefaultClusterNamespaces(log *logrus.Entry, api argocdapi.API, namespaces []string) error {
	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		return errors.Wrap(err, "error getting auth token")
	}
	bearer := toBearerToken(tokenResponse.Token)

	defaultClusterURL := "https://kubernetes.default.svc"
	body := make(argocdapi.UpdateClusterBody, 1)
	body["namespaces"] = namespaces
	log.WithField("namespaces", namespaces).Info("Updating default cluster namespaces")
	resp, err := api.UpdateCluster(defaultClusterURL, &body, []string{"namespaces"}, bearer)
	if err != nil {
		return errors.Wrap(err, "error updating default cluster")
	}
	defer util.CloseBody(resp.Body)

	return nil
}

func TryCreateProject(ctx context.Context, api argocdapi.API, log *logrus.Entry, name string, group string) (exists bool, err error) {
	l := log.WithField("project", name)

	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		return false, errors.Wrap(err, "error getting auth token")
	}
	bearer := toBearerToken(tokenResponse.Token)

	l.Info("Checking does AppProject already exist")
	exists, resp1, err1 := api.DoesProjectExist(name, bearer)
	if err1 != nil {
		return false, errors.Wrapf(err1, "error checking does argocd project [%s] exists", name)
	}
	defer util.CloseBody(resp1.Body)

	if exists {
		l.Info("ArgoCD AppProject already exist")
		return true, nil
	}

	body := argocdapi.CreateProjectBody{Project: toProject(name, group)}
	l.Info("ArgoCD AppProject does not exist, creating new one")
	resp2, err2 := api.CreateProject(&body, bearer)
	if err2 != nil {
		return false, errors.Wrap(err2, "error creating argocd project")
	}
	defer util.CloseBody(resp2.Body)

	return false, nil
}

func RegisterNewCluster(ctx context.Context, k8sClient awseks.API, argocdClient argocdapi.API, cluster string, log *logrus.Entry) (*awseks.ClusterInfo, error) {
	tokenResponse, err := argocdClient.GetAuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "error getting auth token")
	}
	bearer := toBearerToken(tokenResponse.Token)

	log.Infof("Describing cluster %s", cluster)
	info, err := k8sClient.DescribeCluster(ctx, cluster)
	if err != nil {
		return nil, errors.Wrapf(err, "error describing cluster %s", cluster)
	}

	log.Infof("Checking does k8s cluster %s exist", cluster)
	clusters, err := argocdClient.ListClusters(&cluster, bearer)
	if err != nil {
		return nil, errors.Wrapf(err, "error listing clusters")
	}

	for _, item := range clusters.Items {
		if item.Name == cluster {
			log.Infof("K8s cluster %s exist and will not register it", cluster)
			return info, nil
		}
	}

	log.Infof("K8s cluster %s exist and needs to be registered", cluster)

	log.Infof("Registering k8s cluster %s in ArgoCD", cluster)
	server := strings.ToLower(strings.TrimPrefix(info.Endpoint, "https://"))
	body := argocdapi.RegisterClusterBody{
		Name: info.Name,
		Config: &argocdapi.ClusterConfig{
			BearerToken: info.BearerToken,
			TLSClientConfig: &argocdapi.TLSClientConfig{
				CAData:     info.CertificateAuthority,
				ServerName: server,
			},
		},
		Namespaces:    []string{},
		Server:        info.Endpoint,
		ServerVersion: info.Version,
	}
	resp, err := argocdClient.RegisterCluster(&body, bearer)
	if err != nil {
		return nil, errors.Wrapf(err, "error registering cluster %s in argocd", cluster)
	}
	defer util.CloseBody(resp.Body)

	return info, nil
}
