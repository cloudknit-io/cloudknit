# Patches the error status of argocd app
function PatchError() {
  if [ $is_destroy = true ]
  then
      data='{"metadata":{"labels":{"env_status":"destroy_failed"}}}'
  else
      data='{"metadata":{"labels":{"env_status":"provision_failed"}}}'
  fi

  argocd app patch $team_env_name --patch $data --type merge > null

  data='{"metadata":{"labels":{"component_status":"'$component_error_status'"}}}'
  argocd app patch $team_env_config_name --patch $data --type merge >null

  sh /audit.sh $team_name $env_name $config_name "Failed" $component_error_status $reconcile_id $config_reconcile_id $is_destroy 0 "noSkip" $customer_id
}


# This function appends logs to filename provided as an argument
function appendLogs() {
  IFS=''
  while read line; do
    echo $line | tee -a $1 
  done
}

# Base error function that returns from the main process
function Error() {
  if [ -n "$1" ]; then
    echo "Error: "$1
    PatchError
  fi

  exit 1
}


# 'Sub Error Function' this Function should be called instead of Error function as we need to save the logs to s3.
function SaveAndExit() {
  echo $show_output_start
  echo $1 2>&1 | appendLogs /tmp/$s3FileName.txt
  echo $show_output_end
  aws s3 cp /tmp/$s3FileName.txt s3://zlifecycle-$zl_env-tfplan-$customer_id/$team_name/$env_name/$config_name/$config_reconcile_id/$s3FileName --profile compuzest-shared --quiet
  Error $1
}

# A uility function used to return error code
function returnErrorCode() {
  return 99;
}


# This Function is used to set AWS credentials to environment variables.
function setAWSCreds() {
  aws_region=$(aws ssm get-parameter --profile compuzest-shared --region us-east-1 --name "/$1/aws_region" --with-decryption --query "Parameter.Value" | jq -r ".")
  if [ ! -z $aws_region ];
  then
    export AWS_REGION=$aws_region
  fi

  aws_access_key_id=$(aws ssm get-parameter --profile compuzest-shared --region us-east-1 --name "/$1/aws_access_key_id" --with-decryption --query "Parameter.Value" | jq -r ".")
  aws_secret_access_key=$(aws ssm get-parameter --profile compuzest-shared --region us-east-1 --name "/$1/aws_secret_access_key" --with-decryption --query "Parameter.Value" | jq -r ".")

  if [ ! -z $aws_access_key_id -a ! -z $aws_secret_access_key ];
  then
    aws configure set aws_access_key_id $aws_access_key_id 
    aws configure set aws_secret_access_key $aws_secret_access_key

    aws_session_token=$(aws ssm get-parameter --profile compuzest-shared --region us-east-1 --name "/$1/aws_session_token" --with-decryption --query "Parameter.Value" | jq -r ".")

    if [ ! -z $aws_session_token ];
    then
      aws configure set aws_session_token $aws_session_token
    fi

    return 1
  fi
  return 0
}