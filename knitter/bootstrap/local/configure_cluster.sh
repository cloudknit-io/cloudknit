#!/bin/bash
set -eo pipefail

LOCATION=$1
PARENT_DIRECTORY=$2

cd $PARENT_DIRECTORY

kubectl apply -f pull-ecr-cron.yaml # create resources to allow local clusters to pull from ECR
kubectl create job --from=cronjob/aws-registry-credential-cron -n zlifecycle-il-operator-system aws-registry-initial-job


ip_addr=$(ipconfig getifaddr en0)

if [ ! $ip_addr ]
then
    ip_addr=$(ipconfig getifaddr en1)
fi

sed -i .bak "s+https://0.0.0.0:59999+https://$ip_addr:59999+g" ~/.kube/config

sleep 10s

curl --insecure https://$ip_addr:59999

sleep 10s

APISERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
kubectl create secret generic k8s-api --from-literal=url=$APISERVER -n zlifecycle-il-operator-system

argocd cluster add k3d-$LOCATION-k3d --insecure --name $LOCATION