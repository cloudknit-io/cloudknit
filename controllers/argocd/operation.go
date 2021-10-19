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
	"strings"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
)

func GenerateNewRbacConfig(log logr.Logger, oldPolicyCsv string, oidcGroup string, role string, additionalRoles []string) (newPolicyCsv string, err error) {
	rbacMap, err := parsePolicyCsv(oldPolicyCsv)
	if err != nil {
		return "", err
	}
	subject := fmt.Sprintf("role:%s", role)
	var projects []string
	projects = append(projects, role)
	projects = append(projects, additionalRoles...)
	log.Info(
		"Generating new RBAC configuration",
		"subject", subject,
		"additionalRoles", additionalRoles,
		"projects", projects,
		"oidcGroup", oidcGroup,
	)
	rbacMap.updateRbac(subject, projects, oidcGroup)

	return rbacMap.generatePolicyCsv(), nil
}

func GenerateAdminRbacConfig(log logr.Logger, oldPolicyCsv string, oidcGroup string, admin string) (newPolicyCsv string, err error) {
	rbacMap, err := parsePolicyCsv(oldPolicyCsv)
	if err != nil {
		return "", err
	}
	adminSubject := fmt.Sprintf("role:%s", admin)
	log.Info(
		"Generating admin RBAC configuration",
		"subject", adminSubject,
		"oidcGroup", oidcGroup,
	)

	return rbacMap.generateAdminRbac(adminSubject, oidcGroup), nil
}

func DeleteApplication(log logr.Logger, api API, name string) error {
	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		return err
	}
	bearer := toBearerToken(tokenResponse.Token)

	exists, err := api.DoesApplicationExist(name, bearer)
	if err != nil {
		return err
	}
	if !exists {
		log.Info(
			"Application does not exist, probably it has been already deleted",
			"application", name,
		)
		return nil
	}

	return api.DeleteApplication(name, bearer)
}

func RegisterRepo(log logr.Logger, api API, repoOpts RepoOpts) (bool, error) {
	repoURI := repoOpts.RepoURL[strings.LastIndex(repoOpts.RepoURL, "/")+1:]
	repoName := strings.TrimSuffix(repoURI, ".git")
	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		return false, err
	}
	bearer := toBearerToken(tokenResponse.Token)
	repositories, resp1, err1 := api.ListRepositories(bearer)
	if err1 != nil {
		return false, err1
	}
	defer common.CloseBody(resp1.Body)
	if isRepoRegistered(*repositories, repoOpts.RepoURL) {
		log.Info("Repository already registered on ArgoCD",
			"repoName", repoName,
			"repoUrl", repoOpts.RepoURL,
		)
		return false, nil
	}
	log.Info(
		"Repository is not registered on ArgoCD, registering now...",
		"repoName", repoName,
		"repoUrl", repoOpts.RepoURL,
	)
	createRepoBody := CreateRepoBody{Repo: repoOpts.RepoURL, Name: repoName, SSHPrivateKey: repoOpts.SSHPrivateKey}
	resp2, err2 := api.CreateRepository(createRepoBody, bearer)
	if err2 != nil {
		return false, err2
	}
	defer common.CloseBody(resp2.Body)
	log.Info("Successfully registered repository on ArgoCD", "repo", repoOpts.RepoURL)
	return true, nil
}

func TryCreateBootstrapApps(ctx context.Context, log logr.Logger) error {
	api := NewHTTPClient(ctx, log, env.Config.ArgocdServerURL)

	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		return err
	}
	bearer := toBearerToken(tokenResponse.Token)

	exists, err := api.DoesApplicationExist("company-bootstrap", bearer)
	if err != nil {
		return fmt.Errorf("error while calling ArgoCD API to check if Application Exists: %w", err)
	}
	if exists {
		log.Info("Application already registered on ArgoCD",
			"application", "company-bootstrap",
		)
	} else {
		companyResp, companyErr := api.CreateApplication(GenerateCompanyBootstrapApp(), bearer)
		if companyErr != nil {
			return fmt.Errorf("error while creating Company Bootstrap Application: %w", companyErr)
		}
		defer common.CloseBody(companyResp.Body)
		log.Info("Successfully registered application on ArgoCD",
			"application", "company-bootstrap",
		)
	}

	exists2, err2 := api.DoesApplicationExist("config-watcher-bootstrap", bearer)
	if err2 != nil {
		return err2
	}
	if exists2 {
		log.Info("Application already registered on ArgoCD",
			"application", "config-watcher-bootstrap",
		)
	} else {
		companyResp2, companyErr2 := api.CreateApplication(GenerateConfigWatcherBootstrapApp(), bearer)
		if companyErr2 != nil {
			return fmt.Errorf("error while creating Config Watcher Bootstrap Application: %w", companyErr2)
		}
		defer common.CloseBody(companyResp2.Body)
		log.Info("Successfully registered application on ArgoCD",
			"application", "config-watcher-bootstrap",
		)
	}

	return nil
}

func TryCreateProject(ctx context.Context, log logr.Logger, name string, group string) (exists bool, err error) {
	api := NewHTTPClient(ctx, log, env.Config.ArgocdServerURL)

	tokenResponse, err := api.GetAuthToken()
	if err != nil {
		return false, err
	}
	bearer := toBearerToken(tokenResponse.Token)

	log.Info("Checking does AppProject already exist", "name", name)
	exists, resp1, err := api.DoesProjectExist(name, bearer)
	if err != nil {
		return false, err
	}
	defer common.CloseBody(resp1.Body)

	if exists {
		log.Info("AppProject already exist", "name", name)
		return true, nil
	}

	body := CreateProjectBody{Project: toProject(name, group)}
	log.Info("Checking does AppProject does not exist, creating new one", "name", name)
	resp2, err := api.CreateProject(body, bearer)
	if err != nil {
		return false, err
	}
	defer common.CloseBody(resp2.Body)

	return false, nil
}

func toProject(name string, group string) *appv1.AppProject {
	typeMeta := metav1.TypeMeta{APIVersion: "argoproj.io/v1alpha1", Kind: "AppProject"}
	objectMeta := metav1.ObjectMeta{Name: name, Namespace: "argocd"}
	spec := appv1.AppProjectSpec{
		SourceRepos:              []string{"*"},
		Destinations:             []appv1.ApplicationDestination{{Server: "https://kubernetes.default.svc", Namespace: "*"}},
		ClusterResourceWhitelist: []metav1.GroupKind{{Group: "*", Kind: "*"}},
		Roles: []appv1.ProjectRole{
			{
				Name:   "frontend",
				Groups: []string{fmt.Sprintf("%s:%s", group, name)},
				Policies: []string{
					fmt.Sprintf("p, proj:%s:frontend, applications, get, %s/*, allow", name, name),
					fmt.Sprintf("p, proj:%s:frontend, applications, delete, %s/*, allow", name, name),
					fmt.Sprintf("p, proj:%s:frontend, applications, sync, %s/*, allow", name, name),
				},
			},
		},
	}
	return &appv1.AppProject{TypeMeta: typeMeta, ObjectMeta: objectMeta, Spec: spec, Status: appv1.AppProjectStatus{}}
}

func toBearerToken(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}
