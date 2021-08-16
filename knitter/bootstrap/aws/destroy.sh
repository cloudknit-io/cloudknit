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

checkForFailures() {
    if [ $? -ne 0 ]
    then
        echo ""
        echo "-------------------------------------"   
        read -p "Bootstrap phase has failed, type C to exit, any other key to continue" -n 1 -r
        echo ""

        if [[ $REPLY =~ ^[Cc]$ ]]
        then
            exit 1
        fi
    fi
}

cd ../../zlifecycle-provisioner/zlifecycle-addons
terraform init
terraform workspace select $LOCATION
terraform init
terraform destroy -auto-approve -var-file tfvars/$LOCATION.tfvars
checkForFailures

kubectl get application -n argocd | awk '//{print $1}'| xargs kubectl patch applications.argoproj.io -p '{"metadata":{"finalizers":[]}}' --type=merge -n argocd application

cd ../k8s-addons
terraform init
terraform workspace select $LOCATION
terraform init
terraform destroy -auto-approve -var-file tfvars/$LOCATION.tfvars
checkForFailures

cd ../aws-eks
terraform init
terraform workspace select $LOCATION
terraform init
terraform destroy -auto-approve -var-file tfvars/$LOCATION.tfvars
checkForFailures

cd ../aws-vpc
terraform init
terraform workspace select $LOCATION
terraform init
terraform destroy -auto-approve -var-file tfvars/$LOCATION.tfvars
