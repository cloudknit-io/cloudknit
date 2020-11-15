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

app="${cust_id_env_name}-terraform-config-${name}"
argocd_server_name=$(kubectl get pods -l app.kubernetes.io/name=argocd-server -n argo --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')
response=`curl --insecure https://argo-cd-argocd-server:443/api/v1/session -d $'{"username":"admin","password":"'$argocd_server_name'"}'`
token=$(echo $response | jq -r '.token')

if [ $result -eq 2 ]
then

  RES=$(curl --insecure -H "Authorization: Bearer ${token}" -H "Content-Type: application/json" -d '{"project":"default","syncPolicy":{"automated":{"prune":true,"selfHeal":true}},"destination":{"namespace":"default","server":"https://192.168.1.155:59999"},"source":{"repoURL":"git@github.com:CompuZest/helm-charts.git","path":"charts/terraform-config","targetRevision":"HEAD","helm":{"values":"is_out_of_sync: true\ncustomer_id: \"'$customer_id'\"\nenv_name: \"'$env_name'\"\nname: \"'$name'\"\nmodule:\n  source: \"'$module_source'\"\n  path: '$module_source_path'\nvariables:\n  - name: cidr\n    value: \"10.1.0.0/16\"\noutputs: \"\"\n"}}}' -X PUT https://argo-cd-argocd-server.argo.svc.cluster.local/api/v1/applications/$app/spec)

  sleep 5

  curl --insecure https://argo-cd-argocd-server.argo.svc.cluster.local/api/v1/applications/$cust_id_env_name/sync -H 'Content-Type: application/json' -H "Authorization: Bearer ${token}" --data-binary '{"revision":"HEAD","prune":false,"dryRun":false,"strategy":{"hook":{}},"resources":null}'

else
  RES=$(curl --insecure -H "Authorization: Bearer ${token}" -H "Content-Type: application/json" -d '{"project":"default","syncPolicy":{"automated":{"prune":true,"selfHeal":true}},"destination":{"namespace":"default","server":"https://192.168.1.155:59999"},"source":{"repoURL":"git@github.com:CompuZest/helm-charts.git","path":"charts/terraform-config","targetRevision":"HEAD","helm":{"values":customer_id: \"'$customer_id'\"\nenv_name: \"'$env_name'\"\nname: \"'$name'\"\nmodule:\n  source: \"'$module_source'\"\n  path: '$module_source_path'\nvariables:\n  - name: cidr\n    value: \"10.1.0.0/16\"\noutputs: \"\"\n"}}}' -X PUT https://argo-cd-argocd-server.argo.svc.cluster.local/api/v1/applications/$app/spec)

fi
else
terraform apply -var-file tfvars/$customer_id/$env_name.tfvars -auto-approve
echo -n 0 > /tmp/plan_code.txt
fi
