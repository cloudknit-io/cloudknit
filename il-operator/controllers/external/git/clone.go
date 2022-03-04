package git

import (
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pkg/errors"
)

const (
	username  = "git"
	ModeSSH   = "ssh"
	ModeToken = "token"
)

func (g *GoGit) Clone(repo string, directory string) error {
	auth, err := g.getAuthOptions()
	if err != nil {
		return errors.Wrap(err, "error getting auth options")
	}

	opts := gogit.CloneOptions{
		URL:  repo,
		Auth: auth,
	}

	r, err := gogit.PlainClone(directory, false, &opts)
	if err != nil {
		return err
	}
	g.r = r
	return nil
}

func getSSHAuth(privateKey []byte) (*ssh.PublicKeys, error) {
	pk, err := ssh.NewPublicKeys(username, privateKey, "")
	if err != nil {
		return nil, errors.Wrap(err, "error creating public keys from private key")
	}
	return pk, nil
}

func getTokenAuth(token string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: username,
		Password: token,
	}
}

func (g *GoGit) Open(path string) error {
	r, err := gogit.PlainOpen(path)
	if err != nil {
		return err
	}
	g.r = r
	return nil
}
