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

main() {
    announcePhase() {
        echo ""
        echo ""
        echo "-------------------------------------"
        echo $1
        echo ""
        echo "-------------------------------------"   
    }

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

    if [[ -z "$AWS_ACCOUNT_ID" ]]
    then
        echo "Error: Please set \$AWS_ACCOUNT_ID for ECR"
        exit 1
    fi

    echo "Please select the environment you wish to bootstrap:"
    select LOCATION in "dev-a" "dev-b" "sandbox"; do
        if [[ $LOCATION == "sandbox" ]]
        then
            readonly PARENT_DIRECTORY=aws
            break;
        else
            readonly PARENT_DIRECTORY=local
            break;
        fi
    done

    # Provision with Terraform
    announcePhase "Provision"
    ./$PARENT_DIRECTORY/provision.sh $LOCATION;
    checkForFailures

    # Deploy Operator
    announcePhase "Deploying zlifecycle-il-operator"
    ./common/deploy_operator.sh
    checkForFailures

    # Bootstrap customers
    announcePhase "Bootstrap customer Environments"
    ./$PARENT_DIRECTORY/bootstrap_customers.sh $LOCATION $PARENT_DIRECTORY;
    checkForFailures

    #Prepare secrets
    echo ""
    echo ""
    echo "-------------------------------------"
    read -p "Please create secrets and enter Y to continue? " -n 1 -r
    echo ""
    echo "-------------------------------------"
    echo ""
    echo ""

    #Set up Argo
    announcePhase "Set up Argo"
    ./common/configure_argo.sh
    checkForFailures

    #Configure Cluster
    announcePhase "Configure cluster"
    ./$PARENT_DIRECTORY/configure_cluster.sh $LOCATION $PARENT_DIRECTORY
    checkForFailures
}

main
