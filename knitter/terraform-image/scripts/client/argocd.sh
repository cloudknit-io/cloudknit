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

is_sync=$1
result=$2
team_env_name=$3
team_env_config_name=$4
workflow_id=$5

argoPassword=$(kubectl get secret argocd-server-login -n argocd -o json | jq '.data.password | @base64d' | tr -d '"')

echo y | argocd login --insecure argocd-server:443 --grpc-web --username admin --password $argoPassword

if [ $is_sync -eq 0 ]
then
    if [ $result -eq 2 ]
    then
        env_sync_status=$(argocd app get $team_env_name -o json | jq -r '.status.sync.status')
        config_sync_status=$(argocd app get $team_env_config_name -o json | jq -r '.status.sync.status')
        
        if [ $config_sync_status != "OutOfSync" ]
        then
            tfconfig="${team_env_config_name}-terraformconfig"

            argocd app patch-resource $team_env_config_name --kind TerraformConfig --resource-name $tfconfig --patch '{ "spec": { "isInSync": false } }' --patch-type 'application/merge-patch+json'

            if [ $env_sync_status != "OutOfSync" ]
            then
                argocd app sync $team_env_name
            fi
        fi
    else
        # update Argo Application to show in sync now that sync has manually taken place
        if [ $config_sync_status == "OutOfSync" ]
        then
            argocd app sync $team_env_config_name
        fi
    fi
else
    # add last argo workflow run id to config application so it can fetch workflow details on UI
    data='{"metadata":{"labels":{"last_workflow_run_id":"'$workflow_id'"}}}'
    argocd app patch $team_env_config_name --patch $data --type merge
fi
