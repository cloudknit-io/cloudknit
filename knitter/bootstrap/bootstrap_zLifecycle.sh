echo "Please enter 1 for local and 2 for AWS:"
select LOCATION in "1" "2"; do
    case $LOCATION in
        1 ) ./local/bootstrap_zLifecycle_step1.sh; break;;
        2 ) ./aws/bootstrap_zLifecycle_step1.sh; break;;
    esac
done

cd ../../environment-operator
make deploy IMG=shahadarsh/environment-operator:latest

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
