# Update Environment CRD

## Overview

When Environment CRD changes happen, it might break the operator & cause issues with existing environments that are already created. This runbook has steps on how to update environment CR in existing zLifecycle setup

## When to use this runbook

When Environment CRD changes and you want to promote changes to an existing zLifecycle setup

## Prerequisites

## Initial Steps Overview

- [Change Operator Replica to 0](#change-operator-replica-to-0)
- [Delete existing Environment CRs](#delete-existing-environment-crs)
- [Apply new crds](#apply-new-crds)
- [Apply Environment CR & restart Operator](#apply-environment-cr-restart-operator)

## Detailed Steps

### Change Operator Replica to 0
1. Edit Operator Deployment to have 0 replicas

### Delete existing Environment CRs
1. Update existing Environment CRs to remove finalizer using below command:

```bash
kubectl get environments -n zlifecycle | awk '//{print $1}' | xargs kubectl patch environment  -p '{"metadata":{"finalizers":[]}}' --type=merge -n zlifecycle
```

2. 

### Running the operator

#### IntelliJ
1. Install the plugin [EnvFile](https://plugins.jetbrains.com/plugin/7861-envfile)
2. Edit -> Edit configurations -> Add New Configuration -> Go Build -> select `Package` for `Run kind:`
3. Select the `EnvFile` tab -> Enable EnvFile -> Add -> Select your env file for your environment
4. Now you can run/debug your operator code: Run -> Run: '<configuration-name>' | Debug: '<configuration-name>'

## Other tools
1. Make sure all of the variables in the env file start with `export <key>=<value`
2. Run `source <environment>.env`
3. Build the operator and run the executable
