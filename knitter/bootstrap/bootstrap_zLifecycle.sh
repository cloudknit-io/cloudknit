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
if [[ -z "$AWS_ACCOUNT_ID" ]]
then
    echo "Error: Please set \$AWS_ACCOUNT_ID"
    exit 1
fi

echo "Please select the environment you wish to bootstrap:"
select LOCATION in "dev-a" "dev-b" "sandbox"; do
    case $LOCATION in
        "dev-a" ) LOCAL=1 ./local/bootstrap_zLifecycle_step1.sh $LOCATION; break;;
        "dev-b" ) LOCAL=1 ./local/bootstrap_zLifecycle_step1.sh $LOCATION; break;;
        "sandbox" ) LOCAL=0 ./aws/bootstrap_zLifecycle_step1.sh; break;;
    esac
done

cd ../../zlifecycle-il-operator
echo "Deploying zlifecycle-il-operator"
make deploy IMG=$AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/zlifecycle-il-operator:latest

LOCAL=1
if [[ $LOCATION == "sandbox" ]]
then
    LOCAL=0
fi

if [[ $LOCAL -eq 1 ]]
then
    cd ../zLifecycle/bootstrap/local
    kubectl apply -f company-config-$LOCATION.yaml

    # Create all team environments
    cd ../../../compuzest-$LOCATION-zlifecycle-config
    kubectl apply -R -f teams/account-team
    kubectl apply -R -f teams/user-team
else
    cd ../zLifecycle/bootstrap/aws
    kubectl apply -f company-config.yaml

    # Create all team environments
    cd ../../../compuzest-zlifecycle-config
    kubectl apply -R -f teams/account-team
    kubectl apply -R -f teams/user-team
fi

echo ""
echo ""
echo "-------------------------------------"
read -p "Please create secrets and enter Y to continue? " -n 1 -r
echo ""
echo "-------------------------------------"

if [[ $REPLY =~ ^[Yy]$ ]]
then
    if [[ $LOCAL -eq 1 ]]
    then
        echo $(pwd)
        cd ../zLifecycle/bootstrap/local
        kubectl apply -f pull-ecr-cron.yaml # create resources to allow local clusters to pull from ECR
        kubectl create job --from=cronjob/aws-registry-credential-cron -n zlifecycle-il-operator-system aws-registry-initial-job
    fi
    cd ..
    ./common/bootstrap_zLifecycle_step2.sh $LOCATION $LOCAL
fi
