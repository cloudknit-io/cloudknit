echo "TERRAFORM APPLY PLAN"
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
UpdateComponentReconcile "${team_name}" "${env_name}" "${config_name}" '{ "status" : "running_plan" }' 

echo $show_output_start

if [ ! -z "$workspace" ];
then
    terraform workspace select $workspace || terraform workspace new $workspace
    echo "Workspace $workspace selected" 2>&1 | appendLogs /tmp/plan_output.txt
fi

((((terraform plan -lock=$lock_state -input=false -no-color -out=terraform-plan -detailed-exitcode; echo $? >&3) 2>&1 | appendLogs "/tmp/plan_output.txt" >&4) 3>&1) | (read xs; exit $xs)) 4>&1
result=$?
echo -n $result >/tmp/plan_code.txt
echo $show_output_end

echo "AWS S3 COPY"
echo "   zl_env=${zl_env}"
echo "   customer_id=${customer_id}"
echo "   team_name=${team_name}"
echo "   env_name=${env_name}"
echo "   config_name=${config_name}"
echo "   config_reconcile_id=${config_reconcile_id}"


aws s3 cp /tmp/plan_output.txt s3://zlifecycle-$zl_env-tfplan-$customer_id/$team_name/$env_name/$config_name/$config_reconcile_id/plan_output --profile compuzest-shared
aws s3 cp terraform-plan s3://zlifecycle-$zl_env-tfplan-$customer_id/$team_name/$env_name/$config_name/tfplans/$config_reconcile_id --profile compuzest-shared

UpdateComponentReconcile "${team_name}" "${env_name}" "${config_name}" '{ "status" : "calculating_cost" }' 

infracost breakdown --path terraform-plan --format json --log-level=warn >>output.json

estimated_cost=$(cat output.json | jq -r ".projects[0].breakdown.totalMonthlyCost")
costResources=$(cat output.json | jq -r ".projects[0].breakdown.resources")

UpdateComponentReconcile "${team_name}" "${env_name}" "${config_name}" '{ "estimatedCost" : '${estimated_cost}', "costResources" : '${costResources}' }'