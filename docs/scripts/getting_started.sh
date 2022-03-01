#/bin/bash

echo Enter your company name used by zLifecycle:

read company

echo Enter the team name:

read team

echo Enter the environment name:

read env

echo "Downloading hello world template yaml"

curl https://zlifecycle.github.io/docs/examples/first-environment.yaml --output hello-world.yaml

echo "Replacing variables"

sed -e 's/${company}/'"${company}"'/g' -e 's/${team}/'"${team}"'/g'  -e 's/${env}/'"${env}"'/g' hello-world.yaml | tee hello-world.yaml > /dev/null

echo "hello-world.yaml environment file created"
