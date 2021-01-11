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
make deploy IMG=shahadarsh/zlifecycle-il-operator:latest
```
