#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

go run r2server.go status.go \
	--action=api \
	--debug \
	--endpoint=https://admin.redirect2.me/api/lookup.json \
	--port=4000 \
	--verbose

# --endpoint=http://localhost:4001/api/lookup.json \
