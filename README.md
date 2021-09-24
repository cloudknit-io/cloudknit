# zlifecycle-il-operator ![Build status badge](https://github.com/CompuZest/zlifecycle-il-operator/actions/workflows/main.yaml/badge.svg)

zLifecycle Operator that generates Intermediate Language (ArgoCD Apps/Argo Workflows) from the CRD (Environment)

## Prerequisites

### Linter
This project uses [golangci-lint](https://github.com/golangci/golangci-lint) for running code linting,
in conjunction with (File Watchers)[https://www.jetbrains.com/help/idea/using-file-watchers.html] plugin for IntelliJ
Lint rules are configured in the `.golangci.yaml` file in the project root.

### ArgoCD credentials
This project calls ArgoCD API endpoints, so it needs credentials from which it can request an Auth Token.
Credentials should be stored in a secret called `argocd-creds` in the operator namespace, usually `zlifecycle-il-operator-system`.
This secret also contains the ARGOCD_WEBHOOK_URL, and an optional ARGOCD_API_URL variable which should point argocd server,
the default value http://argocd-server.argocd.svc.cluster.local.
Check LastPass for secret values.
TODO: Refactor ARGOCD_WEBHOOK_URL and ARGOCD_API_URL to be a config variable instead of a secret value

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

## Troubleshooting
The `controller-gen not found` error can usually be fixed by updating your path to include your go packages `export PATH=$PATH:$(go env GOPATH)/bin`

## Auto-generated files
This operator combines a few things.
1. go code that we write in any directory
2. some go code/files generated by the `operator-sdk create` command into the `types/` and `controllers/` directories
3. the `role.yaml` file and CRD files for our custom resources generated by `controller-gen` command in the `make manifests` command
4. a file header generated by `controller-gen` command in the `make manifests` command
*Gotcha* Because operator-sdk is meant to work with kustomize it auto-generates an rbac resource without a prefix and will need manual attention.

`operator-sdk` is designed to work wtih `kustomize` (also confusing that it supports a helm operator to deploy dynamic helm charts, not the same as an operator deployed by a helm chart)
When creating new controllers, `operator-sdk` works well with `kustomize` and will generate templates into kustomize default directories such as `config/crd/bases` and `config/rbac`. For now, the `make manifests` command has been configured to generate files into `helm/templates` and future `operator-sdk` commands should also be run with a custom output folder to generate into that directory (can also be moved manually).

## Deploy controller to k8s cluster
helm install 

## Local Development
For faster docker builds, and the ability to shell into a contianer, use `Dockerfile.dev`, you can do this with `make docker-dev-build` or

Running `make manifests-local` uses the Mac friendly version of `sed` for the cleanup on the autogenerated file (see more details above in #auto-generated-files)

```
make docker-dev-build docker-push IMG=$AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-il-operator:branch-tag
```
## Tests
Mocks for test are auto-generated by running the `go generate ./...` command. This is built into `make docker-build` and should be run as any service that is mocked changes
