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

echo "Please enter 1 for local and 2 for AWS:"
select LOCATION in "1" "2"; do
    case $LOCATION in
        1 ) ./local/bootstrap_zLifecycle_step1.sh; break;;
        2 ) ./aws/bootstrap_zLifecycle_step1.sh; break;;
    esac
done

cd ../../zlifecycle-il-operator
make deploy IMG=shahadarsh/zlifecycle-il-operator:latest

cd ../zLifecycle/bootstrap

kubectl apply -f common/company-config.yaml

echo ""
echo ""
echo "-------------------------------------"
read -p "Please create secrets and enter Y to continue? " -n 1 -r
echo ""
echo "-------------------------------------"

if [[ $REPLY =~ ^[Yy]$ ]]
then
    ./common/bootstrap_zLifecycle_step2.sh $LOCATION
fi
