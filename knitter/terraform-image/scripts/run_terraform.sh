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
is_apply=$4
lock_state=$5
is_sync=$6
workflow_id=$7
terraform_il_path=$8
is_destroy=$9
config_reconcile_id=${10}
reconcile_id=${11}
customer_id=${12}
auto_approve=${13}

#---------- INIT PHASE START ----------#

echo "Initializing..." 2>&1 | tee /tmp/$s3FileName.txt

. /initialize-component-variables.sh

. /initialize-functions.sh

sh /client/setup_github.sh || SaveAndExit "Cannot setup github ssh key"

sh /client/setup_aws.sh || SaveAndExit "Cannot setup aws credentials"

cd $ENV_COMPONENT_PATH

sh /argocd/login.sh $customer_id

# add last argo workflow run id to config application so it can fetch workflow details on UI
data='{"metadata":{"labels":{"last_workflow_run_id":"'$workflow_id'"}}}'
argocd app patch $team_env_config_name --patch $data --type merge >null

. /set-aws-creds.sh

. /initialize-terraform.sh

#---------- INIT PHASE END ----------#

if [ $is_apply -eq 0 ]; then
  if [ $is_sync -eq 1 ]; then
    sh /argocd/patch_env_component.sh $team_env_config_name
  fi

  echo "DEBUG: is_destroy: $is_destroy"

  result=1
  if [ $is_destroy = true ]; then
    . /terraform_destroy_plan.sh
  else
    . /terraform_apply_plan.sh
  fi

  if [ $result -eq 1 ]; then
    SaveAndExit "There is an issue with generating terraform plan"
  fi

  sh /argocd/process_based_on_plan_result.sh $is_sync $result $team_env_name $team_env_config_name $workflow_id $is_destroy $team_name $env_name $config_name $reconcile_id $config_reconcile_id $auto_approve $customer_id || SaveAndExit "There is an issue with ArgoCD CLI"

else
  if [ $is_destroy = true ]; then
    . /terraform_destroy.sh
  else
    . /terraform_apply.sh
  fi
  # argocd app sync $team_env_config_name >null
fi
