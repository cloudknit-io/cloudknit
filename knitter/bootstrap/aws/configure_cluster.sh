#!/bin/bash
set -eo pipefail

LOCATION=$1
APISERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')

kubectl create secret generic k8s-api --from-literal=url=$APISERVER -n zlifecycle-il-operator-system || true

argocd cluster add arn:aws:eks:us-east-1:413422438110:cluster/0-$LOCATION-eks --name $LOCATION
