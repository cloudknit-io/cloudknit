UpdateComponentStatus "${env_name}" "${team_name}" "${config_name}" "destroying" true

echo $show_output_start
echo "Executing terraform destroy..." 2>&1 | appendLogs "/tmp/apply_output.txt"
echo $show_output_end

zlifecycle-internal-cli state component patch \
  --company $customer_id \
  --team $team_name \
  --environment $env_name \
  --component $config_name \
  --status destroying \
  -u http://zlifecycle-state-manager."${customer_id}"-system.svc.cluster.local:8080 \
  -v

aws s3 cp s3://zlifecycle-$zl_env-tfplan-$customer_id/$team_name/$env_name/$config_name/tfplans/$config_reconcile_id terraform-plan --profile compuzest-shared

echo $show_output_start

if [ ! -z "$workspace" ];
then
    terraform workspace select $workspace || terraform workspace new $workspace
    echo "Workspace $workspace selected" 2>&1 | appendLogs /tmp/plan_output.txt
fi

((((terraform apply -auto-approve -input=false -no-color terraform-plan || returnErrorCode; echo $? >&3) 2>&1 | appendLogs "/tmp/apply_output.txt" >&4) 3>&1) | (read xs; exit $xs)) 4>&1
result=$?
if [ $result -eq 99 ]
then
 echo $show_output_end
 SaveAndExit "Failure during terraform destroy.";
fi
echo -n $result >/tmp/plan_code.txt
echo $show_output_end

aws s3 cp /tmp/apply_output.txt s3://zlifecycle-$zl_env-tfplan-$customer_id/$team_name/$env_name/$config_name/$config_reconcile_id/apply_output --profile compuzest-shared

if [ $result -eq 0 ]; then
    UpdateComponentStatus "${env_name}" "${team_name}" "${config_name}" "destroyed" true
else
    SaveAndExit "There is an issue with destroying"
fi

zlifecycle-internal-cli state component patch \
  --company $customer_id \
  --team $team_name \
  --environment $env_name \
  --component $config_name \
  --status destroyed \
  -u http://zlifecycle-state-manager."${customer_id}"-system.svc.cluster.local:8080 \
  -v
