#!/bin/bash

## Starts the app, reloading on changes

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

envfile="$1.env"

# include env file
if [ -f $envfile ]; then
	echo "found $1.env file, sourcing it..."
	source $envfile
else
	echo "environment not selected!"
	return 1
fi

echo "running operator in $MODE mode"

go run main.go
