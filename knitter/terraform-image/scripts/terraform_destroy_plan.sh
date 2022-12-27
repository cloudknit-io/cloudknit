echo "TERRAFORM DESTROY PLAN"
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

echo $show_output_start
echo "Executing plan..." 2>&1 | appendLogs /tmp/plan_output.txt
echo $show_output_end
UpdateComponentStatus "${env_name}" "${team_name}" "${config_name}" "running_destroy_plan"

echo $show_output_start

if [ ! -z "$workspace" ];
then
    terraform workspace select $workspace || terraform workspace new $workspace
    echo "Workspace $workspace selected" 2>&1 | appendLogs /tmp/plan_output.txt
fi

((((terraform plan -destroy -lock=$lock_state -input=false -no-color -out=terraform-plan -detailed-exitcode; echo $? >&3) 2>&1 | appendLogs "/tmp/plan_output.txt" >&4) 3>&1) | (read xs; exit $xs)) 4>&1
result=$?
echo -n $result >/tmp/plan_code.txt
echo $show_output_end

aws s3 cp /tmp/plan_output.txt s3://zlifecycle-$zl_env-tfplan-$customer_id/$team_name/$env_name/$config_name/$config_reconcile_id/plan_output --profile compuzest-shared
aws s3 cp terraform-plan s3://zlifecycle-$zl_env-tfplan-$customer_id/$team_name/$env_name/$config_name/tfplans/$config_reconcile_id --profile compuzest-shared

isDestroyed=true
costing_payload='{"teamName": "'$team_name'", "environmentName": "'$env_name'", "component": { "componentName": "'$config_name'", "isDestroyed" : '${isDestroyed}'  }}'
echo $costing_payload >temp_costing_payload.json

echo "TERRAFORM DESTROY PLAN - COSTING PAYLOAD"
echo $costing_payload

# TODO : add orgId to URL
# TODO : replace customer_id with generic multi-tenant API url
curl -X 'POST' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/costing/saveComponent" -H 'accept: */*' -H 'Content-Type: application/json' -d @temp_costing_payload.json
