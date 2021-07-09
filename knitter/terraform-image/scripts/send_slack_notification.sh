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

# exit when any command fails
set -e

workflow_name=$1
team_name=$2
env_name=$3
config_name=$4

team_env_config_name=$team_name-$env_name-$config_name

namespace="argocd"

message="${team_env_config_name} terraform is out of sync. To see the diff & approve the sync to desired state go here: http://localhost:8081/workflows/${namespace}/${workflow_name}"

data='{"channel": "slack-notification","message": "'$message'"}'
echo $data

curl -d "${data}" -H "Content-Type: application/json" -X POST http://terraform-diff-eventsource-svc.${namespace}.svc.cluster.local:12000/terraform-diff
