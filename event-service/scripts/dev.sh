#!/bin/bash

## Starts the app, reloading on changes

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

# include env file
if [ -f ".env" ]; then
	export $(cat .env | grep -v "#" | xargs)
fi

export AWS_PROFILE="compuzest-dev"
export DEV_MODE="true"

ulimit -n 2048
air
