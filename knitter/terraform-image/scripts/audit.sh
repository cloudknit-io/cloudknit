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
        config_status="Provisioned"
        if [[ $is_destroy = true ]]; then
            config_status="Destroyed"
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
    config_status='skipped_destroy'
    if [ "$skip_component" = 'selectiveReconcile' ]; then
        config_status='skipped_provision'
    fi
fi

#echo "running argocd login script"
#sh /argocd/login.sh $customer_id

# TODO: Look at this block to see if we need to replicate it
#echo "current config status: $config_status"
#if [[ $config_name != 0 ]]; then
#    echo "Fetching component status via zlifecycle-internal-cli"
    #. /component-state-zlifecycle-internal-cli.sh
#fi

if [[ $config_name != 0 && $config_reconcile_id = null ]]; then
    echo "running validate environment component script: team $team_name, environment $env_name, component $config_name"
    # sh ./validate_env_component.sh $team_name $env_name $config_name $customer_id
    comp_status=0
    if [[ $config_status == *"skipped"* ]]; then
        echo "getting environment component previous status"
        config_previous_status=$(curl "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/costing/component?teamName=${team_name}&envName=${env_name}&compName=${config_name}" | jq -r ".status") || null
        echo "config_prev_status: $config_previous_status"
        if [[ $config_previous_status == null ]]; then
            comp_status="not_provisioned"
            data='{"metadata":{"labels":{"is_skipped":"'$is_skipped'","component_status":"not_provisioned"}}}'
        else
            data='{"metadata":{"labels":{"is_skipped":"'$is_skipped'"}}}'
        fi
    else
        comp_status="initializing"
        UpdateComponentWfRunId "${env_name}" "${team_name}" "${config_name}" "initializing"
        data='{"metadata":{"labels":{"is_skipped":"'$is_skipped'","audit_status":"initializing","last_workflow_run_id":"initializing"}}}'
    fi
    # echo "patch argocd resource $team_env_config_name with data $data"
    # argocd app patch $team_env_config_name --patch $data --type merge > null
    if [[ $comp_status != 0 ]]; then
        UpdateComponentStatus "${env_name}" "${team_name}" "${config_name}" "${comp_status}" ${is_destroy}
    fi
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
    echo "calling patch environment script"
    . /patch_environment.sh
fi

if [[ $is_destroy = true ]]; then
    status="destroy_"$status
else
    status="provision_"$status
fi

result=""

if [ $config_name -eq 0 ]; then # environment recon
    if [ $reconcile_id = null ]; then # create env reconcile
        payload='{"name": "'${env_name}'", "teamName": "'${team_name}'", "startDateTime": "'${start_date}'"}'
        echo ${payload} >tmp_new_env_recon.json

	echo "PAYLOAD: $payload"
        result=$(curl -X 'POST' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/reconciliation/environment" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_new_env_recon.json)
    else # update env reconcile
        payload='{"status": "'${status}'", "teamName": "'${team_name}'", "endDateTime": "'${end_date}'"}'
        echo ${payload} >tmp_update_env_recon.json

	echo "PAYLOAD: $payload"
        result=$(curl -X 'POST' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/reconciliation/environment/${reconcile_id}" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_update_env_recon.json)
    fi
else # component recon
    if [ $config_reconcile_id = null ]; then # create comp reconcile
        payload='{"name": "'${config_name}'", "startDateTime": "'${start_date}'", "envReconcileId": "'${reconcile_id}'"}'
        echo ${payload} >tmp_new_comp_recon.json

	echo "PAYLOAD: $payload"
        result=$(curl -X 'POST' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/reconciliation/component" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_new_comp_recon.json)
    else # update comp reconcile
        payload='{"status": "'${status}'", "endDateTime": "'${end_date}'"}'
        echo ${payload} >tmp_update_comp_recon.json

	echo "PAYLOAD: $payload"
        result=$(curl -X 'POST' "http://zlifecycle-api.zlifecycle-system.svc.cluster.local/v1/orgs/${customer_id}/reconciliation/component/${config_reconcile_id}" -H 'accept: */*' -H 'Content-Type: application/json' -d @tmp_update_comp_recon.json)
    fi
fi

echo "AUDIT RECONCILE_ID: $result"
echo $result > /tmp/reconcile_id.txt

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
