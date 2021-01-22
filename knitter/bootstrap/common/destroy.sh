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

kubectl port-forward service/argo-cd-argocd-server 8080:80 -n argocd &

argoPassword=$(kubectl get secret argocd-server-login -n argocd -o json | jq '.data.password | @base64d' | tr -d '"')
yes Y | argocd login --insecure localhost:8080 --grpc-web --username admin --password $argoPassword

kubectl get applications -n argocd -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | xargs kubectl patch applications  -p '{"metadata":{"finalizers":[]}}' --type=merge -n argocd
kubectl patch crd applications.argoproj.io -p '{"metadata":{"finalizers":[]}}' --type=merge -n argocd
argocd app delete $(argocd app list -o name)
