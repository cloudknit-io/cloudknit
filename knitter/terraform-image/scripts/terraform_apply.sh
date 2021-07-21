echo $show_output_start
echo "Executing terraform apply..."
echo $show_output_end
data='{"metadata":{"labels":{"component_status":"provisioning"}}}'

argocd app patch $team_env_config_name --patch $data --type merge >null

aws s3 cp s3://zlifecycle-tfplan-zmart/$team_name/$env_name/$config_name/$config_reconcile_id terraform-plan

echo $show_output_start
terraform apply -auto-approve -input=false -parallelism=2 -no-color terraform-plan || Error "Can not apply terraform plan"
result=$?
echo -n $result >/tmp/plan_code.txt
echo $show_output_end

if [ $result -eq 0 ]; then
    data='{"metadata":{"labels":{"component_status":"provisioned"}}}'
    argocd app patch $team_env_config_name --patch $data --type merge >null
else
    data='{"metadata":{"labels":{"component_status":"provision_failed"}}}'
    argocd app patch $team_env_config_name --patch $data --type merge >null

    Error "There is issue with provisioning"
fi
