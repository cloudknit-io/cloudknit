#!/bin/bash
## Runs all the tests

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

if [ -f "test.env" ]; then
	export $(cat test.env | grep -v "#" | xargs)
fi

if [ -f "./bin/test-setup.sh" ]; then
	./scripts/test-setup.sh
fi

gotestsum ./internal/...
