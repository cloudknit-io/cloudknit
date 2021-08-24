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
is_apply=$6
lock_state=$7
is_sync=$8
workflow_id=$9
terraform_il_path=$10
is_destroy=$11
config_reconcile_id=$12
reconcile_id=$13
customer_id=$14

team_env_name=$team_name-$env_name
team_env_config_name=$team_name-$env_name-$config_name
show_output_start='----->show_output_start<-----'
show_output_end='----->show_output_end<-----'
is_debug=0

ENV_COMPONENT_PATH=/home/terraform-config/$terraform_il_path

function PatchError() {
  if [ $is_destroy = true ]
  then
      data='{"metadata":{"labels":{"env_status":"destroy_failed"}}}'
  else
      data='{"metadata":{"labels":{"env_status":"provision_failed"}}}'
  fi
  
  argocd app patch $team_env_name --patch $data --type merge > null
  sh /audit.sh $team_name $env_name $config_name "" "Failed" $reconcile_id $config_reconcile_id $is_destroy 0
}

function appendLogs() {
  IFS=''
  while read line; do
    echo $line | tee -a $1 
  done
}

function Error() {
  if [ -n "$1" ]; then
    echo "Error: "$1
    PatchError
  fi

  exit 1
}

function SaveAndExit() {
  aws s3 cp /tmp/apply_output.txt s3://zlifecycle-tfplan-$customer_id/$team_name/$env_name/$config_name/$config_reconcile_id/apply_output --profile compuzest-shared
  Error $1
}

sh /client/setup_github.sh || Error "Cannot setup github ssh key"
sh /client/setup_aws.sh || Error "Cannot setup aws credentials"

cd $ENV_COMPONENT_PATH

sh /argocd/login.sh

data='{"metadata":{"labels":{"component_status":"initializing"}}}'
argocd app patch $team_env_config_name --patch $data --type merge >null

terraform init || Error "Cannot initialize terraform"

if [ $is_apply -eq 0 ]; then
  if [ $is_sync -eq 1 ]; then
    sh /argocd/patch_env_component.sh $team_env_config_name

    # add last argo workflow run id to config application so it can fetch workflow details on UI
    data='{"metadata":{"labels":{"last_workflow_run_id":"'$workflow_id'"}}}'
    argocd app patch $team_env_config_name --patch $data --type merge >null
  fi

  echo "DEBUG: is_destroy: $is_destroy"

  result=1
  if [ $is_destroy = true ]; then
    . /terraform_destroy_plan.sh
  else
    . /terraform_apply_plan.sh
  fi

  if [ $result -eq 1 ]; then
    data='{"metadata":{"labels":{"component_status":"plan_failed"}}}'
    argocd app patch $team_env_config_name --patch $data --type merge >null

    Error "There is an issue with generating terraform plan"
  fi

  sh /argocd/process_based_on_plan_result.sh $is_sync $result $team_env_name $team_env_config_name $workflow_id $is_destroy $team_name $env_name $config_name $reconcile_id $config_reconcile_id || Error "There is an issue with ArgoCD CLI"

else
  if [ $is_destroy = true ]; then
    . /terraform_destroy.sh
  else
    . /terraform_apply.sh
  fi
  argocd app sync $team_env_config_name >null
fi
