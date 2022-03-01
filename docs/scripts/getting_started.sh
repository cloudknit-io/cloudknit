#/bin/bash

company="$1"
team="$2"
env="$3"

echo "Downloading hello world yaml"

curl https://zlifecycle.github.io/docs/examples/first-environment.yaml --output hello-world.yaml

echo "Replacing variables"

sed 's/${company}/'"${company}"'/g' hello-world.yaml | tee hello-world.yaml > /dev/null
sed 's/${team}/'"${team}"'/g' hello-world.yaml | tee hello-world.yaml > /dev/null
sed 's/${env}/'"${env}"'/g' hello-world.yaml | tee hello-world.yaml > /dev/null

echo "hello-world.yaml environment file created"
