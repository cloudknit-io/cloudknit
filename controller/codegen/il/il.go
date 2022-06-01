package il

import (
	"context"

	git2 "github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/git/gogit"
	"github.com/compuzest/zlifecycle-il-operator/controller/components/operations/git"

	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controller/env"
)

type Service struct {
	TFILGitAPI   git2.API
	TFILTempDir  string
	TFILCleanupF git.CleanupFunc
	ZLILGitAPI   git2.API
	ZLILTempDir  string
	ZLILCleanupF git.CleanupFunc
}

func NewService(ctx context.Context, token string, log *logrus.Entry) (*Service, error) {
	zlILGitAPI, err := gogit.NewGoGit(ctx, &gogit.Options{Mode: gogit.ModeToken, Token: token})
	if err != nil {
		return nil, err
	}

	tfILGitAPI, err := gogit.NewGoGit(ctx, &gogit.Options{Mode: gogit.ModeToken, Token: token})
	if err != nil {
		return nil, err
	}

	// temp clone IL repo
	tempZLILRepoDir, zlILCleanup, err := git.CloneTemp(zlILGitAPI, env.Config.ILZLifecycleRepositoryURL, log)
	if err != nil {
		return nil, err
	}

	tempTFILRepoDir, tfILCleanup, err := git.CloneTemp(tfILGitAPI, env.Config.ILTerraformRepositoryURL, log)
	if err != nil {
		zlILCleanup()
		return nil, err
	}

	return &Service{
		TFILGitAPI:   tfILGitAPI,
		TFILTempDir:  tempTFILRepoDir,
		TFILCleanupF: tfILCleanup,
		ZLILGitAPI:   zlILGitAPI,
		ZLILTempDir:  tempZLILRepoDir,
		ZLILCleanupF: zlILCleanup,
	}, nil
}

type config struct {
	TeamDirectory          string
	CompanyDirectory       string
	ConfigWatcherDirectory string
}

var Config = config{
	TeamDirectory:          env.Config.ILTeamFolder,
	CompanyDirectory:       env.Config.ILCompanyFolder,
	ConfigWatcherDirectory: env.Config.ILConfigWatcherFolder,
}
