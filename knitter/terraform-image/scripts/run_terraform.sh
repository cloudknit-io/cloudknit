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

team_name=$1
env_name=$2
config_name=$3
module_source=$4
module_source_path=$5
variables_file_source=$6
variables_file_path=$7
is_apply=$8
lock_state=$9
is_sync=$10
workflow_id=$11

team_env_name=$team_name-$env_name
team_env_config_name=$team_name-$env_name-$config_name

ENV_COMPONENT_PATH=/home/terraform-config

function Error() {
  if [ -n "$1" ];
  then
    echo "Error: "$1
  fi

    exit 1;
}

sh /client/setup_github.sh || Error "Cannot setup github ssh key"
sh /client/setup_aws.sh || Error "Cannot setup aws credentials"

sh /terraform/provider.tf.sh $ENV_COMPONENT_PATH || Error "Cannot generate terraform provider"
sh /terraform/module.tf.sh $ENV_COMPONENT_PATH $config_name $module_source $module_source_path $variables_file_path || Error "Cannot generate terraform module"
sh /terraform/terraform.tf.sh $ENV_COMPONENT_PATH $team_name $team_env_name $config_name || Error "Cannot generate terraform state block"

cd $ENV_COMPONENT_PATH

terraform init || Error "Cannot initialize terraform"

sh /argocd/login.sh

if [ $is_apply -eq 0 ]
then

    if [ $is_sync -eq 1 ]
    then
        sh /argocd/patch_env_component.sh $team_env_config_name

        # add last argo workflow run id to config application so it can fetch workflow details on UI
        data='{"metadata":{"labels":{"last_workflow_run_id":"'$workflow_id'"}}}'
        argocd app patch $team_env_config_name --patch $data --type merge
    fi

    terraform plan -lock=$lock_state -parallelism=2 -input=false -no-color -out=terraform-plan -detailed-exitcode
    result=$?
    echo -n $result > /tmp/plan_code.txt

    if [ $result -eq 1 ]
    then
        Error "There is issue with generating terraform plan"
    fi

    sh /argocd/control_loop.sh $is_sync $result $team_env_name $team_env_config_name $workflow_id || Error "There is an issue with ArgoCD CLI"

else

    terraform apply -auto-approve -input=false -parallelism=2 -no-color || Error "Can not apply terraform plan"
    echo -n 0 > /tmp/plan_code.txt
    argocd app sync $team_env_config_name

fi
