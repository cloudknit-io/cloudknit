team_name=$1
env_name=$2
config_name=$3
module_source=$4
module_source_path=$5
variables_file_source=$6
variables_file_path=$7
is_apply=$8
lock_state=$9
team_env_name=$team_name-$env_name
team_env_config_name=$team_name-$env_name-$config_name

cd /home/$module_source_path
mkdir ~/.ssh
cat /root/ssh_secret/id_rsa >> ~/.ssh/id_rsa
chmod 400 ~/.ssh/id_rsa
ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

mkdir ~/.aws
cat <<EOT >> ~/.aws/credentials
[default]
aws_access_key_id = ${CUSTOMER_AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${CUSTOMER_AWS_SECRET_ACCESS_KEY}
[compuzest-shared]
aws_access_key_id = ${SHARED_AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${SHARED_AWS_SECRET_ACCESS_KEY}
EOT

terraform init
terraform workspace select $team_env_name || terraform workspace new $team_env_name
terraform init

if [ $is_apply -eq 0 ]
then
    terraform plan -lock=$lock_state -detailed-exitcode -var-file vars/$variables_file_path
    result=$?
    echo -n $result > /tmp/plan_code.txt

    argocd_server_name=$(kubectl get pods -l app.kubernetes.io/name=argocd-server -n argo --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')

    argocd login --insecure argo-cd-argocd-server:443 --grpc-web --username admin --password $argocd_server_name

    env_sync_status=$(argocd app get $team_env_name -o json | jq '.status.sync.status')
    config_sync_status=$(argocd app get $team_env_config_name -o json | jq '.status.sync.status')

    if [ $result -eq 2 ]
    then
        if [ $config_sync_status != "\"OutOfSync\"" ]
        then
            tfconfig="${team_env_config_name}-terraformconfig"

            argocd app patch-resource $team_env_config_name --kind TerraformConfig --resource-name $tfconfig --patch '{ "spec": { "isInSync": false } }' --patch-type 'application/merge-patch+json'

            if [ $env_sync_status != "\"OutOfSync\"" ]
            then
                argocd app sync $team_env_name
            fi
        fi
    else
        if [ $sync_status == "\"OutOfSync\"" ]
        then
            argocd app sync $team_env_config_name
        fi
    fi
else
    terraform apply -var-file vars/$variables_file_path -auto-approve
    echo -n 0 > /tmp/plan_code.txt
fi
