echo $show_output_start
echo "Executing terraform destroy..." 2>&1 | tee /tmp/apply_output.txt
echo $show_output_end
data='{"metadata":{"labels":{"component_status":"destroying"}}}'

argocd app patch $team_env_config_name --patch $data --type merge >null
aws s3 cp s3://zlifecycle-tfplan-zmart/$team_name/$env_name/$config_name/tf_plans/$config_reconcile_id terraform-plan

echo $show_output_start
terraform apply -auto-approve -input=false -parallelism=2 -no-color terraform-plan || Error "Can not apply terraform destroy" 2>&1 | tee -a /tmp/apply_output.txt
result=$?
echo -n $result >/tmp/plan_code.txt
echo $show_output_end

aws s3 cp /tmp/apply_output.txt s3://zlifecycle-tfplan-zmart/$team_name/$env_name/$config_name/$config_reconcile_id/apply_output

if [ $result -eq 0 ]; then
    data='{"metadata":{"labels":{"component_status":"destroyed"}}}'
    argocd app patch $team_env_config_name --patch $data --type merge >null
else
    data='{"metadata":{"labels":{"component_status":"destroy_failed"}}}'
    argocd app patch $team_env_config_name --patch $data --type merge >null

    Error "There is issue with destroying"
fi
