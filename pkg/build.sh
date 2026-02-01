#!/usr/bin/env bash
#
# build a deb package for redirect2me
#

set -o errexit
set -o pipefail
set -o nounset

echo "INFO: build starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

SCRIPT_DIR=$(dirname "$0")
REPO_DIR=$(realpath "${SCRIPT_DIR}/..")
BUILD_DIR="${REPO_DIR}/tmp/build"

if [[ ! -d "${BUILD_DIR}" ]]; then
	echo "INFO: creating build directory at ${BUILD_DIR}"
	mkdir -p "${BUILD_DIR}"
else
	echo "INFO: using existing build directory at ${BUILD_DIR}"
fi

export COMMIT=$(git rev-parse --short HEAD)
export LASTMOD=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# if git is dirty, append +dirty to COMMIT
if [[ -n $(git status --porcelain) ]]; then
	export COMMIT="${COMMIT}+dirty"
fi

# Build the Go binary for Linux AMD64
echo "INFO: building Go binary"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a \
    -installsuffix cgo \
    -ldflags "-X main.COMMIT=$COMMIT -X main.LASTMOD=$LASTMOD -extldflags '-static'" \
	-o "${BUILD_DIR}/redirect2me" \
	"${REPO_DIR}/cmd/server/"*.go

BINARY_SIZE=$(wc -c <"${BUILD_DIR}/redirect2me")
echo "INFO: Go binary build complete size=${BINARY_SIZE} bytes"


echo "INFO: building deb package with nfpm"

cd "${REPO_DIR}/pkg/nfpm"
nfpm pkg --packager deb --target "${BUILD_DIR}/redirect2me_1.0.0.deb"

DEB_SIZE=$(wc -c <"${BUILD_DIR}/redirect2me_1.0.0.deb")
echo "INFO: deb package build complete size=${DEB_SIZE} bytes"

echo "INFO: build complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
