echo ""
echo ""
echo "-------------------------------------"
read -p "If you want create a k3d cluster enter Y: " -n 1 -r
echo ""
echo "-------------------------------------"

if [[ $REPLY =~ ^[Yy]$ ]]
then
    k3d cluster create sandbox-k3d -a 3 --api-port 59999
fi


cd ../../infra-deploy-platform/k8s-addons
terraform init
terraform workspace select 0-local || terraform workspace new 0-local
terraform init
terraform apply -auto-approve -var-file tfvars/local.tfvars

cd ../../terraform-k8s-operator
#make generate
#make install
#make docker-build docker-push IMG=shahadarsh/terraform-k8s-operator:latest
make deploy IMG=shahadarsh/terraform-k8s-operator:latest

cd ../environment-operator
make deploy IMG=shahadarsh/environment-operator:latest
