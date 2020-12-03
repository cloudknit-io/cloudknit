LOCATION=$1

cd ../../infra-deploy-platform/k8s-addons/argo-workflow

argocd_server_name=$(kubectl get pods -l app.kubernetes.io/name=argocd-server -n argo --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')
kubectl port-forward service/argo-cd-argocd-server 8080:80 -n argo &

sleep 2m

argocd login --insecure localhost:8080 --grpc-web --username admin --password $argocd_server_name

sleep 10s
argocd repo add --name terraform-environment git@github.com:CompuZest/terraform-environment.git --ssh-private-key-path argo --insecure-ignore-host-key

sleep 10s
argocd repo add --name helm-charts git@github.com:CompuZest/helm-charts.git --ssh-private-key-path argo --insecure-ignore-host-key


if [ $LOCATION -eq 1 ]
then
    ip_addr=$(ipconfig getifaddr en0)

    if [ ! $ip_addr ]
    then
        ip_addr=$(ipconfig getifaddr en1)
    fi

    sed -i .bak "s+https://0.0.0.0:59999+https://$ip_addr:59999+g" ~/.kube/config

    sleep 10s

    curl --insecure https://$ip_addr:59999

    sleep 10s

    APISERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
    kubectl create secret generic k8s-api --from-literal=url=$APISERVER -n environment-operator-system

    argocd cluster add k3d-sandbox-k3d --insecure --name sandbox

else 

    APISERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
    kubectl create secret generic k8s-api --from-literal=url=$APISERVER -n environment-operator-system

    argocd cluster add arn:aws:eks:us-east-1:413422438110:cluster/0-sandbox-eks --name sandbox
fi

# Create all bootstrap argo workflow template
cd ../../../infra-deploy-bootstrap/argo-templates
kubectl apply -f .

# Create all team environments
cd ../../infra-deploy-terraform-config
kubectl apply -R -f teams/account-team
kubectl apply -R -f teams/user-team

# kubectl apply -f teams/account-team.yaml
# kubectl apply -f teams/user-team.yaml

kubectl port-forward service/argo-workflow-server 8081:2746 -n argo &

