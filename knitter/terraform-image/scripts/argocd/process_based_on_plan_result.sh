is_sync=$1
result=$2
team_env_name=$3
team_env_config_name=$4
workflow_id=$5
is_destroy=$6
team_name=$7
env_name=$8
config_name=$9
reconcile_id=${10}
config_reconcile_id=${11}
auto_approve=${12}
customer_id=${13}

echo "PROCESS BASED ON PLAN RESULT"
echo "   is_sync=${is_sync}"
echo "   result=${result}"
echo "   team_env_name=${team_env_name}"
echo "   team_env_config_name=${team_env_config_name}"
echo "   workflow_id=${workflow_id}"
echo "   is_destroy=${is_destroy}"
echo "   team_name=${team_name}"
echo "   env_name=${env_name}"
echo "   config_name=${config_name}"
echo "   reconcile_id=${reconcile_id}"
echo "   config_reconcile_id=${config_reconcile_id}"
echo "   auto_approve=${auto_approve}"
echo "   customer_id=${customer_id}"
echo ""

function PatchError() {
    UpdateComponentStatus $env_name $team_name $config_name "plan_failed"
    # TODO : Pass orgId
    sh ../audit.sh $team_name $env_name $config_name "" "plan_failed" $reconcile_id $config_reconcile_id $is_destroy 0 "noSkip" ${customer_id}
    if [ $is_destroy = true ]
    then
        data='{"metadata":{"labels":{"env_status":"destroy_failed"}}}'
    else
        data='{"metadata":{"labels":{"env_status":"provision_failed"}}}'
    fi
    argocd app patch $team_env_name --patch $data --type merge > null
}

function Error() {
  if [ -n "$1" ];
  then
    echo "Error: "$1
    PatchError
  fi

    exit 1;
}

env_sync_status=$(argocd app get $team_env_name -o json | jq -r '.status.sync.status') || Error "Failed getting env_sync_status"

config_sync_status=$(argocd app get $team_env_config_name -o json | jq -r '.status.sync.status') || Error "Failed getting config_sync_status"

if [ $result -eq 0 ]
then
    if [ $is_destroy = true ]
    then
        data="destroyed"
    else
        data="provisioned"
    fi
    UpdateComponentStatus $env_name $team_name $config_name ${data}

    if [ $config_sync_status == "OutOfSync" ]
    then
        argocd app sync $team_env_config_name > null || true
    fi
elif [ $result -eq 2 ]
then
    if [ $is_sync -eq 0 ]
    then
        if [ $config_sync_status != "OutOfSync" ]
        then
            UpdateComponentStatus $env_name $team_name $config_name "out_of_sync"

            if [ $env_sync_status != "OutOfSync" ]
            then
                argocd app sync $team_env_name > null || true
            fi
        fi
    elif [ $is_sync -eq 1 ]
    then
        if [ $auto_approve = false ]
        then
            zlifecycle-internal-cli state component patch \
              --company $customer_id \
              --team $team_name \
              --environment $env_name \
              --component $config_name \
              --status waiting_for_approval \
              -u http://zlifecycle-state-manager."${customer_id}"-system.svc.cluster.local:8080 \
              -v

            UpdateComponentStatus $env_name $team_name $config_name "waiting_for_approval"
        fi
    fi
fi
