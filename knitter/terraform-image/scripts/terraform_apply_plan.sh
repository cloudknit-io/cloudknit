echo $show_output_start
echo "Executing apply plan..." 2>&1 | appendLogs /tmp/plan_output.txt
echo $show_output_end
data='{"metadata":{"labels":{"component_status":"running_plan"}}}'
argocd app patch $team_env_config_name --patch $data --type merge >null

echo $show_output_start
((((terraform plan -lock=$lock_state -input=false -no-color -out=terraform-plan -detailed-exitcode; echo $? >&3) 2>&1 | appendLogs "/tmp/plan_output.txt" >&4) 3>&1) | (read xs; exit $xs)) 4>&1
result=$?
echo -n $result >/tmp/plan_code.txt
echo $show_output_end

aws s3 cp /tmp/plan_output.txt s3://zlifecycle-tfplan-$customer_id/$team_name/$env_name/$config_name/$config_reconcile_id/plan_output --profile compuzest-shared
aws s3 cp terraform-plan s3://zlifecycle-tfplan-$customer_id/$team_name/$env_name/$config_name/tfplans/$config_reconcile_id --profile compuzest-shared

data='{"metadata":{"labels":{"component_status":"calculating_cost"}}}'
argocd app patch $team_env_config_name --patch $data --type merge >null

infracost breakdown --path terraform-plan --format json --log-level=warn >>output.json
estimated_cost=$(cat output.json | jq -r ".projects[0].breakdown.totalMonthlyCost")
resources=$(cat output.json | jq -r ".projects[0].breakdown.resources")
unsupportedResources=$(cat output.json | jq ' .projects[0].summary.unsupportedResourceCounts == {}')
if [[ $unsupportedResources != "true" && $estimated_cost == '0' ]]
then
#    estimated_cost='-1'
fi

costing_payload='{"teamName": "'$team_name'", "environmentName": "'$env_name'", "component": { "componentName": "'$config_name'", "cost": '$estimated_cost', "resources" : '$resources'  }}'
echo $costing_payload >temp_costing_payload.json

curl -X 'POST' 'http://zlifecycle-api.zlifecycle-ui.svc.cluster.local/costing/api/v1/saveComponent' -H 'accept: */*' -H 'Content-Type: application/json' -d @temp_costing_payload.json

data='{"metadata":{"labels":{"component_status":"provisioning"}}}'
argocd app patch $team_env_config_name --patch $data --type merge >null
