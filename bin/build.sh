#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

go build \
	-ldflags "-X main.COMMIT=$(git rev-parse --short HEAD) -X main.LASTMOD=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
	-o linux/r2m-server \
	r2server.go status.go 
