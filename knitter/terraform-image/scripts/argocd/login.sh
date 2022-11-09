set -eo pipefail

organization_name=$1

echo y | argocd login --insecure argocd-"$organization_name"-server."$organization_name"-system.svc.cluster.local:443 --grpc-web --username admin --password $ARGOCD_PASSWORD
