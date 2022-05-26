#!/usr/bin/env bash

set -e;

COMPANY=$1


LOC=$(realpath $(dirname $0))
OUT_FILE=$(realpath $COMPANY.env)

mkdir -p $(realpath ../certs)
CERTS_LOC=$(realpath ../certs/${COMPANY})
PODNAME=$(kubectl get pods -n $COMPANY-system -l "app.kubernetes.io/instance=zlifecycle-il-operator" -o jsonpath="{.items[0].metadata.name}")
ENV=$(kubectl exec -n $COMPANY-system -it $PODNAME -- env)
REQUIRED_ENV=$(cat $LOC/env_var_names)

# pre-run cleanup
rm -f $OUT_FILE
rm -Rf $CERTS_LOC

# set these to false
MAKE_FALSE="ENABLE_NEW_RELIC KUBERNETES_DISABLE_WEBHOOKS"

echo "Operating on pod : $PODNAME"

for cluster_var in $ENV
do
    readarray -d = -t varname <<< "$cluster_var"

    CONTINUE=false
    for f in $MAKE_FALSE
    do
        if [ $varname == $f ]; then
            echo "${varname}=false" >> $OUT_FILE
            CONTINUE=true
            break
        fi
    done

    if $CONTINUE; then
        continue
    fi

    for req in $REQUIRED_ENV
    do
        if [[ "${varname[0]}" == $req ]]; then
            echo "$cluster_var" >> $OUT_FILE
            break
        fi
    done
done

echo "KUBERNETES_CERT_DIR=$CERTS_LOC" >> $OUT_FILE
echo "MODE=local" >> $OUT_FILE

echo "Created $OUT_FILE"

mkdir -p "${CERTS_LOC}"
echo "$(kubectl exec -n $COMPANY-system -it $PODNAME -- cat /tmp/k8s-webhook-server/serving-certs/ca.crt)" > "$CERTS_LOC/ca.crt"
echo "Created ca.crt"
echo "$(kubectl exec -n $COMPANY-system -it $PODNAME -- cat /tmp/k8s-webhook-server/serving-certs/tls.crt)" > "$CERTS_LOC/tls.crt"
echo "Created tls.crt"
echo "$(kubectl exec -n $COMPANY-system -it $PODNAME -- cat /tmp/k8s-webhook-server/serving-certs/tls.key)" > "$CERTS_LOC/tls.key"
echo "Created tls.key"
