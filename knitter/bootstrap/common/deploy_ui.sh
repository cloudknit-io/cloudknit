
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

certArn=$(kubectl get secret ssl-cert-arn -o json | jq '.data.arn | @base64d' | tr -d '"')
export AWS_CERT_ARN=$certArn

#cd ../../zlifecycle-ui

#kubectl apply -f kubernetes/deployment.yaml
#kubectl apply -f kubernetes/service.yaml

cd ../../zlifecycle-api

kubectl apply -f kubernetes/deployment.yaml
kubectl apply -f kubernetes/service.yaml
