#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Runs all the tests

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

export ZLI_TEST_MODE=$1

if [ -f "test.env" ]; then
	echo "found test.env file, sourcing it..."
	source test.env
fi

if [ -f "./bin/test-setup.sh" ]; then
	echo "found test-setup.sh file, executing it..."
	./bin/test-setup.sh
fi

gotestsum ./app/...
