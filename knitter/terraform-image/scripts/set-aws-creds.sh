setAWSCreds $customer_id/$team_name/$env_name
aws_response=$?
if [ $aws_response -eq 0 ];
then
  setAWSCreds $customer_id/$team_name
  aws_response=$?
  if [ $aws_response -eq 0 ];
  then
    setAWSCreds $customer_id
    aws_response=$?
    if [ $aws_response -eq 0 ];
    then
      SaveAndExit "No AWS Credentials available. Please set AWS Credentials in the Settings Page."
    fi
  fi
fi