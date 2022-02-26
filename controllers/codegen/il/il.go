package il

import (
	"context"

	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	git2 "github.com/compuzest/zlifecycle-il-operator/controllers/external/git"
)

type Service struct {
	TFILGitAPI   git2.API
	TFILTempDir  string
	TFILCleanupF git2.CleanupFunc
	ZLILGitAPI   git2.API
	ZLILTempDir  string
	ZLILCleanupF git2.CleanupFunc
}

func NewService(ctx context.Context, token string) (*Service, error) {
	zlILGitAPI, err := git2.NewGoGit(ctx, token)
	if err != nil {
		return nil, err
	}

	tfILGitAPI, err := git2.NewGoGit(ctx, token)
	if err != nil {
		return nil, err
	}

	// temp clone IL repo
	tempZLILRepoDir, zlILCleanup, err := git2.CloneTemp(zlILGitAPI, env.Config.ILZLifecycleRepositoryURL)
	if err != nil {
		return nil, err
	}

	tempTFILRepoDir, tfILCleanup, err := git2.CloneTemp(tfILGitAPI, env.Config.ILTerraformRepositoryURL)
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
