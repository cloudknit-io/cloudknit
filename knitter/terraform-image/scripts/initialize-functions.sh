echo "INITIALIZE FUNCTIONS"
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

# Patches the error status of argocd app
function PatchError() {
  if [ $is_destroy = true ]
  then
      status="destroy_failed"
  else
      status="provision_failed"
  fi

  # argocd app patch $team_env_name --patch $data --type merge > null

  UpdateComponentReconcile "${team_name}" "${env_name}" "${config_name}" '{ "status" : "'${component_error_status}'" }' 
  UpdateEnvironmentReconcileStatus "${team_name}" "${env_name}" "1"

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


# Gets latest component reconcile Id
# Args: 
#   envName
#   teamName
#   compName
latestCompReconcileId=null;
function getLatestCompReconId() {
  local teamName="${1}"
  local envName="${2}"
  local compName="${3}"
  

  local response=$(curl -X 'GET' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/teams/${teamName}/environments/${envName}/components/${compName}" -H 'accept: */*')
  latestCompReconcileId=$(echo $response | jq -r '.latestCompRecon.reconcileId')
}

# Saves Component Status
#   Args:
#     $1 - env name (required)
#     $2 - team name (required)

function UpdateComponentReconcile() {
  local teamName="${1}"
  local envName="${2}"
  local compName="${3}"
  local payload="${4}"

  if [[ $latestCompReconcileId == null ]]; then
    getLatestCompReconId "${teamName}" "${envName}" "${compName}"
  fi

  echo "UpdateComponentReconcile payload: "
  if [[ $payload == *.json ]]; then
    echo $(cat $payload)
    curl -X 'POST' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/reconciliation/component/${latestCompReconcileId}" -H 'accept: */*' -H 'Content-Type: application/json' -d @$payload
  else
    echo $payload>tmp_comp_status.json
    curl -X 'POST' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/reconciliation/component/${latestCompReconcileId}" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_comp_status.json
  fi
}

# Gets the latest env reconcile entry
#   Args:
#     $1 - team name (required)
#     $2 - env name (required)
latestEnvReconcileId=null;
function GetLatestEnvReconId() {
  local teamName="${1}"
  local envName="${2}"
  
  echo "Running GetLatestEnvReconId"
  local response=$(curl -X 'GET' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/teams/${teamName}/environments/${envName}" -H 'accept: */*')
  latestEnvReconcileId=$(echo $response | jq -r '.latestEnvRecon.reconcileId')
}

function UpdateEnvironmentReconcileStatus() {
  local teamName="${1}"
  local envName="${2}"
  local failed="${3}"

  GetLatestEnvReconId "${teamName}" "${envName}"
  end_date=$(date '+%Y-%m-%d %H:%M:%S')
  local status='0'
  local payload=''

  echo "Patching environment"

  echo "phase is: "$phase

  if [[ "$failed" == "1" ]]; then
    if [[ $is_destroy == true ]]; then
      status="destroy_failed"
    else
      status="provision_failed"
    fi
    payload='{"status": "'${status}'", "teamName": "'${team_name}'", "endDateTime": "'${end_date}'"}'
  elif [[ "$phase" == "0" ]]; then
    if [[ $is_destroy == true ]]; then
      status="destroying"
    else
      status="provisioning"
    fi
    payload='{"status": "'${status}'", "teamName": "'${team_name}'"}'
  elif [[ "$phase" == "1" ]]; then
    if [[ $is_destroy == true ]]; then
      status="destroyed"
    else
      status="provisioned"
    fi
    payload='{"status": "'${status}'", "teamName": "'${team_name}'", "endDateTime": "'${end_date}'"}'
  fi

  echo "status is: "$status

  if [[ "$status" != "0" ]]; then

    echo ${payload} >tmp_update_env_recon.json

    echo "PAYLOAD: $payload"
    result=$(curl -X 'POST' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/reconciliation/environment/${latestEnvReconcileId}" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_update_env_recon.json)
  fi
}
