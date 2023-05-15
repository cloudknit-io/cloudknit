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
customer_id=${12} # TODO : Change to organization_name
auto_approve=${13}
zl_env=${14}
git_auth_mode=${15}
il_repo=${16}
company_git_org=${17}
use_custom_state=${18}
custom_state_bucket=${19}
custom_state_lock_table=${20}
workspace=${21}

echo "RUN TERRAFORM"
echo "   team_name=${team_name}"
echo "   env_name=${env_name}"
echo "   config_name=${config_name}"
echo "   is_apply=${is_apply}"
echo "   lock_state=${lock_state}"
echo "   is_sync=${is_sync}"
echo "   workflow_id=${workflow_id}"
echo "   terraform_il_path=${terraform_il_path}"
echo "   is_destroy=${is_destroy}"
echo "   config_reconcile_id=${config_reconcile_id}"
echo "   reconcile_id=${reconcile_id}"
echo "   customer_id=${customer_id}"
echo "   auto_approve=${auto_approve}"
echo "   zl_env=${zl_env}"
echo "   git_auth_mode=${git_auth_mode}"
echo "   il_repo=${il_repo}"
echo "   company_git_org=${company_git_org}"
echo "   use_custom_state=${use_custom_state}"
echo "   custom_state_bucket=${custom_state_bucket}"
echo "   custom_state_lock_table=${custom_state_lock_table}"
echo "   workspace=${workspace}"

#---------- INIT PHASE START ----------#

# TODO : Get organization ID from organization_name

echo "Initializing..." 2>&1 | tee /tmp/$s3FileName.txt

sh /argocd/login.sh $customer_id

. /initialize-component-variables.sh

. /initialize-functions.sh


UpdateComponentReconcile "${team_name}" "${env_name}" "${config_name}" '{ "lastWorkflowRunId" : "'${workflow_id}'" }'

sh /client/setup_github.sh || SaveAndExit "Cannot setup github ssh key"

sh /client/setup_aws.sh || SaveAndExit "Cannot setup aws credentials"

internal_git_auth="github-app-internal"
zlifecycle-internal-cli git clone $il_repo              \
  --git-auth $internal_git_auth                         \
  --git-ssh /root/internal_github_app_ssh/sshPrivateKey \
  --dir $IL_REPO_PATH                                   \
  -v
cd $ENV_COMPONENT_PATH

if [ $git_auth_mode != "ssh" ]; then
  zlifecycle-internal-cli git login $company_git_org    \
    --git-auth $git_auth_mode                           \
    --git-ssh /root/public_github_app_ssh/sshPrivateKey \
    --git-config-path $HOME                             \
    -v
fi

. /set-aws-creds.sh

if [ $use_custom_state == "true" ]; then
  zlifecycle-internal-cli aws configure \
    --auth-mode profile                 \
    --profile compuzest-shared          \
    --generated-profile customer-state  \
    --company $customer_id              \
    --team $team_name                   \
    --environment $env_name             \
    --verbose
fi

. /initialize-terraform.sh

#---------- INIT PHASE END ----------#

if [ $is_apply -eq 0 ]; then

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

  # TODO : Pass orgId
  sh /argocd/process_based_on_plan_result.sh $is_sync $result $team_env_name $team_env_config_name $workflow_id $is_destroy $team_name $env_name $config_name $reconcile_id $config_reconcile_id $auto_approve $customer_id || SaveAndExit "There is an issue with ArgoCD CLI"

else
  if [ $is_destroy = true ]; then
    . /terraform_destroy.sh
  else
    . /terraform_apply.sh
  fi
  # argocd app sync $team_env_config_name >null
fi
