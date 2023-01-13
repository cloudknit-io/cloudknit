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

function PatchProcessPlanError() {
    UpdateComponentStatus "${env_name}" "${team_name}" "${config_name}" "plan_failed"

    sh ../audit.sh $team_name $env_name $config_name "" "plan_failed" $reconcile_id $config_reconcile_id $is_destroy 0 "noSkip" ${customer_id}
    if [ $is_destroy = true ]
    then
        data='{"metadata":{"labels":{"env_status":"destroy_failed"}}}'
    else
        data='{"metadata":{"labels":{"env_status":"provision_failed"}}}'
    fi

    # argocd app patch $team_env_name --patch $data --type merge > null
}

function ProcessPlanError() {
  if [ -n "$1" ];
  then
    echo "ProcessPlanError: "$1
    PatchProcessPlanError
  fi

    exit 1;
}

. /initialize-functions.sh

if [ $result -eq 0 ]
then
    if [ $is_destroy = true ]
    then
        compStatus="destroyed"
    else
        compStatus="provisioned"
    fi
    UpdateComponentStatus "${env_name}" "${team_name}" "${config_name}" "${compStatus}"

elif [ $result -eq 2 ]
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

	UpdateComponentStatus "${env_name}" "${team_name}" "${config_name}" "waiting_for_approval"
    else
	UpdateComponentStatus "${env_name}" "${team_name}" "${config_name}" "initializing_apply"
    fi
fi