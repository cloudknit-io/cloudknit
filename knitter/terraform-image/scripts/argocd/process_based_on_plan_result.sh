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

is_sync=$1
result=$2
team_env_name=$3
team_env_config_name=$4
workflow_id=$5
is_destroy=$6
team_name=$7
env_name=$8
config_name=$9
reconcile_id=${10}
config_reconcile_id=${11}
auto_approve=${12}

function PatchError() {
    data='{"metadata":{"labels":{"component_status":"plan_failed"}}}'
    argocd app patch $team_env_config_name --patch $data --type merge > null
    sh ../audit.sh $team_name $env_name $config_name "" "plan_failed" $reconcile_id $config_reconcile_id $is_destroy 0 "noSkip"
    if [ $is_destroy = true ]
    then
        data='{"metadata":{"labels":{"env_status":"destroy_failed"}}}'
    else
        data='{"metadata":{"labels":{"env_status":"provision_failed"}}}'
    fi
    argocd app patch $team_env_name --patch $data --type merge > null
}

function Error() {
  if [ -n "$1" ];
  then
    echo "Error: "$1
    PatchError
  fi

    exit 1;
}

env_sync_status=$(argocd app get $team_env_name -o json | jq -r '.status.sync.status') || Error "Failed getting env_sync_status"

config_sync_status=$(argocd app get $team_env_config_name -o json | jq -r '.status.sync.status') || Error "Failed getting config_sync_status"

if [ $result -eq 0 ]
then
    if [ $is_destroy = true ]
    then
        data='{"metadata":{"labels":{"component_status":"destroyed"}}}'
    else
        data='{"metadata":{"labels":{"component_status":"provisioned"}}}'
    fi
    argocd app patch $team_env_config_name --patch $data --type merge > null

    if [ $config_sync_status == "OutOfSync" ]
    then
        argocd app sync $team_env_config_name > null || true
    fi
elif [ $result -eq 2 ]
then
    if [ $is_sync -eq 0 ]
    then
        if [ $config_sync_status != "OutOfSync" ]
        then
            data='{"metadata":{"labels":{"component_status":"out_of_sync"}}}'
            argocd app patch $team_env_config_name --patch $data --type merge > null

            sh /argocd/patch_env_component.sh $team_env_config_name || Error "Failed patching env component"
            if [ $env_sync_status != "OutOfSync" ]
            then
                argocd app sync $team_env_name > null || true
            fi
        fi
    elif [ $is_sync -eq 1 ]
    then
        if [ $auto_approve = false ]
        then
            data='{"metadata":{"labels":{"component_status":"waiting_for_approval","audit_status":"waiting_for_approval"}}}'
        else
            if [ $is_destroy = true ]
            then
                data='{"metadata":{"labels":{"component_status":"destroying","audit_status":"destroying"}}}'
            else
                data='{"metadata":{"labels":{"component_status":"provisioning","audit_status":"provisioning"}}}'
            fi
        fi
        argocd app patch $team_env_config_name --patch $data --type merge > null
    fi
fi
