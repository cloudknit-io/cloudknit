#!/bin/bash

input=".env.production"
output="static/js/*.js"
remove="static/js/*.js--"
while IFS= read -r line
do
variables=$(echo $line | tr "=" "\n")
    for variable in $variables
    do
        if [[ $variable == __DOCKER_REACT_APP* ]]
        then
            echo $variable="${!variable}"
            sed -i -- "s~${variable}~${!variable}~g" $output
        fi
    done
rm -rf $remove
done < "$input"
