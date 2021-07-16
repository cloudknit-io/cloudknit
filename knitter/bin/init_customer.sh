#!/bin/bash

BASE_ORG='CompuZest'

while [ $# -gt 0 ]; do
  case "$1" in
    -t|--github-token)
      GITHUB_TOKEN="$2"
      ;;
    -c|--customer)
      CUSTOMER="$2"
      ;;
    -z|--zlifecycle-username)
      ZLIFECYCLE_USER="$2"
      ;;
    *)
      printf "***************************\n"
      printf "* Error: Invalid argument.*\n"
      printf "***************************\n"
      exit 1
  esac
  shift
  shift
done

if [ -z "$GITHUB_TOKEN" ];
then
  echo "Missing -t|--github-token parameter!"
  exit 2
fi

if [ -z "$CUSTOMER" ];
then
  echo "Missing -c|--customer parameter!"
  exit 2
fi

if [ -n "$ZLIFECYCLE_USER" ];
then
  ZLIFECYCLE_GIT_USER="$ZLIFECTCLE_USER"
else
  ZLIFECYCLE_GIT_USER="${CUSTOMER}-zlifecycle"
fi
CUSTOMER_IL_REPO="${CUSTOMER}-il"

curl -H "Authorization: token $GITHUB_TOKEN" "https://api.github.com/orgs/${BASE_ORG}/repos" \
  -d "{ \"name\": \"${CUSTOMER_IL_REPO}\", \"private\": true }"

curl -X PUT -H "Authorization: token $GITHUB_TOKEN" -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/${BASE_ORG}/${CUSTOMER_IL_REPO}/collaborators/${ZLIFECYCLE_GIT_USER}" \
  -d "{ \"permissions\": \"push\" }"

curl -X PUT -H "Authorization: token $GITHUB_TOKEN" -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/${BASE_ORG}/helm-charts/collaborators/${ZLIFECYCLE_GIT_USER}" \
  -d "{ \"permissions\": \"pull\" }"
