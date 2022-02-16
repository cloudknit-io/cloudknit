#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Formatting code from all projects

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

echo "=== formatting ==="
gofumpt -l -w ./app
