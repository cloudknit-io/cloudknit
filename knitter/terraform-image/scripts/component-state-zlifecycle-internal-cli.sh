echo "config name: $config_name"
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
echo $component_status > /tmp/component_status.txt

last_flow_skipped=$(argocd app get $team_env_config_name -o json | jq -r '.metadata.labels.is_skipped') || null

if [[ $component_status == "6" && $is_destroy == true ]] && [[ $config_reconcile_id == null || $last_flow_skipped == "true" ]]; then
    config_status='skipped_destroy'
    is_skipped='true'
fi
