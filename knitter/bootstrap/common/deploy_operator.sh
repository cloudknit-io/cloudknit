#!/bin/bash

set -eo pipefail

cd ../../zlifecycle-il-operator

# hack to keep account id out of version control, see yaml.example file for more info
cp config/manager/kustomization.yaml.example config/manager/kustomization.yaml

make deploy IMG=$AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-il-operator:latest
