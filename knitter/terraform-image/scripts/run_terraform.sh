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

team_env_name=$team_name-$env_name
team_env_config_name=$team_name-$env_name-$config_name

ENV_COMPONENT_PATH=/home/terraform-config/$terraform_il_path

function Error() {
  if [ -n "$1" ];
  then
    echo "Error: "$1
  fi

    exit 1;
}

sh /client/setup_github.sh || Error "Cannot setup github ssh key"
sh /client/setup_aws.sh || Error "Cannot setup aws credentials"

cd $ENV_COMPONENT_PATH

sh /argocd/login.sh

data='{"metadata":{"labels":{"component_status":"initializing"}}}'
argocd app patch $team_env_config_name --patch $data --type merge

terraform init || Error "Cannot initialize terraform"

if [ $is_apply -eq 0 ]
then
    if [ $is_sync -eq 1 ]
    then
        sh /argocd/patch_env_component.sh $team_env_config_name

        # add last argo workflow run id to config application so it can fetch workflow details on UI
        data='{"metadata":{"labels":{"last_workflow_run_id":"'$workflow_id'"}}}'
        argocd app patch $team_env_config_name --patch $data --type merge
    fi

    echo "DEBUG: is_destroy: $is_destroy"

    result=1
    if [ $is_destroy = true ]
    then
      echo "Executing destroy plan..."
      data='{"metadata":{"labels":{"component_status":"running_destroy_plan"}}}'
      argocd app patch $team_env_config_name --patch $data --type merge

      terraform plan -destroy -lock=$lock_state -parallelism=2 -input=false -no-color -out=terraform-plan -detailed-exitcode
      result=$?
      echo -n $result > /tmp/plan_code.txt
    else
      echo "Executing apply plan..."
      data='{"metadata":{"labels":{"component_status":"running_plan"}}}'
      argocd app patch $team_env_config_name --patch $data --type merge

      terraform plan -lock=$lock_state -parallelism=2 -input=false -no-color -out=terraform-plan -detailed-exitcode
      result=$?
      echo -n $result > /tmp/plan_code.txt

      data='{"metadata":{"labels":{"component_status":"calculating_cost"}}}'
      argocd app patch $team_env_config_name --patch $data --type merge

      infracost breakdown --path . --format json >> output.json
      estimated_cost=$(cat output.json | jq -r ".projects[0].breakdown.totalMonthlyCost")

      data='{"metadata":{"labels":{"component_cost":"'$estimated_cost'"}}}'
      argocd app patch $team_env_config_name --patch $data --type merge
    fi

    if [ $result -eq 1 ]
    then
        data='{"metadata":{"labels":{"component_status":"plan_failed"}}}'
        argocd app patch $team_env_config_name --patch $data --type merge

        Error "There is issue with generating terraform plan"
    fi

    sh /argocd/process_based_on_plan_result.sh $is_sync $result $team_env_name $team_env_config_name $workflow_id || Error "There is an issue with ArgoCD CLI"

else
    if [ $is_destroy = true ]
    then
      echo "Executing terraform destroy..."
      data='{"metadata":{"labels":{"component_status":"destroying"}}}'
      argocd app patch $team_env_config_name --patch $data --type merge

      terraform destroy -auto-approve -input=false -parallelism=2 -no-color || Error "Cannot run terraform destroy"
      result=$?
      echo -n $result > /tmp/plan_code.txt

      if [ $result -eq 0 ]
      then
        data='{"metadata":{"labels":{"component_status":"destroyed"}}}'
        argocd app patch $team_env_config_name --patch $data --type merge
      else
        data='{"metadata":{"labels":{"component_status":"destroy_failed"}}}'
        argocd app patch $team_env_config_name --patch $data --type merge

        Error "There is issue with destroying"
      fi
    else
      echo "Executing terraform apply..."
      data='{"metadata":{"labels":{"component_status":"provisioning"}}}'
      argocd app patch $team_env_config_name --patch $data --type merge

      terraform apply -auto-approve -input=false -parallelism=2 -no-color || Error "Can not apply terraform plan"
      result=$?
      echo -n $result > /tmp/plan_code.txt

      if [ $result -eq 0 ]
      then
        data='{"metadata":{"labels":{"component_status":"provisioned"}}}'
        argocd app patch $team_env_config_name --patch $data --type merge
      else
        data='{"metadata":{"labels":{"component_status":"provision_failed"}}}'
        argocd app patch $team_env_config_name --patch $data --type merge

        Error "There is issue with provisioning"
      fi

    fi

    argocd app sync $team_env_config_name
fi
