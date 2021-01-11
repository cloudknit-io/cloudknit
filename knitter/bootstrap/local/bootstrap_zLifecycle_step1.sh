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
LOCAL=1

echo ""
echo ""
echo "-------------------------------------"
read -p "If you want create a k3d cluster enter Y: " -n 1 -r
echo ""
echo "-------------------------------------"

if [[ $REPLY =~ ^[Yy]$ ]]
then
    if ! docker info >/dev/null 2>&1; then
        echo "Docker does not seem to be running, run it first and retry"
        exit 1
    fi
    k3d cluster create $LOCATION-k3d -a 3 --api-port 59999
fi

cd ../../zlifecycle-provisioner/k8s-addons
terraform init
terraform workspace select 0-$LOCATION || terraform workspace new 0-$LOCATION
terraform init
terraform apply -auto-approve -var-file tfvars/$LOCATION.tfvars

