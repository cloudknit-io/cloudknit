#!/bin/bash

function notifySlack() {
  if [ -z "$SLACK_HOOK" ]; then
    return 0
  fi

  #curl -s --retry 3 --retry-delay 3 -X POST --data-urlencode 'payload={"text": "'"$1"'"}' $SLACK_HOOK > /dev/null
}

function rotateNodes() {
  asgName=$1
  asgRegion="us-east-1"

  echo "get nodes"

  # Get number of nodes older than MAX_AGE_DAYS in the current ASG
  oldNodes=$(kubectl get nodes 2> /dev/null | \
    sed 's/d//g' | \
    wc -l)

  if [[ $oldNodes != "" && $oldNodes -gt 0 ]]; then
    currentAsgNodes=$(aws autoscaling describe-auto-scaling-groups \
      --auto-scaling-group-name $asgName --region $asgRegion | \
      jq '.AutoScalingGroups[].DesiredCapacity')

    if [[ $currentAsgNodes != "" ]]; then
      desiredNodes=$(expr $currentAsgNodes + $oldNodes 2> /dev/null)

      if [[ $desiredNodes != "" && $desiredNodes -gt 0 ]]; then
        aws autoscaling set-desired-capacity --auto-scaling-group-name $asgName \
          --desired-capacity $desiredNodes --region $asgRegion

        if [[ $? -eq 0 ]]; then
          echo "`date` -- Found $oldNodes nodes older than $MAX_AGE_DAYS days in $asgName. Scaled up $oldNodes and waiting for scale down..."
          notifySlack "Found $oldNodes nodes older than $MAX_AGE_DAYS days in $asgName. Scaled up $oldNodes and waiting for scale down..."
        else
          echo "`date` -- Found $oldNodes nodes older than $MAX_AGE_DAYS days in $asgName. Failed to scale up for nodes rotation, hit maximum."
          notifySlack "Found $oldNodes nodes older than $MAX_AGE_DAYS days in $asgName. Failed to scale up for nodes rotation, hit maximum."
        fi
      fi
    fi
  else
    echo "`date` -- No old nodes found in $asgName."
  fi

  return 0
}

echo "start"
autoscalingGroupsNoWs=$(echo "$AUTOSCALING_GROUPS" | tr -d "[:space:]")
IFS=';' read -ra autoscalingGroups <<< "$autoscalingGroupsNoWs"

for asg in "${autoscalingGroups[@]}"; do
  IFS='|' read asgName asgRegion <<< "$asg"

  rotateNodes $asgName $asgRegion
done
