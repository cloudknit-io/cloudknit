kubectl port-forward service/argo-cd-argocd-server 8080:80 -n argocd

kubectl get applications -n argocd -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | xargs kubectl patch applications  -p '{"metadata":{"finalizers":[]}}' --type=merge -n argocd
kubectl patch crd applications.argoproj.io -p '{"metadata":{"finalizers":[]}}' --type=merge -n argocd
argocd app delete $(argocd app list -o name)

argocd cluster rm sandbox
argocd repo rm git@github.com:CompuZest/terraform-environment.git
argocd repo rm git@github.com:CompuZest/helm-charts.git

cd ../../../infra-deploy-zlifecycle/k8s-addons
terraform init
terraform workspace select 0-local
terraform init

terraform destroy -auto-approve -var-file tfvars/sandbox.tfvars

echo ""
echo ""
echo "-------------------------------------"
read -p "If you want delete the k3d cluster enter Y: " -n 1 -r
echo ""
echo "-------------------------------------"

if [[ $REPLY =~ ^[Yy]$ ]]
then
    k3d cluster delete sandbox-k3d
fi

