#!/usr/bin/env sh

case "$(uname -s)" in
Darwin)
  OS=darwin
  ;;
Linux)
  OS=linux
  ;;
*)
  echo "Unsupported"
  exit 1
  ;;
esac

REPO_NAME="cli-client"
REPO_OWNER="sentinel-official"

VERSION=$(
  curl --silent "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" |
    tr -d "\n" |
    sed 's/.*"tag_name": "v//' |
    sed 's/".*//'
)

FILE_NAME="sentinelcli-${VERSION}-${OS}-amd64"
FILE_PATH="$(mktemp -d)/${FILE_NAME}"
ASSET_URL="https://github.com/sentinel-official/cli-client/releases/download/v${VERSION}/${FILE_NAME}"

curl --location --output "${FILE_PATH}" "${ASSET_URL}"
chmod +x "${FILE_PATH}" &&
  sudo chown root "${FILE_PATH}" &&
  sudo mv "${FILE_PATH}" /usr/local/bin/sentinelcli
