team_env_name=$team_name-$env_name

data='0'

echo "Patching environment"

echo "phase is: "$phase;

if [ $phase = '0' ]
then
    if [ $is_destroy = true ]
    then
        data='{"metadata":{"labels":{"env_status":"destroying"}}}'    
    else
        data='{"metadata":{"labels":{"env_status":"provisioning"}}}'
    fi  
fi

if [ $phase = '1' ]
then
    if [ $is_destroy = true ]
    then
        data='{"metadata":{"labels":{"env_status":"destroyed"}}}'    
    else
        data='{"metadata":{"labels":{"env_status":"provisioned"}}}'
    fi    
fi

echo "data is: "$data

if [ $data != '0' ]
then
    argocd app patch $team_env_name --patch $data --type merge > null
fi
