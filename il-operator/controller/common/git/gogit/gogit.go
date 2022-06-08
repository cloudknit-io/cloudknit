package gogit

import (
	"context"

	"github.com/compuzest/zlifecycle-il-operator/controller/common/git"

	gogit "github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

type GoGit struct {
	ctx     context.Context
	r       *gogit.Repository
	options *Options
}

type Options struct {
	Token      string `json:"token"`
	PrivateKey []byte `json:"privateKey"`
	Mode       string `json:"mode"`
}

func NewGoGit(ctx context.Context, opts *Options) (git.API, error) {
	return &GoGit{options: opts, ctx: ctx}, nil
}

func (g *GoGit) getAuthOptions() (transport.AuthMethod, error) {
	if g.options.Mode == ModeToken {
		if g.options.Token == "" {
			return nil, errors.New("token is required if auth mode is token")
		}
		return getTokenAuth(g.options.Token), nil
	}
	if len(g.options.PrivateKey) == 0 {
		return nil, errors.New("private key is required if auth mode is ssh")
	}
	auth, err := getSSHAuth(g.options.PrivateKey)
	if err != nil {
		return nil, errors.New("error getting ssh auth")
	}
	return auth, nil
}

var _ git.API = (*GoGit)(nil)
