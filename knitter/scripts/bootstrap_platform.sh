# Check if docker is running
if ! docker info >/dev/null 2>&1; then
    echo "Docker does not seem to be running, run it first and retry"
    exit 1
fi

echo "Please enter 1 for local and 2 for AWS:"
select LOCATION in "1" "2"; do
    case $LOCATION in
        1 ) ./local/bootstrap_zplatform_step1.sh; break;;
        2 ) ./aws/bootstrap_zplatform_step1.sh; break;;
    esac
done

echo ""
echo ""
echo "-------------------------------------"
read -p "Please create secrets and enter Y to continue? " -n 1 -r
echo ""
echo "-------------------------------------"

if [[ $REPLY =~ ^[Yy]$ ]]
then
    ./common/bootstrap_zplatform_step2.sh $LOCATION
fi
