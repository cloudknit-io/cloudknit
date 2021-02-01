#!/bin/bash
# Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
#
# Unauthorized copying of this file, via any medium, is strictly prohibited
# Proprietary and confidential
#
# NOTICE: All information contained herein is, and remains the property of
# CompuZest, Inc. The intellectual and technical concepts contained herein are
# proprietary to CompuZest, Inc. and are protected by trade secret or copyright
# law. Dissemination of this information or reproduction of this material is
# strictly forbidden unless prior written permission is obtained from CompuZest, Inc.

set -eo pipefail

LOCATION=$1
PARENT_DIRECTORY=$2



cd ../bootstrap/$PARENT_DIRECTORY

kubectl apply -f ecr-auth # create resources to allow local clusters to pull from ECR
kubectl patch workflowtemplate terraform-run-template -n argocd -p '{"imagePullSecrets": [{"name": "aws-registry"}]}' --type=merge # add ecr image pull secrets to argo workflow templates
kubectl patch workflowtemplate terraform-sync-template -n argocd -p '{"imagePullSecrets": [{"name": "aws-registry"}]}' --type=merge # add ecr image pull secrets to argo workflow templates
kubectl patch workflowtemplate workflow-trigger-template -n argocd -p '{"spec": {"imagePullSecrets": [{"name": "aws-registry"}]}}' --type=merge # add ecr image pull secrets to argo workflow templates

if [[ $(kubectl get job -n zlifecycle-ui | grep aws | wc -l) -eq 0 ]]
then
    kubectl create job --from=cronjob/aws-registry-credential-cron -n zlifecycle-ui aws-registry-initial-job
fi
if [[ $(kubectl get job -n zlifecycle-il-operator-system | grep aws | wc -l) -eq 0 ]]
then
    kubectl create job --from=cronjob/aws-registry-credential-cron -n zlifecycle-il-operator-system aws-registry-initial-job
fi
if [[ $(kubectl get job -n argocd | grep aws | wc -l) -eq 0 ]]
then
    kubectl create job --from=cronjob/aws-registry-credential-cron -n argocd aws-registry-initial-job
fi

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
