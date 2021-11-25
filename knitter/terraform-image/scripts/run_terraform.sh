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

team_env_name=$team_name-$env_name
team_env_config_name=$team_name-$env_name-$config_name
show_output_start='----->show_output_start<-----'
show_output_end='----->show_output_end<-----'
is_debug=0
if [ $is_apply -eq 0 ];
then
  component_error_status="plan_failed"
  s3FileName="plan_output"
else
  s3FileName="apply_output"
  if [ $is_destroy = true ]
  then
      component_error_status="destroy_failed"
  else
      component_error_status="provision_failed"
  fi
fi

echo "Initializing..." 2>&1 | tee /tmp/$s3FileName.txt


ENV_COMPONENT_PATH=/home/terraform-config/$terraform_il_path

function PatchError() {
  if [ $is_destroy = true ]
  then
      data='{"metadata":{"labels":{"env_status":"destroy_failed"}}}'
  else
      data='{"metadata":{"labels":{"env_status":"provision_failed"}}}'
  fi

  argocd app patch $team_env_name --patch $data --type merge > null

  data='{"metadata":{"labels":{"component_status":"'$component_error_status'"}}}'
  argocd app patch $team_env_config_name --patch $data --type merge >null

  sh /audit.sh $team_name $env_name $config_name "Failed" "Failed" $reconcile_id $config_reconcile_id $is_destroy 0
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
  echo $show_output_start
  echo $1 2>&1 | appendLogs /tmp/$s3FileName.txt
  echo $show_output_end
  aws s3 cp /tmp/$s3FileName.txt s3://zlifecycle-tfplan-$customer_id/$team_name/$env_name/$config_name/$config_reconcile_id/$s3FileName --profile compuzest-shared --quiet
  Error $1
}

function returnErrorCode() {
  return 99;
}

function setAWSCreds() {
  aws_region=$(aws ssm get-parameter --profile compuzest-shared --region us-east-1 --name "/$1/aws_region" --with-decryption --query "Parameter.Value" | jq -r ".")
  if [ ! -z $aws_region ];
  then
    export AWS_REGION=$aws_region
  fi

  aws_access_key_id=$(aws ssm get-parameter --profile compuzest-shared --region us-east-1 --name "/$1/aws_access_key_id" --with-decryption --query "Parameter.Value" | jq -r ".")
  aws_secret_access_key=$(aws ssm get-parameter --profile compuzest-shared --region us-east-1 --name "/$1/aws_secret_access_key" --with-decryption --query "Parameter.Value" | jq -r ".")

  if [ ! -z $aws_access_key_id -a ! -z $aws_secret_access_key ];
  then
    aws configure set aws_access_key_id $aws_access_key_id 
    aws configure set aws_secret_access_key $aws_secret_access_key

    aws_session_token=$(aws ssm get-parameter --profile compuzest-shared --region us-east-1 --name "/$1/aws_session_token" --with-decryption --query "Parameter.Value" | jq -r ".")

    if [ ! -z $aws_session_token ];
    then
      aws configure set aws_session_token $aws_session_token
    fi

    return 1
  fi
  return 0
}

sh /client/setup_github.sh || SaveAndExit "Cannot setup github ssh key"
sh /client/setup_aws.sh || SaveAndExit "Cannot setup aws credentials"

cd $ENV_COMPONENT_PATH

sh /argocd/login.sh

setAWSCreds $customer_id/$team_name/$env_name
aws_response=$?
if [ $aws_response -eq 0 ];
then
  setAWSCreds $customer_id/$team_name
  aws_response=$?
  if [ $aws_response -eq 0 ];
  then
    setAWSCreds $customer_id
    aws_response=$?
    if [ $aws_response -eq 0 ];
    then
      SaveAndExit "No AWS Credentials available. Please set AWS Credentials in the Settings Page."
    fi
  fi
fi

data='{"metadata":{"labels":{"component_status":"initializing"}}}'
argocd app patch $team_env_config_name --patch $data --type merge >null

echo $show_output_start
((((terraform init; echo $? >&3) 2>&1 1>/dev/null | appendLogs "/tmp/$s3FileName.txt" >&4) 3>&1) | (read xs; exit $xs)) 4>&1
if [ $? -ne 0 ]; then
  echo $show_output_end
  SaveAndExit "Failed to initialize terraform"
fi
echo $show_output_end

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
    SaveAndExit "There is an issue with generating terraform plan"
  fi

  sh /argocd/process_based_on_plan_result.sh $is_sync $result $team_env_name $team_env_config_name $workflow_id $is_destroy $team_name $env_name $config_name $reconcile_id $config_reconcile_id $auto_approve || SaveAndExit "There is an issue with ArgoCD CLI"

else
  if [ $is_destroy = true ]; then
    . /terraform_destroy.sh
  else
    . /terraform_apply.sh
  fi
  argocd app sync $team_env_config_name >null
fi
