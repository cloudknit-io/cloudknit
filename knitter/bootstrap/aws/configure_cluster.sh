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
APISERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')

if [[ $(kubectl get secret k8s-api -n zlifecycle-il-operator-system | wc -l) -eq 0 ]]
then
    kubectl create secret generic k8s-api --from-literal=url=$APISERVER -n zlifecycle-il-operator-system || true
fi

argocd cluster add arn:aws:eks:us-east-1:413422438110:cluster/0-$LOCATION-eks --name $LOCATION

cd ..
kubectl apply ingress/.