
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

cd ../../zlifecycle-ui/kubernetes

sed -i '' "s/\${AWS_ACCOUNT_ID}/${AWS_ACCOUNT_ID}/g" deployment.yaml

# set cert for ALB Ingress
certArn=$(kubectl get secret ssl-cert-arn -o json | jq '.data.arn | @base64d' | tr -d '"')
sed -i '' "s/\${AWS_CERT_ARN}/${certArn}/g" ingress-alb.yaml

kubectl apply -f .
