LOCATION=$1

cd ../../zlifecycle-provisioner/k8s-addons/argo-workflow

kubectl port-forward service/argo-cd-argocd-server 8080:80 -n argocd &

sleep 2m
argocd_server_name=$(kubectl get pods -l app.kubernetes.io/name=argocd-server -n argocd --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')
argocd login --insecure localhost:8080 --grpc-web --username admin --password $argocd_server_name

sleep 10s
ilRepo=$(kubectl get ConfigMap company-config -n zlifecycle-il-operator-system -o jsonpath='{.data.ilRepo}')
ilRepoName=$(kubectl get ConfigMap company-config -n zlifecycle-il-operator-system -o jsonpath='{.data.ilRepoName}')
argocd repo add --name $ilRepoName $ilRepo --ssh-private-key-path zLifecycle --insecure-ignore-host-key

sleep 10s
helmChartsRepo=$(kubectl get ConfigMap company-config -n zlifecycle-il-operator-system -o jsonpath='{.data.helmChartsRepo}')
argocd repo add --name helm-charts $helmChartsRepo --ssh-private-key-path zLifecycle --insecure-ignore-host-key

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
    kubectl create secret generic k8s-api --from-literal=url=$APISERVER -n zlifecycle-il-operator-system

    argocd cluster add k3d-sandbox-k3d --insecure --name sandbox

else 

    APISERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
    kubectl create secret generic k8s-api --from-literal=url=$APISERVER -n zlifecycle-il-operator-system

    argocd cluster add arn:aws:eks:us-east-1:413422438110:cluster/0-sandbox-eks --name sandbox

fi

# Create all bootstrap argo workflow template
cd ../../../zLifecycle/argo-templates
kubectl apply -f .

# Create all team environments
cd ../../compuzest-zlifecycle-config
kubectl apply -R -f teams/account-team
kubectl apply -R -f teams/user-team

kubectl port-forward service/argo-workflow-server 8081:2746 -n argocd &
