#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Runs code statistics, checks for outdated dependencies, then runs linters

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

echo "=== linting ==="
golangci-lint run --fix --max-issues-per-linter=0 --sort-results ./...
