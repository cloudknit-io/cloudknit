# exit when any command fails
set -e

workflow_name=$1
team_name=$2
env_name=$3
config_name=$4

team_env_config_name=$team_name-$env_name-$config_name

message="${team_env_config_name} terraform is out of sync. To see the diff & approve the sync to desired state go here: http://localhost:8081/workflows/argo/${workflow_name}"

data='{"channel": "slack-notification","message": "'$message'"}'
echo $data

curl -d "${data}" -H "Content-Type: application/json" -X POST http://webhook-eventsource-svc.argo.svc.cluster.local:12000/terraform-diff

