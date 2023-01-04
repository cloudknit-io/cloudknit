# this script has access to:
# $team_name
# $env_name

status='0'

echo "Patching environment"

echo "phase is: "$phase;

if [ $phase = '0' ]
then
    if [ $is_destroy = true ]
    then
        status="destroying"
    else
        status="provisioning"
    fi  
fi

if [ $phase = '1' ]
then
    if [ $is_destroy = true ]
    then
        status="destroyed"
    else
        status="provisioned"
    fi    
fi

echo "status is: "$status

if [ $status != '0' ]
then
    # argocd app patch $team_env_name --patch $data --type merge > null
    UpdateEnvironmentStatus "${team_name}" "${env_name}" "${status}"
fi
