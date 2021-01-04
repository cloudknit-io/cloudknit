argocd app delete 1-customer
#argocd cluster rm arn:aws:eks:us-east-1:413422438110:cluster/0-sandbox-eks
argocd repo rm git@github.com:CompuZest/infra-deploy-terraform-config.git
argocd repo rm git@github.com:CompuZest/helm-charts.git

cd ../../../zlifecycle-provisioner/k8s-addons
terraform init
terraform workspace select 0-sandbox
terraform init
terraform destroy -auto-approve -var-file tfvars/sandbox.tfvars

cd ../aws-eks
terraform init
terraform workspace select 0-sandbox
terraform init
terraform destroy -auto-approve -var-file tfvars/sandbox.tfvars


cd ../aws-vpc
terraform workspace select 0-sandbox
terraform destroy -auto-approve -var-file tfvars/sandbox.tfvars
