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
if [[ $(kubectl get secret k8s-api -n zlifecycle-il-operator-system | wc -l) -eq 0 ]]
then
    kubectl create secret generic k8s-api --from-literal=url=$APISERVER -n zlifecycle-il-operator-system
fi

./create_ecr_secret.sh

argocd cluster add k3d-$LOCATION-k3d --insecure --name $LOCATION
