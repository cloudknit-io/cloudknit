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

team_env_name=$team_name-$env_name
team_env_config_name=$team_name-$env_name-$config_name

# yaml status values
# config --> Initialising... , Success , Failed
# env --> initializing , ended 


url_environment='http://zlifecycle-api.zlifecycle-ui.svc.cluster.local/reconciliation/api/v1/environment/save'
url='http://zlifecycle-api.zlifecycle-ui.svc.cluster.local/reconciliation/api/v1/component/save'
start_date=$(date)
end_date='"'$(date)'"'

is_skipped="false"
if [ $config_reconcile_id -eq 0 ]; then
    end_date=null
    config_reconcile_id=null
else
    if [[ $config_status == "Initialising..." ]]; then
        config_status="provisioning_in_progress"
        if [[ $is_destroy == true ]]; then
            config_status="destroying_in_progress"
        fi
    fi

    if [[ $config_status == "Success" ]]; then
        config_status="Provisioned"
        if [[ $is_destroy = true ]]; then
            config_status="Destroyed"
        fi
    fi

    if [ "$skip_component" != "noSkip" ]; then
        is_skipped="true"
        config_status='skipped'
        if [ "$skip_component" = 'selectiveReconcile' ]; then
            config_status='skipped_reconcile'
        fi
    fi

    if [[ $config_status == *"failed"* ]]; then
        if [[ $is_destroy == true ]]; then
            config_status="destroy_'$config_status'"
        else
            config_status="provision_'$config_status'"
        fi
    fi

fi

if [[ $config_name != 0 && $config_reconcile_id = null ]]; then
    sh ./validate_env_component.sh $team_name $env_name $config_name
    . /argocd/login.sh
    if [[ $config_status == *"skipped"* ]]; then
        data='{"metadata":{"labels":{"is_skipped":"'$is_skipped'"}}}'
    else
        data='{"metadata":{"labels":{"is_skipped":"'$is_skipped'","component_status":"initializing","is_destroy":"'$is_destroy'","audit_status":"initializing"}}}'
    fi
    argocd app patch $team_env_config_name --patch $data --type merge > null
else
    echo -n '0' >/tmp/error_code.txt
fi

component_payload='[{"id" : '$config_reconcile_id', "name" : "'$team_env_config_name'", "status" : "'$config_status'", "startDateTime" : "'$start_date'", "endDateTime" : '$end_date'}]'

end_date='"'$(date)'"'
if [ $reconcile_id -eq 0 ]; then
    end_date=null
    reconcile_id=null
fi

if [ $config_name -eq 0 ]; then
    component_payload=[]
    url=$url_environment
    . /patch_environment.sh
fi

status="provision_'$status'"
if [[ $is_destroy = true ]]; then
    status="destroy_'$status'"
fi

payload='{"reconcileId": '$reconcile_id', "name" : "'$team_env_name'", "teamName" : "'$team_name'", "status" : "'$status'", "startDateTime" : "'$start_date'", "endDateTime" : '$end_date', "componentReconciles" : '$component_payload'}'

echo $payload >reconcile_payload.txt

result=$(curl -X 'POST' "$url" -H 'accept: */*' -H 'Content-Type: application/json' -d @reconcile_payload.txt)
echo $result >/tmp/reconcile_id.txt

