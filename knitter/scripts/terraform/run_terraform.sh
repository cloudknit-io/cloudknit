module_source=$1
module_source_path=$2
is_apply=$3
lock_state=$4
customer_id=$5
env_name=$6
name=$7
cust_id_env_name=$customer_id-$env_name

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
terraform workspace new $cust_id_env_name
terraform workspace select $cust_id_env_name
terraform init

if [ $is_apply -eq 0 ]
then
    terraform plan -lock=$lock_state -detailed-exitcode -var-file tfvars/$customer_id/$env_name.tfvars
    result=$?
    echo -n $result > /tmp/plan_code.txt

    if [ $result -eq 2 ]
    then
        tfconfig="${name}-terraformconfig"

        argocd_server_name=$(kubectl get pods -l app.kubernetes.io/name=argocd-server -n argo --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')

        argocd login --insecure argo-cd-argocd-server:443 --grpc-web --username admin --password $argocd_server_name

        argocd app patch-resource $name --kind TerraformConfig --resource-name $tfconfig --patch '{ "spec": { "isInSync": false } }' --patch-type 'application/merge-patch+json'

         argocd app sync $cust_id_env_name

    fi
else
    terraform apply -var-file tfvars/$customer_id/$env_name.tfvars -auto-approve
    echo -n 0 > /tmp/plan_code.txt
fi
