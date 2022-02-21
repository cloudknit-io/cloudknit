package util

import (
	"context"

	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	perrors "github.com/pkg/errors"
	coreV1 "k8s.io/api/core/v1"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	AuthTier = string
	AuthMode = string
)

const (
	AuthTierCompany        AuthTier = "company"
	AuthTierInternal       AuthTier = "internal"
	AuthTierServiceAccount AuthTier = "serviceAccount"
	AuthModeGitHubApp      AuthMode = "githubApp"
	AuthModeSSH            AuthMode = "ssh"
)

func GetPrivateKey(ctx context.Context, c kClient.Client, key kClient.ObjectKey) ([]byte, error) {
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

func GetAuthTierPrivateKey(ctx context.Context, c kClient.Client, tier AuthTier) ([]byte, error) {
	switch tier {
	case AuthTierCompany:
		key := kClient.ObjectKey{Name: env.Config.GitHubAppSecretNameCompany, Namespace: env.Config.GitHubAppSecretNamespaceCompany}
		return GetPrivateKey(ctx, c, key)
	case AuthTierInternal:
		key := kClient.ObjectKey{Name: env.Config.GitHubAppSecretNameInternal, Namespace: env.Config.GitHubAppSecretNamespaceInternal}
		return GetPrivateKey(ctx, c, key)
	case AuthTierServiceAccount:
		key := kClient.ObjectKey{Name: env.Config.GitSSHSecretName, Namespace: env.SystemNamespace()}
		return GetPrivateKey(ctx, c, key)
	default:
		return nil, perrors.Errorf("invalid tier: %s", tier)
	}
}

func GetAuthTierModePrivateKey(ctx context.Context, c kClient.Client, mode AuthMode, tier AuthTier) ([]byte, error) {
	switch mode {
	case AuthModeGitHubApp:
		if tier != AuthTierCompany && tier != AuthTierInternal {
			return nil, perrors.Errorf("invalid auth tier [%s] for auth mode p%s]", tier, mode)
		}
		return GetAuthTierPrivateKey(ctx, c, tier)
	case AuthModeSSH:
		return GetAuthTierPrivateKey(ctx, c, AuthTierServiceAccount)
	default:
		return nil, perrors.Errorf("invalid auth mode [%s]", mode)
	}
}
