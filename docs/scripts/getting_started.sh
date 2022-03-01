#/bin/bash

echo Enter your company name used by zLifecycle:

read company

echo Enter the team name for which you want to provision the environment:

read team

echo Enter name of the environment you want to provision:

read env

echo "Downloading hello world template yaml"

curl https://zlifecycle.github.io/docs/examples/first-environment.yaml --output hello-world.yaml

cat hello-world.yaml

echo "Replacing variables"

sed 's/${company}/'"${company}"'/g' hello-world.yaml | tee hello-world.yaml > /dev/null
sed 's/${team}/'"${team}"'/g' hello-world.yaml | tee hello-world.yaml > /dev/null
sed 's/${env}/'"${env}"'/g' hello-world.yaml | tee hello-world.yaml > /dev/null

echo "hello-world.yaml environment file created"
