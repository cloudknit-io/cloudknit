# Bootstrap only zlifecycle-il-operator locally

## Overview

Create only bare minimum needed to run `zlifecycle-il-operator` locally without provisioning the entire zlifecycle platform

## When to use this runbook
When you want to bootstrap only `zLifecycle-il-operator` locally so you can test it without provisioning everything else

## Initial Steps Overview

1. [Prerequisites](#prerequisites)
1. [Bootstrap zlifecycle-il-operator locally](#bootstrap-zlifecycle-il-operator-locally)
1. [Testing](#testing)
1. [Destroy](#destroy)

## Detailed Steps

#### Prerequisites

1. [K3D or a similar tool to run k8s locally](https://k3d.io/)

#### Bootstrap zLifecycle-il-operator locally

Run script to bootstrap operator.

Note: When asked enter secrets (currently only 3 secrets but might change) for `zlifecycle-il-operator` namespace from Lastpass.

For now ignore the following error towards the end:
```
FATA[0000] rpc error: code = Unknown desc = Post "http://localhost:8080/cluster.ClusterService/Create": dial tcp [::1]:8080: connect: connection refused
```

```bash
cd ../../bootstrap
./bootstrap_local_operator.sh
```

#### Testing

For testing purposes following these steps: 

1. Create an environment custom resource by running following command against a yaml in config repo. For example: https://github.com/CompuZest/zmart-design-team-config/blob/main/dev/dev-env.yaml 

```bash
kubectl apply -f dev-env.yaml
```
1. Check the IL repo to see any changes
1. Make changes you want in the operator code and deploy the updated operator to local cluster by using following commands:
```bash
make docker-build IMG=$AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-il-operator:testing
make docker-push IMG=$AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-il-operator:testing
make deploy IMG=$AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-il-operator:testing
```
1. Once updated operator runs make sure there are only expected changes in the IL repo

Note: We will have better tests including unit tests that will allow for faster feedback in future but for now this is the best way to test operator locally

#### Destroy

Just run following command to destroy the cluster. No need to do anything else.

```bash
k3d cluster delete [replace with cluster-name]
```
