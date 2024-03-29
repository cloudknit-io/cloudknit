package gitfactory

import (
	"context"

	git2 "github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/git/gogit"
	github2 "github.com/compuzest/zlifecycle-il-operator/controller/services/operations/github"

	"github.com/compuzest/zlifecycle-il-operator/controller/common/github"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Factory struct {
	client kClient.Client
	log    *logrus.Entry
}

func NewFactory(client kClient.Client, log *logrus.Entry) *Factory {
	return &Factory{
		client: client,
		log:    log,
	}
}

type Options struct {
	SSHOptions    *SSHOptions
	GitHubOptions *GitHubAppOptions
}

type SSHOptions struct {
	SecretName      string
	SecretNamespace string
}

type GitHubAppOptions struct {
	GitHubClient       github.API
	GitHubOrganization string
}

func (f *Factory) NewGitClient(ctx context.Context, opts *Options) (client git2.API, err error) {
	if opts == nil {
		return nil, errors.New("must provide valid options")
	}

	switch {
	case opts.SSHOptions != nil:
		if opts.SSHOptions == nil {
			return nil, errors.New("must provide valid ssh options")
		}
		client, err = f.newSSHGitClient(ctx, opts.SSHOptions)
		if err != nil {
			return nil, errors.Wrap(err, "error instantiating ssh gogit client")
		}
	case opts.GitHubOptions != nil:
		if opts.GitHubOptions == nil {
			return nil, errors.New("must provide valid GitHub App options")
		}
		client, err = f.newGitHubAppClient(ctx, opts.GitHubOptions)
		if err != nil {
			return nil, errors.Wrap(err, "error instantiating github app gogit client")
		}
	default:
		return nil, errors.Errorf("invalid options provided")
	}

	return client, nil
}

func (f *Factory) newSSHGitClient(ctx context.Context, opts *SSHOptions) (git2.API, error) {
	key := kClient.ObjectKey{Name: opts.SecretName, Namespace: opts.SecretNamespace}
	sshKey, err := util.GetPrivateKey(ctx, f.client, key)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching private key for ssh git client")
	}
	gitClient, err := gogit.NewGoGit(ctx, &gogit.Options{Mode: gogit.ModeSSH, PrivateKey: sshKey})
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating gogit client")
	}

	return gitClient, nil
}

func (f *Factory) newGitHubAppClient(ctx context.Context, opts *GitHubAppOptions) (git2.API, error) {
	token, err := github2.GenerateInstallationToken(f.log, opts.GitHubClient, opts.GitHubOrganization)
	if err != nil {
		return nil, errors.Wrap(err, "error generating installation token")
	}
	client, err := gogit.NewGoGit(ctx, &gogit.Options{Mode: gogit.ModeToken, Token: token})
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating git client")
	}

	return client, nil
}
