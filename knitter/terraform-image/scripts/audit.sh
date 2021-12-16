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

url_environment='http://zlifecycle-api.zlifecycle-ui.svc.cluster.local/reconciliation/api/v1/environment/save'
url='http://zlifecycle-api.zlifecycle-ui.svc.cluster.local/reconciliation/api/v1/component/save'
start_date=$(date)
end_date='"'$(date)'"'


if [ $config_reconcile_id -eq 0 ]; then
    end_date=null
    config_reconcile_id=null
fi

if [ "$skip_component" != "noSkip" ]; then
    if [[ "$config_status" = "Initialising..." && $config_reconcile_id = null ]]; then
        . /argocd/login.sh
        data='{"metadata":{"labels":{"is_skipped":"true"}}}'
        argocd app patch $team_env_config_name --patch $data --type merge > null
    fi
    config_status='skipped'
    if [ "$skip_component" = 'selectiveReconcile' ]; then
        config_status='skipped_reconcile'
    fi
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

payload='{"reconcileId": '$reconcile_id', "name" : "'$team_env_name'", "teamName" : "'$team_name'", "status" : "'$status'", "startDateTime" : "'$start_date'", "endDateTime" : '$end_date', "componentReconciles" : '$component_payload'}'

echo $payload >reconcile_payload.txt

result=$(curl -X 'POST' "$url" -H 'accept: */*' -H 'Content-Type: application/json' -d @reconcile_payload.txt)
echo $result >/tmp/reconcile_id.txt

