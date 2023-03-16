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

  UpdateComponentStatus "${env_name}" "${team_name}" "${config_name}" "${component_error_status}"
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

# Saves Component Status
#   Args:
#     $1 - env name (required)
#     $2 - team name (required)
#     $3 - component name (required)
#     $4 - component status (required)
#     $5 - component isDestroyed (default: false)
function UpdateComponentStatus() {
  local envName="${1}"
  local teamName="${2}"
  local compName="${3}"
  local compStatus="${4}"

  local payload='{ "status" : "'${compStatus}'" }'
  
  echo "Running UpdateComponentStatus ${compStatus} : ${payload}"
  echo $payload >tmp_comp_status.json

  curl -X 'PUT' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/teams/${teamName}/environments/${envName}/components/${compName}" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_comp_status.json
}

# Saves Component isDestroyed
#   Args:
#     $1 - env name (required)
#     $2 - team name (required)
#     $3 - component name (required)
#     $4 - component isDestroyed (required)
function UpdateComponentDestroyed() {
  local envName="${1}"
  local teamName="${2}"
  local compName="${3}"
  local isDestroyed=${4}

  local payload='{ "isDestroyed" : '${isDestroyed}' }'

  echo "Running UpdateComponentDestroyed: ${payload}"
  echo $payload >tmp_comp_status.json

  curl -X 'PUT' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/teams/${teamName}/environments/${envName}/components/${compName}" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_comp_status.json
}

# Saves or updates a component
#   Args:
#     $1 - env name (required)
#     $2 - team name (required)
#     $3 - component name (required)
#     $4 - workflow run id (required)
function UpdateComponentWfRunId() {
  local envName="${1}"
  local teamName="${2}"
  local compName="${3}"
  local wfRunId="${4}"

  local payload='{ "lastWorkflowRunId" : "'${wfRunId}'" }'
  
  echo "Running UpdateComponentWfRunId ${wfRunId} : ${payload}"
  echo $payload >tmp_comp_wf_runid.json

  curl -X 'PUT' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/teams/${teamName}/environments/${envName}/components/${compName}" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_comp_wf_runid.json
}

# Saves or updates a component
#   Args:
#     $1 - env name (required)
#     $2 - team name (required)
#     $3 - component name (required)
#     $4 - cost (required)
#     $5 - cost resources (required)
function UpdateComponentCost() {
  local envName="${1}"
  local teamName="${2}"
  local compName="${3}"
  local estimatedCost="${4}"
  local costResources="${5}"

  local payload='{ "estimatedCost" : '${estimatedCost}', "costResources" : '${costResources}' }'
  
  echo "Running UpdateComponentCost : ${payload}"
  echo $payload >tmp_comp_cost.json

  curl -X 'PUT' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/teams/${teamName}/environments/${envName}/components/${compName}" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_comp_cost.json
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
  local response=$(curl -X 'GET' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/teams/${teamName}/environments/${envName}/audit/latest" -H 'accept: */*')
  latestEnvReconcileId=$(echo $response | jq -r '.reconcileId')
}

function UpdateEnvironmentReconcileStatus() {
  local teamName="${1}"
  local envName="${2}"
  local failed="${3}"

  GetLatestEnvReconId "${teamName}" "${envName}"
  
  local status='0'

  echo "Patching environment"

  echo "phase is: "$phase;

  if [ $failed = '1']
  then
    if [ $is_destroy = true ]
    then
        status="destroy_failed"
    else
        status="provision_failed"
    fi
  else
    if [ $phase = '0' ]
    then
        if [ $is_destroy = true ]
        then
            status="destroying"
        else
            status="provisioning"
        fi  
    fi

    if [ $phase = '1' ]
    then
        if [ $is_destroy = true ]
        then
            status="destroyed"
        else
            status="provisioned"
        fi    
    fi
  fi

  echo "status is: "$status

  if [ $status != '0' ]
  then
    payload='{"status": "'${status}'", "teamName": "'${team_name}'", "endDateTime": "'${end_date}'"}'
    echo ${payload} >tmp_update_env_recon.json

    echo "PAYLOAD: $payload"
    result=$(curl -X 'POST' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/reconciliation/environment/${latestEnvReconcileId}" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_update_env_recon.json)
  fi
}
