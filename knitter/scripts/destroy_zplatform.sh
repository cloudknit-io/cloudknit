cd ../infra-deploy-platform/k8s-addons
terraform init
terraform workspace select 0-sandbox
terraform init
terraform destroy -auto-approve -var-file tfvars/sandbox.tfvars

cd ../aws-eks
terraform init
terraform workspace select 0-sandbox
terraform init
terraform destroy -auto-approve -var-file tfvars/sandbox.tfvars


cd ../../infra-deploy-networking/aws-vpc
terraform workspace select 0-sandbox
terraform destroy -auto-approve -var-file tfvars/sandbox.tfvars
