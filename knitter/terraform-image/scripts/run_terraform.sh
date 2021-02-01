# Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
#
# Unauthorized copying of this file, via any medium, is strictly prohibited
# Proprietary and confidential
#
# NOTICE: All information contained herein is, and remains the property of
# CompuZest, Inc. The intellectual and technical concepts contained herein are
# proprietary to CompuZest, Inc. and are protected by trade secret or copyright
# law. Dissemination of this information or reproduction of this material is
# strictly forbidden unless prior written permission is obtained from CompuZest, Inc.

team_name=$1
env_name=$2
config_name=$3
module_source=$4
module_source_path=$5
variables_file_source=$6
variables_file_path=$7
is_apply=$8
lock_state=$9
is_sync=$10
team_env_name=$team_name-$env_name
team_env_config_name=$team_name-$env_name-$config_name

cd /home/terraform-config
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

if [ -n "${module_source_path}" ]; then
    full_module_source="${module_source}//${module_source_path}"
else
    full_module_source="${module_source}"
fi

cat > module.tf << EOL
module "${config_name}" {
  source = "${full_module_source}"
EOL

cat vars/$variables_file_path >> module.tf

echo "}" >> module.tf

cat > provider.tf << EOL
provider "aws" {
  region = "us-east-1"
  version = "~> 3.0"
}
EOL

cat > terraform.tf << EOL
terraform {
  required_version = "= 0.13.2"

  backend "s3" {
    profile                 = "compuzest-shared"

    bucket                  = "compuzest-zlifecycle-tfstate"
    key                     = "${team_name}/${env_name}/${config_name}/terraform.tfstate"
    region                  = "us-east-1"

    dynamodb_table          = "compuzest-zlifecycle-tflock"
    encrypt                 = true
  }
}
EOL

cat module.tf
cat provider.tf
cat terraform.tf

terraform init

if [ $is_apply -eq 0 ]
then
    terraform plan -lock=$lock_state -detailed-exitcode
    result=$?
    echo -n $result > /tmp/plan_code.txt

    if [ $is_sync -eq 0 ]
    then
            argoPassword=$(kubectl get secret argocd-server-login -n argocd -o json | jq '.data.password | @base64d' | tr -d '"')

            argocd login --insecure argo-cd-argocd-server:443 --grpc-web --username admin --password $argoPassword

            env_sync_status=$(argocd app get $team_env_name -o json | jq -r '.status.sync.status')
            config_sync_status=$(argocd app get $team_env_config_name -o json | jq -r '.status.sync.status')

            if [ $result -eq 2 ]
            then
                if [ $config_sync_status != "OutOfSync" ]
                then
                    tfconfig="${team_env_config_name}-terraformconfig"

                    argocd app patch-resource $team_env_config_name --kind TerraformConfig --resource-name $tfconfig --patch '{ "spec": { "isInSync": false } }' --patch-type 'application/merge-patch+json'

                    if [ $env_sync_status != "OutOfSync" ]
                    then
                        argocd app sync $team_env_name
                    fi
                fi
            else
                if [ $config_sync_status == "OutOfSync" ]
                then
                    argocd app sync $team_env_config_name
                fi
            fi
    fi
else
    terraform apply -auto-approve
    echo -n 0 > /tmp/plan_code.txt
fi
