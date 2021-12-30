team_env_name=$team_name-$env_name
team_env_config_name=$team_name-$env_name-$config_name
show_output_start='----->show_output_start<-----'
show_output_end='----->show_output_end<-----'
is_debug=0

# Intialize error statuses and s3FileName paths.
if [ $is_apply -eq 0 ];
then
  component_error_status="plan_failed"
  s3FileName="plan_output"
else
  component_error_status="apply_failed"
  s3FileName="apply_output"
fi


ENV_COMPONENT_PATH=/home/terraform-config/$terraform_il_path