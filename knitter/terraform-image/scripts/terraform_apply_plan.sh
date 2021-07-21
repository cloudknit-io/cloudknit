echo $show_output_start
echo "Executing apply plan..."
echo $show_output_end
data='{"metadata":{"labels":{"component_status":"running_plan"}}}'
argocd app patch $team_env_config_name --patch $data --type merge >null

echo $show_output_start
terraform plan -lock=$lock_state -parallelism=2 -input=false -no-color -out=terraform-plan -detailed-exitcode
result=$?
echo -n $result >/tmp/plan_code.txt
echo $show_output_end

aws s3 cp terraform-plan s3://zlifecycle-tfplan-zmart/$team_name/$env_name/$config_name/$config_reconcile_id

data='{"metadata":{"labels":{"component_status":"calculating_cost"}}}'
argocd app patch $team_env_config_name --patch $data --type merge >null

infracost breakdown --path . --format json --log-level=warn >>output.json
estimated_cost=$(cat output.json | jq -r ".projects[0].breakdown.totalMonthlyCost")
resources=$(cat output.json | jq -r ".projects[0].breakdown.resources")

costing_payload='{"teamName": "'$team_name'", "environmentName": "'$env_name'", "component": { "componentName": "'$config_name'", "cost": '$estimated_cost', "resources" : '$resources'  }}'
echo $costing_payload >temp_costing_payload.json

curl -X 'POST' 'http://zlifecycle-api.zlifecycle-ui.svc.cluster.local/costing/api/v1/saveComponent' -H 'accept: */*' -H 'Content-Type: application/json' -d @temp_costing_payload.json
