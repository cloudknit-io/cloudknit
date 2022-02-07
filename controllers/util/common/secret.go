package common

import (
	"context"

	perrors "github.com/pkg/errors"
	coreV1 "k8s.io/api/core/v1"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

func GetSSHPrivateKey(ctx context.Context, c kClient.Client, key kClient.ObjectKey) ([]byte, error) {
	secret := coreV1.Secret{}
	if err := c.Get(ctx, key, &secret); err != nil {
		return nil, perrors.Wrapf(err, "error getting secret %s/%s from k8s cache", key.Namespace, key.Name)
	}

	sshPrivateKeyField := "sshPrivateKey"
	sshPrivateKey := secret.Data[sshPrivateKeyField]
	if len(sshPrivateKey) == 0 {
		return nil, perrors.Errorf(`secret %s/%s is invalid: missing field "sshPrivateKey"`, key.Namespace, key.Name)
	}

	return sshPrivateKey, nil
}
