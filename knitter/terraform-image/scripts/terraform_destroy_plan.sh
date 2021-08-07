echo $show_output_start
echo "Executing destroy plan..." 2>&1 | tee /tmp/plan_output.txt
echo $show_output_end
data='{"metadata":{"labels":{"component_status":"running_destroy_plan"}}}'
argocd app patch $team_env_config_name --patch $data --type merge >null

echo $show_output_start
((((terraform plan -destroy -lock=$lock_state -parallelism=2 -input=false -no-color -out=terraform-plan -detailed-exitcode; echo $? >&3) | appendLogs "/tmp/plan_output.txt" >&4) 3>&1) | (read xs; exit $xs)) 4>&1
result=$?
echo -n $result >/tmp/plan_code.txt
echo $show_output_end

aws s3 cp /tmp/plan_output.txt s3://zlifecycle-tfplan-$customer_id/$team_name/$env_name/$config_name/$config_reconcile_id/plan_output
aws s3 cp terraform-plan s3://zlifecycle-tfplan-$customer_id/$team_name/$env_name/$config_name/tfplans/$config_reconcile_id --profile compuzest-shared

costing_payload='{"teamName": "'$team_name'", "environmentName": "'$env_name'", "component": { "componentName": "'$config_name'", "isDeleted" : '1'  }}'
echo $costing_payload >temp_costing_payload.json
curl -X 'POST' 'http://zlifecycle-api.zlifecycle-ui.svc.cluster.local/costing/api/v1/saveComponent' -H 'accept: */*' -H 'Content-Type: application/json' -d @temp_costing_payload.json
