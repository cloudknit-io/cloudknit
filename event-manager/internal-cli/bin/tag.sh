#!/bin/bash

## Tags the git repo using the first argument or the incremented minor version

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

TGT=${1-none}
if [[ $TGT == "none" ]]; then
  TGT=$(git describe --tags | sed -e 's/v//g')
  TGT=$(echo ${TGT} | awk -F. -v OFS=. '{$NF++;print}')
fi
if [[ ${TGT:0:1} == "v" ]]; then
  TGT = "${TGT:1}"
fi

echo $TGT

find ./app/env -type f -name "env.go" -print0 | xargs -0 sed -i '' -e "s/Version = \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/Version = \"${TGT}\"/g"

make build

git add .
git commit -m "v${TGT}" || true

git tag "v${TGT}"

git push
git push --tags
