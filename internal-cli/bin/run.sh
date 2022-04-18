#!/bin/bash

## Starts the app, reloading on changes

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

# include env file
if [ -f ".env" ]; then
	source .env
fi

make build
build/debug/zlifecycle-internal-cli $@
