echo "SET AWS CREDS"
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

# TODO : Change customer_id to organization_name
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
