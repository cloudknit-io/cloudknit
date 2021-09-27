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

package argocd_test

import (
	"strings"
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestGenerateNewRbacConfigEmptyPolicyCsv(t *testing.T) {
	t.Parallel()

	log := ctrl.Log.WithName("TestGenerateNewRbacConfigEmptyPolicyCsv")
	policyCsv, err := argocd.GenerateNewRbacConfig(log, "", "test:payment", "payment", []string{"design"})
	assert.NoError(t, err)
	expectedPolicyCsv := `p,role:payment,repositories,get,*,allow
p,role:payment,applications,*,payment/*,allow
p,role:payment,applications,*,design/*,allow
g,test:payment,role:payment
`
	assert.ElementsMatch(t, strings.Split(expectedPolicyCsv, "\n"), strings.Split(policyCsv, "\n"))
}

func TestGenerateNewRbacConfigExistingPolicyCsv(t *testing.T) {
	t.Parallel()

	log := ctrl.Log.WithName("TestGenerateNewRbacConfigExistingPolicyCsv")
	existingPolicyCsv := `p,role:design,repositories,get,*,allow
p,role:design,applications,*,design/*,allow
g,test:design,role:design
`
	policyCsv, err := argocd.GenerateNewRbacConfig(log, existingPolicyCsv, "test:payment", "payment", []string{"design"})
	assert.NoError(t, err)
	expectedPolicyCsv := `p,role:design,repositories,get,*,allow
p,role:design,applications,*,design/*,allow
p,role:payment,repositories,get,*,allow
p,role:payment,applications,*,payment/*,allow
p,role:payment,applications,*,design/*,allow
g,test:design,role:design
g,test:payment,role:payment
`
	assert.ElementsMatch(t, strings.Split(expectedPolicyCsv, "\n"), strings.Split(policyCsv, "\n"))
}

func TestGenerateAdminRbacConfigEmptyPolicyCsv(t *testing.T) {
	t.Parallel()

	log := ctrl.Log.WithName("TestGenerateAdminRbacConfigEmptyPolicyCsv")
	policyCsv, err := argocd.GenerateAdminRbacConfig(log, "", "test:admin", "admin")
	assert.NoError(t, err)
	expectedPolicyCsv := `p,role:admin,certificates,*,*,allow
p,role:admin,applications,*,*/*,allow
p,role:admin,repositories,*,*,allow
p,role:admin,clusters,*,*,allow
p,role:admin,accounts,*,*,allow
p,role:admin,projects,*,*,allow
p,role:admin,gpgkeys,*,*,allow
g,test:admin,role:admin
`
	assert.ElementsMatch(t, strings.Split(expectedPolicyCsv, "\n"), strings.Split(policyCsv, "\n"))
}

func TestGenerateAdminRbacConfigExistingPolicyCsv(t *testing.T) {
	t.Parallel()

	existingPolicyCsv := `p,role:design,repositories,get,*,allow
p,role:design,applications,*,design/*,allow
g,zmart-tech-sandbox:design,role:design
`
	log := ctrl.Log.WithName("TestGenerateAdminRbacConfigEmptyPolicyCsv")
	policyCsv, err := argocd.GenerateAdminRbacConfig(log, existingPolicyCsv, "test:admin", "admin")
	assert.NoError(t, err)
	expectedPolicyCsv := `p,role:design,repositories,get,*,allow
p,role:design,applications,*,design/*,allow
p,role:admin,certificates,*,*,allow
p,role:admin,applications,*,*/*,allow
p,role:admin,repositories,*,*,allow
p,role:admin,clusters,*,*,allow
p,role:admin,accounts,*,*,allow
p,role:admin,projects,*,*,allow
p,role:admin,gpgkeys,*,*,allow
g,zmart-tech-sandbox:design,role:design
g,test:admin,role:admin
`
	assert.ElementsMatch(t, strings.Split(expectedPolicyCsv, "\n"), strings.Split(policyCsv, "\n"))
}

func TestRegisterRepoNewRepo(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockAPI := mocks.NewMockAPI(mockCtrl)

	repoOpts := argocd.RepoOpts{RepoURL: "git@github.com:CompuZest/test_repo.git", SSHPrivateKey: "test_key"}

	mockAPI.EXPECT().GetAuthToken().Return(&argocd.GetTokenResponse{Token: "test_token"}, nil)
	repo := argocd.Repository{Repo: "git@github.com:CompuZest/test_repo2.git", Name: "test_repo2"}
	list := argocd.RepositoryList{Items: []argocd.Repository{repo}}
	mockAPI.EXPECT().ListRepositories(gomock.Any()).Return(&list, common.CreateMockResponse(200), nil)
	createRepoBody := argocd.CreateRepoBody{Name: "test_repo", Repo: "git@github.com:CompuZest/test_repo.git", SSHPrivateKey: "test_key"}
	mockAPI.EXPECT().CreateRepository(createRepoBody, gomock.Any()).Return(common.CreateMockResponse(200), nil)

	log := ctrl.Log.WithName("TestRegisterRepoNewRepo")
	registered, err := argocd.RegisterRepo(log, mockAPI, repoOpts)
	assert.True(t, registered)
	assert.NoError(t, err)
}

func TestRegisterRepoExistingRepo(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockArgocdAPI := mocks.NewMockAPI(mockCtrl)

	mockArgocdAPI.EXPECT().GetAuthToken().Return(&argocd.GetTokenResponse{Token: "test_token"}, nil)
	repo := argocd.Repository{Repo: "git@github.com:CompuZest/test_repo.git", Name: "test_repo"}
	list := argocd.RepositoryList{Items: []argocd.Repository{repo}}
	mockArgocdAPI.EXPECT().ListRepositories(gomock.Any()).Return(&list, common.CreateMockResponse(200), nil)

	log := ctrl.Log.WithName("TestRegisterRepoExistingRepo")

	repoOpts := argocd.RepoOpts{RepoURL: "git@github.com:CompuZest/test_repo.git", SSHPrivateKey: "test_key"}
	registered, err := argocd.RegisterRepo(log, mockArgocdAPI, repoOpts)
	assert.False(t, registered)
	assert.NoError(t, err)
}
