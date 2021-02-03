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

AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
REGION=us-east-1
SECRET_NAME=${REGION}-ecr-registry
EMAIL=zLifecycle@compuzest.com

TOKEN=`aws ecr --region=$REGION get-authorization-token --output text --query authorizationData[].authorizationToken | base64 -d | cut -d: -f2`

kubectl delete secret -n zlifecycle-il-operator-system --ignore-not-found $SECRET_NAME
kubectl create secret -n zlifecycle-il-operator-system docker-registry $SECRET_NAME \
 --docker-server=https://$AWS_ACCOUNT_ID.dkr.ecr.${REGION}.amazonaws.com \
 --docker-username=AWS \
 --docker-password="${TOKEN}" \
 --docker-email="${EMAIL}"

kubectl patch deployment zlifecycle-il-operator-controller-manager -n zlifecycle-il-operator-system  -p '{"spec": { "template": { "spec": {"imagePullSecrets": [{"name": "'${SECRET_NAME}'"}]}}}}'

kubectl create secret -n zlifecycle-ui docker-registry $SECRET_NAME \
 --docker-server=https://$AWS_ACCOUNT_ID.dkr.ecr.${REGION}.amazonaws.com \
 --docker-username=AWS \
 --docker-password="${TOKEN}" \
 --docker-email="${EMAIL}"

kubectl create secret -n argocd docker-registry $SECRET_NAME \
 --docker-server=https://$AWS_ACCOUNT_ID.dkr.ecr.${REGION}.amazonaws.com \
 --docker-username=AWS \
 --docker-password="${TOKEN}" \
 --docker-email="${EMAIL}"

