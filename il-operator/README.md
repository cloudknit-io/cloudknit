# zlifecycle-il-operator ![Build status badge](https://github.com/CompuZest/zlifecycle-il-operator/actions/workflows/main.yaml/badge.svg)

zLifecycle Operator that generates Intermediate Language (ArgoCD Apps/Argo Workflows) from the CRD (Environment)

## Prerequisites

### ArgoCD credentials
This project calls ArgoCD API endpoints, so it needs credentials from which it can request an Auth Token.
Credentials should be stored in a secret called `argocd-creds` in the operator namespace, usually `zlifecycle-il-operator-system`.
This secret also contains the ARGOCD_WEBHOOK_URL, and an optional ARGOCD_API_URL variable which should point argocd server,
the default value http://argocd-server.argocd.svc.cluster.local.
Check LastPass for secret values.
TODO: Refactor ARGOCD_WEBHOOK_URL and ARGOCD_API_URL to be a config variable instead of a secret value

### Auto-registration of team config repos
`Team` resource has a `repoSecret` field which is the name of the secret which holds the SSH key for the team config repo.
The secret should have the private SSH key stored in base64 format in the `sshPrivateKey` field.

## Vendoring

We are using `go mod vendor` for our code so that all dependencies are available to the operator without relying on external sources. 

Note: Any time go dependencies change remember to run `go mod vendor` at the repo root directory and commit the latest folder to source control.

## [Bootstrap zlifecycle-il-operator locally](./zlifecycle/runbook/setup/bootstrap-operator-locally.md)

## Build & Push Docker image

Run following in the root directory.

```bash
export ECR_REPO=[ THE AWS ECR REPO ]
make docker-push
```

## Deploy controller to k8s cluster

Run following in the root directory.

```bash
export ECR_REPO=[ THE AWS ECR REPO ]
make deploy
```

## Local Development
For faster docker builds, and the ability to shell into a contianer, use `Dockerfile.dev`, you can do this with `make docker-dev-build` or

```
make docker-dev-build docker-push IMG=$AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-il-operator:branch-tag
```
