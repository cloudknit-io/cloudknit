data='{"metadata":{"labels":{"component_status":"destroying","audit_status":"destroying"}}}'
argocd app patch $team_env_config_name --patch $data --type merge >null

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

# aws s3 cp s3://zlifecycle-$env_name-tfplan-$customer_id/$team_name/$env_name/$config_name/tfplans/$config_reconcile_id terraform-plan --profile compuzest-shared
sh /api_file.sh $team_name/$env_name/$config_name/tfplans/$config_reconcile_id terraform-plan $customer_id

echo $show_output_start
((((terraform apply -auto-approve -input=false -no-color terraform-plan || returnErrorCode; echo $? >&3) 2>&1 | appendLogs "/tmp/apply_output.txt" >&4) 3>&1) | (read xs; exit $xs)) 4>&1
result=$?
if [ $result -eq 99 ]
then
 echo $show_output_end
 SaveAndExit "Failure during terraform destroy.";
fi
echo -n $result >/tmp/plan_code.txt
echo $show_output_end



# aws s3 cp /tmp/apply_output.txt s3://zlifecycle-$env_name-tfplan-$customer_id/$team_name/$env_name/$config_name/$config_reconcile_id/apply_output --profile compuzest-shared
sh /api_file.sh "@/tmp/apply_output.txt" $team_name/$env_name/$config_name/$config_reconcile_id/apply_output $customer_id

if [ $result -eq 0 ]; then
    data='{"metadata":{"labels":{"component_status":"destroyed"}}}'
    argocd app patch $team_env_config_name --patch $data --type merge >null
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
