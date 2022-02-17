#!/bin/bash
## Formatting code from all projects

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

echo "=== formatting ==="
gofumpt -l -w ./app
