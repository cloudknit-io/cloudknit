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

export ECR_REPO=413422438110.dkr.ecr.us-east-1.amazonaws.com
cd ../../zlifecycle-il-operator

# hack to keep account id out of version control, see yaml.example file for more info
cp config/manager/kustomization.yaml.example config/manager/kustomization.yaml
make deploy
