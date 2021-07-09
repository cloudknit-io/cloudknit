# Destroy an environment 

## Overview

Destroy is currently semi-automated and has a few manual steps that needs to be followed

## When to use this runbook

When you want to destroy an environment

## Initial Steps Overview

- [Delete Environment Config](#delete-environment-config)
- [Delete IL code](#delete-il-code)
- [Delete Orphan Resources from k8s](#delete-orphan-resources-from-k8s)

## Detailed Steps

### Delete Environment Config
1. Go to Team Config Repo like `https://github.com/zmart-tech/zmart-payment-team-config`
2. Delete the Environment folder
3. Push the changes to Github
4. Manually Sync the Environment (this doesn't trigger an automated sync currently)
5. Verify and Approve the destroys on the UI

### Delete IL code
1. Delete both the Environment yaml and folder from the IL repo
2. Push the changes to Github

### Delete Orphan Resources from k8s
1. Run following queries to delete argo workflows and argocd apps.
[Note] Make sure to change the Env Name instead of `payment-prod`

```bash
kubectl get workflow -n argocd | awk '/payment-prod/{print $1}'| xargs  kubectl delete -n argocd workflow
```

```bash
kubectl get application -n argocd | awk '/payment-prod/{print $1}'| xargs  kubectl delete -n argocd application
```
