argocd app delete 1-customer
argocd cluster rm sandbox
argocd repo rm git@github.com:CompuZest/terraform-environment.git
argocd repo rm git@github.com:CompuZest/helm-charts.git

cd ../../../infra-deploy-platform/k8s-addons
terraform init
terraform workspace select 0-local
terraform init

terraform destroy -auto-approve -var-file tfvars/sandbox.tfvars

k3d cluster delete sandbox-k3d
