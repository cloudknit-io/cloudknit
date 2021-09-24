package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/go-logr/logr"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO: First check is the repo already registered
func TryRegisterRepo(ctx context.Context, c kClient.Client, log logr.Logger, api argocd.API, repoURL string, namespace string, repoSecret string) error {
	secret := &coreV1.Secret{}
	secretNamespacedName :=
		types.NamespacedName{Namespace: namespace, Name: repoSecret}
	if err := c.Get(ctx, secretNamespacedName, secret); err != nil {
		log.Info(
			"Secret does not exist in namespace\n",
			"secret", repoSecret,
			"namespace", namespace,
		)
		return err
	}

	sshPrivateKeyField := "sshPrivateKey"
	sshPrivateKey := string(secret.Data[sshPrivateKeyField])
	if sshPrivateKey == "" {
		errMsg := fmt.Sprintf("Secret is missing %s data field!", sshPrivateKeyField)
		err := errors.New(errMsg)
		log.Error(err, errMsg)
		return err
	}

	repoOpts := argocd.RepoOpts{
		RepoURL:       repoURL,
		SSHPrivateKey: sshPrivateKey,
	}

	if _, err := argocd.RegisterRepo(log, api, repoOpts); err != nil {
		return err
	}

	return nil
}
