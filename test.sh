#!/bin/ash

# ((((terraform apply -auto-approve -input=false -no-color terraform-plan || returnErrorCode; echo $? >&3) 2>&1 | appendLogs "/tmp/apply_output.txt" >&4) 3>&1) | (read xs; exit $xs)) 4>&1
# a=$(aws s3 ls 2>&1 >/dev/null)
a=$(echo '"dfghjk,.mnbvc"' | jq -r ".");
# ab= $a | jq -r "."

# echo $?
echo $a
# if [ ! -z $aws_region ]
# then
#     echo $aws_region
# fi