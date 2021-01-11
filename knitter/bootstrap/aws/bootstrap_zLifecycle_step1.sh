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

cd ../../zlifecycle-provisioner/aws-vpc
terraform init
terraform workspace select 0-$LOCATION || terraform workspace new 0-$LOCATION
terraform init

terraform apply -auto-approve -var-file tfvars/$LOCATION.tfvars

cd ../aws-eks
terraform init
terraform workspace select 0-$LOCATION || terraform workspace new 0-$LOCATION
terraform init
terraform apply -auto-approve -var-file tfvars/$LOCATION.tfvars
sleep 2m
terraform apply -auto-approve -var-file tfvars/$LOCATION.tfvars
sleep 2m
aws eks --region us-east-1 update-kubeconfig --name 0-$LOCATION-eks

cd ../k8s-addons
terraform init
terraform workspace select 0-$LOCATION || terraform workspace new 0-$LOCATION
terraform init
terraform apply -auto-approve -var-file tfvars/$LOCATION.tfvars
