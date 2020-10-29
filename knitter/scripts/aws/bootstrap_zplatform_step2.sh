cd ../../../infra-deploy-platform/k8s-addons/argo-workflow

argocd_server_name=$(kubectl get pods -l app.kubernetes.io/name=argocd-server -n argo --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')
kubectl port-forward service/argo-cd-argocd-server 8080:80 -n argo &

sleep 2m

argocd login --insecure localhost:8080 --grpc-web --username admin --password $argocd_server_name

argocd repo add --name infra-deploy-terraform-config git@github.com:CompuZest/infra-deploy-terraform-config.git --ssh-private-key-path argo --insecure-ignore-host-key

argocd repo add --name helm-charts git@github.com:CompuZest/helm-charts.git --ssh-private-key-path argo --insecure-ignore-host-key

#argocd cluster add arn:aws:eks:us-east-1:413422438110:cluster/0-sandbox-eks
argocd cluster add k3d-sandbox-k3d

cd ../../../infra-deploy-bootstrap/scripts

kubectl apply -f terraform-sync-template.yaml
kubectl apply -f terraform-template.yaml

argocd app create 1-customer --repo git@github.com:CompuZest/infra-deploy-terraform-config.git --path 1 --dest-server https://kubernetes.default.svc --dest-namespace default --sync-policy automated --auto-prune

kubectl port-forward service/argo-workflow-server 8081:2746 -n argo &

