AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
REGION=us-east-1
SECRET_NAME=${REGION}-ecr-registry
EMAIL=zLifecycle@compuzest.com

TOKEN=`aws ecr --region=$REGION get-authorization-token --output text --query authorizationData[].authorizationToken | base64 -d | cut -d: -f2`

kubectl delete secret --ignore-not-found $SECRET_NAME
kubectl create secret -n zlifecycle-il-operator-system docker-registry $SECRET_NAME \
 --docker-server=https://$AWS_ACCOUNT_ID.dkr.ecr.${REGION}.amazonaws.com \
 --docker-username=AWS \
 --docker-password="${TOKEN}" \
 --docker-email="${EMAIL}"

kubectl create secret -n zlifecycle-ui docker-registry $SECRET_NAME \
 --docker-server=https://$AWS_ACCOUNT_ID.dkr.ecr.${REGION}.amazonaws.com \
 --docker-username=AWS \
 --docker-password="${TOKEN}" \
 --docker-email="${EMAIL}"

kubectl create secret -n argocd docker-registry $SECRET_NAME \
 --docker-server=https://$AWS_ACCOUNT_ID.dkr.ecr.${REGION}.amazonaws.com \
 --docker-username=AWS \
 --docker-password="${TOKEN}" \
 --docker-email="${EMAIL}"
