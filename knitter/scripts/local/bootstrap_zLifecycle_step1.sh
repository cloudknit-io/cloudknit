echo ""
echo ""
echo "-------------------------------------"
read -p "If you want create a k3d cluster enter Y: " -n 1 -r
echo ""
echo "-------------------------------------"

if [[ $REPLY =~ ^[Yy]$ ]]
then
    if ! docker info >/dev/null 2>&1; then
        echo "Docker does not seem to be running, run it first and retry"
        exit 1
    fi
    k3d cluster create sandbox-k3d -a 3 --api-port 59999
fi

cd ../../infra-deploy-platform/k8s-addons
terraform init
terraform workspace select 0-local || terraform workspace new 0-local
terraform init
terraform apply -auto-approve -var-file tfvars/local.tfvars

cd ../../environment-operator
make deploy IMG=shahadarsh/environment-operator:latest
