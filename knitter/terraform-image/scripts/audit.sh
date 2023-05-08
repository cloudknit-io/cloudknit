team_name=$1
env_name=$2
config_name=$3
status=$4
config_status=$5
reconcile_id=$6
config_reconcile_id=$7
is_destroy=$8
phase=$9
skip_component=${10}
customer_id=${11}

zlifecycle_namespace='zlifecycle-system'
api_version='v1'

team_env_name=$team_name-$env_name
team_env_config_name=$team_name-$env_name-$config_name

. /initialize-functions.sh

# yaml status values
# config --> Initialising... , Success , Failed
# env --> initializing , ended 

echo "config name: $config_name"

start_date=$(date '+%Y-%m-%d %H:%M:%S')
end_date=$(date '+%Y-%m-%d %H:%M:%S')

if [ $config_reconcile_id -eq 0 ]; then
    end_date=null
    config_reconcile_id=null
    if [[ $config_status == "Initialising..." ]]; then
        config_status="provisioning_in_progress"
        if [[ $is_destroy == true ]]; then
            config_status="destroying_in_progress"
        fi
    fi
else
    if [[ $config_status == "Success" ]]; then
        config_status="provisioned"
        if [[ $is_destroy = true ]]; then
            config_status="destroyed"
        fi
    fi

    if [[ $config_status == *"failed"* ]]; then
        if [[ $is_destroy == true ]]; then
            config_status="destroy_"$config_status
        else
            config_status="provision_"$config_status
        fi
    fi
fi

is_skipped="false"
if [ "$skip_component" != "noSkip" ]; then
    is_skipped="true"
fi

if [[ $config_name != 0 && $config_reconcile_id = null ]]; then
    echo "running validate environment component script: team $team_name, environment $env_name, component $config_name"
    comp_status=0
    lastWorkflowRunId=null

    if [[ $is_skipped == false ]]; then
        comp_status="initializing"
        lastWorkflowRunId="initiailizing"
    fi
    echo "Updating component reconcile entry: "
    UpdateComponentReconcile "${team_name}" "${env_name}" "${config_name}" '{ "status" : "'${comp_status}'", "isDestroy" : "'${is_destroy}'", "isSkipped" : '${is_skipped}', "lastWorkflowRunId" : "'${lastWorkflowRunId}'", "startDateTime" : "'"$start_date"'"  }'
fi

echo "write 0 to /tmp/error_code.txt"
echo -n '0' >/tmp/error_code.txt

end_date=$(date '+%Y-%m-%d %H:%M:%S')
if [ $reconcile_id -eq 0 ]; then
    echo "set end_date and reconcile_id to null"
    end_date=null
    reconcile_id=null
fi

# there is no config so must be an environment?
if [ $config_name -eq 0 ]; then
    component_payload=[]
    echo "Updating environment reconcile status"
    UpdateEnvironmentReconcileStatus "${team_name}" "${env_name}"
    reconcileId=$latestEnvReconcileId
else
    if [[ $config_reconcile_id != null ]]; then
        UpdateComponentReconcile "${team_name}" "${env_name}" "${config_name}" '{"status": "'${config_status}'", "endDateTime": "'"$end_date"'"}'
    fi
    reconcileId=$latestCompReconcileId
fi



echo "AUDIT RECONCILE_ID: $reconcileId"
echo $reconcileId > /tmp/reconcile_id.txt

echo "config name: $config_name"
if [[ $config_name != 0 ]]; then
    set +e
    echo "calling zlifecycle-internal-cli state component pull"
    zlifecycle-internal-cli state component pull \
      --company "$customer_id" \
      --team "$team_name" \
      --environment "$env_name" \
      --component "$config_name" \
      -u http://zlifecycle-state-manager."$customer_id"-system.svc.cluster.local:8080 \
      -v

    component_status=$?
    echo "saving component status to /tmp/component_status: component_status $component_status"
    echo $component_status > /tmp/component_status
fi
