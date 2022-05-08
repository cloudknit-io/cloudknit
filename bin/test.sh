#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Runs all the tests

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

if [ -f "test.env" ]; then
	export $(cat test.env | grep -v "#" | xargs)
fi

if [ -f "./bin/test-setup.sh" ]; then
	./bin/test-setup.sh
fi

gotestsum ./controller/...
