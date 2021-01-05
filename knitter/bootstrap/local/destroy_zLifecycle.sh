# Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
#
# Unauthorized copying of this file, via any medium, is strictly prohibited
# Proprietary and confidential
#
# NOTICE: All information contained herein is, and remains the property of
# CompuZest, Inc. The intellectual and technical concepts contained herein are
# proprietary to CompuZest, Inc. and are protected by trade secret or copyright
# law. Dissemination of this information or reproduction of this material is
# strictly forbidden unless prior written permission is obtained from CompuZest, Inc.

kubectl port-forward service/argo-cd-argocd-server 8080:80 -n argocd

kubectl get applications -n argocd -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | xargs kubectl patch applications  -p '{"metadata":{"finalizers":[]}}' --type=merge -n argocd
kubectl patch crd applications.argoproj.io -p '{"metadata":{"finalizers":[]}}' --type=merge -n argocd
argocd app delete $(argocd app list -o name)

argocd cluster rm sandbox
argocd repo rm git@github.com:CompuZest/terraform-environment.git
argocd repo rm git@github.com:CompuZest/helm-charts.git

cd ../../../zlifecycle-provisioner/k8s-addons
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

