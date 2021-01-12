#!/bin/bash

set -eo pipefail

LOCATION=$1
PARENT_DIRECTORY=$2

cd $PARENT_DIRECTORY

kubectl apply -f company-config-$LOCATION.yaml

# Create all team environments
cd ../../../compuzest-$LOCATION-zlifecycle-config
kubectl apply -R -f teams/account-team
kubectl apply -R -f teams/user-team