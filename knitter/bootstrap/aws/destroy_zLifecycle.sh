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

LOCATION=$1
if [[ -z "$LOCATION" ]]
then
    echo "Error: Please pass the name of the environment you'd like to destroy to this script"
    exit 1
fi

argocd_server_name=$(kubectl get pods -l app.kubernetes.io/name=argocd-server -n argocd --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')
argocd login --insecure localhost:8080 --grpc-web --username admin --password $argocd_server_name

argocd app delete 1-customer
argocd cluster rm arn:aws:eks:us-east-1:413422438110:cluster/0-$LOCATION-eks
argocd repo rm git@github.com:CompuZest/infra-deploy-terraform-config.git
argocd repo rm git@github.com:CompuZest/helm-charts.git

cd ../../../zlifecycle-provisioner/k8s-addons
terraform init
terraform workspace select 0-$LOCATION
terraform init
terraform destroy -auto-approve -var-file tfvars/$LOCATION.tfvars

cd ../aws-eks
terraform init
terraform workspace select 0-$LOCATION
terraform init
terraform destroy -auto-approve -var-file tfvars/$LOCATION.tfvars


cd ../aws-vpc
terraform init
terraform workspace select 0-$LOCATION
terraform init
terraform destroy -auto-approve -var-file tfvars/$LOCATION.tfvars
