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

PARENT_DIRECTORY=$2
cd $PARENT_DIRECTORY

kubectl apply -f company-config.yaml

# Create all team environments
cd ../../../compuzest-zlifecycle-config
kubectl apply -R -f teams/account-team
kubectl apply -R -f teams/user-team
