# zlifecycle-il-operator
zLifecycle Operator that generates Intermediate Language (ArgoCD Apps/Argo Workflows) from the CRD (Environment)

## Build & Push Docker image

Run following in the root directory.

```bash
make docker-build docker-push IMG=shahadarsh/zlifecycle-il-operator:latest
```

## Deploy controller to k8s cluster

Run following in the root directory.

```bash
export AWS_ACCOUNT_ID=xxxx
make docker-build docker-push IMG=$AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-il-operator:latest
```
