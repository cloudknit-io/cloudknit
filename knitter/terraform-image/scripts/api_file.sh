source=$1
dest=$2
customer_id=$3

if [[ $source == *"@"* ]]; then
    curl -X 'POST' "http://zlifecycle-api.${customer_id}-system.svc.cluster.local/reconciliation/api/v1/component/putObject" -H 'accept: */*' -H "Content-Type: multipart/form-data" \
    -F 'file='$source'' -F 'path='$dest'' -F 'customerId='$customer_id''
else
    curl -X 'POST' "http://zlifecycle-api.${customer_id}-system.svc.cluster.local/reconciliation/api/v1/component/downloadObject" -H 'accept: */*' -H 'Content-Type: application/json' -d '{"path":"'$source'","customerId":"'$customer_id'"}' > $dest
fi

