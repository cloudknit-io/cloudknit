team_name=$1
env_name=$2
env_component_name=$3
organization_name=$4

team_env_name=$team_name-$env_name
team_env_component_name=$team_name-$env_name-$env_component_name

sh /argocd/login.sh $organization_name

# Check if environment component application exists. If not then skip plan/apply and so that the
# environment component application gets created
argocd app get $team_env_component_name
result=$?
echo -n $result >/tmp/error_code.txt
