cd ../../infra-deploy-networking/aws-vpc
terraform workspace select 0-sandbox
terraform apply -auto-approve -var-file tfvars/sandbox.tfvars

cd ../../infra-deploy-platform/aws-eks
terraform init
terraform workspace select 0-sandbox
terraform init
terraform apply -auto-approve -var-file tfvars/sandbox.tfvars
sleep 2m
terraform apply -auto-approve -var-file tfvars/sandbox.tfvars
sleep 2m
aws eks --region us-east-1 update-kubeconfig --name 0-sandbox-eks

cd ../k8s-addons
terraform init
terraform workspace select 0-sandbox
terraform init
terraform apply -auto-approve -var-file tfvars/sandbox.tfvars

cd ../../terraform-k8s-operator
#make generate
#make install
#make docker-build docker-push IMG=shahadarsh/terraform-k8s-operator:latest
make deploy IMG=shahadarsh/terraform-k8s-operator:latest

cd ../environment-operator
make deploy IMG=shahadarsh/environment-operator:latest
