package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/go-logr/logr"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TryRegisterRepo(
	c client.Client,
	log logr.Logger,
	ctx context.Context,
	api argocd.Api,
	repoUrl string,
	namespace string,
	repoSecret string,
) error {
	secret := &coreV1.Secret{}
	secretNamespacedName :=
		types.NamespacedName{Namespace: namespace, Name: repoSecret}
	if err := c.Get(ctx, secretNamespacedName, secret); err != nil {
		log.Info(
			"Secret %s does not exist in namespace %s\n",
			repoSecret,
			namespace,
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
		RepoUrl:       repoUrl,
		SshPrivateKey: sshPrivateKey,
	}

	if _, err := argocd.RegisterRepo(log, api, repoOpts); err != nil {
		return err
	}

	return nil
}