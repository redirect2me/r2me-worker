#!/usr/bin/env bash
#
# run locally
#

set -o errexit
set -o pipefail
set -o nounset

GO_PATH=$(go env GOPATH)
export PATH="$GO_PATH/bin:$PATH"

if ! [ -x "$(command -v sqlc)" ]; then
  echo 'ERROR: sqlc is not installed.' >&2
  exit 1
fi

if ! [ -x "$(command -v air)" ]; then
  echo 'ERROR: air is not installed.' >&2
  exit 2
fi

if [ ! -f ".env" ]; then
	echo "INFO: .env file not found"
	exit 3
fi

echo "INFO: loading .env file"
export $(cat .env)

echo "INFO: running with air"
~/go/bin/air
