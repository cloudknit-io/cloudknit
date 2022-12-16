set -eo pipefail

echo y | argocd login --insecure argocd-zlifecycle-server.zlifecycle-system.svc.cluster.local:443 --grpc-web --username admin --password $ARGOCD_PASSWORD
