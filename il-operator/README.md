# zlifecycle-il-operator
zLifecycle Operator that generates Intermediate Language (ArgoCD Apps/Argo Workflows) from the CRD (Environment)

## Vendoring

We are using `go mod vendor` for our code so that all dependencies are available to the operator without relying on external sources. 

Note: Any time go dependencies change remember to run `go mod vendor` at the repo root directory and commit the latest folder to source control.

## [Bootstrap zlifecycle-il-operator locally](./zlifecycle/runbook/setup/bootstrap-operator-locally.md)

## Build & Push Docker image

Run following in the root directory.

```bash
export AWS_ACCOUNT_ID=xxxx
make docker-build docker-push IMG=$AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-il-operator:latest
```

## Deploy controller to k8s cluster

Run following in the root directory.

```bash
export AWS_ACCOUNT_ID=xxxx
make deploy IMG=$AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-il-operator:latest
```

## Local Development
The operator image does not come with a shell, to debug the container change the image in `Dockerfile` to:  `gcr.io/distroless/base:debug`