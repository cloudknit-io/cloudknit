#!/bin/bash

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

go run ./cmd/zlifecycle-event-service/main.go --migrate up
