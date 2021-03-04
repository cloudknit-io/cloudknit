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

team_name=$1
env_name=$2
env_component_name=$3

team_env_name=$team_name-$env_name
team_env_component_name=$team_name-$env_name-$env_component_name

argoPassword=$(kubectl get secret argocd-server-login -n argocd -o json | jq '.data.password | @base64d' | tr -d '"')

echo y | argocd login --insecure argocd-server:443 --grpc-web --username admin --password $argoPassword

# Check if environment component application exists. If not then skip plan/apply and so that the 
# environment component application gets created
argocd app get $team_env_component_name 
result=$?
echo -n $result > /tmp/error_code.txt
