#/bin/bash

echo "Enter your company name used by zLifecycle (E.g. Enter "zmart" if your zL url is https://zmart.app.zlifecycle.com):"

read company

if [[ $company == "" ]]
then
    echo "Company Name can't be empty. Please try again!"

    exit 0
fi

echo "Enter the team name (E.g. checkout):"

read team

if [[ $team == "" ]]
then
    echo "Team Name can't be empty. Please try again!"

    exit 0
fi

echo "Enter the environment name (E.g. dev):"

read env

if [[ $env == "" ]]
then
    echo "Environment Name can't be empty. Please try again!"

    exit 0
fi

echo "Downloading hello world template yaml"

curl https://docs.zlifecycle.com/examples/first-environment.yaml --output hello-world.yaml

echo "Replacing variables"

cp hello-world.yaml hello-world.yaml.tmp

sed -e 's/${company}/'"${company}"'/g' -e 's/${team}/'"${team}"'/g'  -e 's/${env}/'"${env}"'/g' hello-world.yaml.tmp > hello-world.yaml

rm hello-world.yaml.tmp

echo "hello-world.yaml environment file created"
