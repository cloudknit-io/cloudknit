k3d cluster create sandbox-k3d -a 3

cd ../../../infra-deploy-platform/k8s-addons
terraform init
terraform workspace new 0-local
terraform workspace select 0-local
terraform init
terraform apply -auto-approve -var-file tfvars/local.tfvars

#cd ../../terraform-k8s-operator
#make generate
#make install
#make docker-build docker-push IMG=shahadarsh/terraform-k8s-operator:latest
#make deploy IMG=shahadarsh/terraform-k8s-operator:latest
