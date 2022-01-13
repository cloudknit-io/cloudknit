#!/bin/sh
#
# This script downloads an asset from latest or specific Github release of a
# private repo. Feel free to extract more of the variables into command line
# parameters.
#
# PREREQUISITES
#
# curl, wget, jq
#
# USAGE
#
# Set all the variables inside the script, make sure you chmod +x it, then
# to download specific version to my_app.tar.gz:
#
#     gh_download.sh -v v0.0.9 -f zlifecycle-internal-cli_0.0.9_linux_arm64.zip -t <GITHUB_TOKEN> -r CompuZest/zlifecycle-internal-cli -o cli.zip
#
# to download latest version:
#
#     gh_download.sh -v latest -f zlifecycle-internal-cli_<latest semver>_linux_arm64.zip -t <GITHUB_TOKEN> -r CompuZest/zlifecycle-internal-cli -o cli.zip
#
# If your version/tag doesn't match, the script will exit with error.

GITHUB="https://api.github.com"

if ! command -v jq &> /dev/null
then
    echo "jq could not be found"
    exit
fi

if ! command -v wget &> /dev/null
then
    echo "wget could not be found"
    exit
fi

if ! command -v curl &> /dev/null
then
    echo "curl could not be found"
    exit
fi

while [ $# -gt 0 ]; do
  case "$1" in
    -v|--version)
      VERSION="$2"
      ;;
    -t|--github-token)
      TOKEN="$2"
      ;;
    -f|--file)
      FILE="$2"
      ;;
    -r|--repository)
      REPO="$2"
      ;;
    -o|--output)
      OUTPUT="$2"
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

if [ -z "$TOKEN" ];
then
  echo "Missing -t|--github-token parameter!"
  exit 2
fi

if [ -z "$VERSION" ];
then
  echo "Missing -v|--version parameter!"
  exit 2
fi

if [ -z "$FILE" ];
then
  echo "Missing -f|--file parameter!"
  exit 2
fi

if [ -z "$REPO" ];
then
  echo "Missing -r|--repository parameter!"
  exit 2
fi

alias errcho='>&2 echo'

function gh_curl() {
  curl -H "Authorization: Bearer $TOKEN" \
       -H "Accept: application/vnd.github.v3.raw" \
       $@
}

parser=".assets | map(select(.name == \"$FILE\"))[0].id"
if [ "$VERSION" = "latest" ]; then
  # Github should return the latest release first.
  asset_id=`gh_curl -s $GITHUB/repos/$REPO/releases/latest | jq "$parser"`
else
  asset_id=`gh_curl -s $GITHUB/repos/$REPO/releases/tags/$VERSION | jq "$parser"`
fi;

if [ -z "$asset_id" ]; then
  errcho "ERROR: version not found $VERSION"
  exit 1
fi;
if [ "$asset_id" = "null" ]; then
  errcho "ERROR: file $FILE not found in version $VERSION"
  exit 2
fi;

wget -q --auth-no-challenge --header='Accept:application/octet-stream' \
  https://$TOKEN:@api.github.com/repos/$REPO/releases/assets/$asset_id \
  -O "$OUTPUT"
