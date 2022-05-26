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
	"fmt"
	argocd2 "github.com/compuzest/zlifecycle-il-operator/controller/common/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controller/components/operations/argocd"
	"strings"
	"testing"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/sirupsen/logrus"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGenerateNewRbacConfigEmptyPolicyCsv(t *testing.T) {
	t.Parallel()

	log := logrus.NewEntry(logrus.New())
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

	log := logrus.NewEntry(logrus.New())
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

	log := logrus.NewEntry(logrus.New())
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
	log := logrus.NewEntry(logrus.New())
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

	mockAPI := argocd2.NewMockAPI(mockCtrl)

	repoOpts := argocd2.RepoOpts{RepoURL: "git@github.com:CompuZest/test_repo.git", SSHPrivateKey: "test_key"}

	mockToken := &argocd2.GetTokenResponse{Token: "test_token"}
	mockAPI.EXPECT().GetAuthToken().Return(&argocd2.GetTokenResponse{Token: "test_token"}, nil)
	token := fmt.Sprintf("Bearer %s", mockToken.Token)
	repo := argocd2.Repository{Repo: "git@github.com:CompuZest/test_repo2.git", Name: "test_repo2"}
	list := argocd2.RepositoryList{Items: []argocd2.Repository{repo}}
	mockAPI.EXPECT().ListRepositories(token).Return(&list, util.CreateMockResponse(200), nil)
	createRepoBody := argocd2.CreateRepoViaSSHBody{Name: "test_repo", Repo: "git@github.com:CompuZest/test_repo.git", SSHPrivateKey: "test_key"}
	mockAPI.EXPECT().CreateRepository(gomock.Eq(createRepoBody), token).Return(util.CreateMockResponse(200), nil)

	log := logrus.NewEntry(logrus.New())
	registered, err := argocd.RegisterRepo(log, mockAPI, &repoOpts)
	assert.True(t, registered)
	assert.NoError(t, err)
}

func TestRegisterRepoExistingRepo(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockArgocdAPI := argocd2.NewMockAPI(mockCtrl)

	mockToken := &argocd2.GetTokenResponse{Token: "test_token"}
	mockArgocdAPI.EXPECT().GetAuthToken().Return(mockToken, nil)
	repo := argocd2.Repository{Repo: "git@github.com:CompuZest/test_repo.git", Name: "test_repo"}
	list := argocd2.RepositoryList{Items: []argocd2.Repository{repo}}
	token := fmt.Sprintf("Bearer %s", mockToken.Token)
	mockArgocdAPI.EXPECT().ListRepositories(token).Return(&list, util.CreateMockResponse(200), nil)

	log := logrus.NewEntry(logrus.New())
	repoOpts := argocd2.RepoOpts{RepoURL: "git@github.com:CompuZest/test_repo.git", SSHPrivateKey: "test_key"}
	registered, err := argocd.RegisterRepo(log, mockArgocdAPI, &repoOpts)
	assert.False(t, registered)
	assert.NoError(t, err)
}
