#!/usr/bin/env bash
#
# build a deb package for redirect2me
#

set -o errexit
set -o pipefail
set -o nounset

SCRIPT_DIR=$(dirname "$0")
REPO_DIR=$(realpath "${SCRIPT_DIR}/..")
BUILD_DIR="${REPO_DIR}/tmp/build"

echo "INFO: install starting at $(date -u +"%Y-%m-%dT%H:%M:%SZ")"

IDENTITY="~/.ssh/id_ed25519"
SERVER=${1:-127.0.0.1}
echo "INFO: target server is ${SERVER}"

DEB="${BUILD_DIR}/redirect2me_1.0.0.deb"
echo "INFO: deb file is ${DEB}"

echo "INFO: copying deb file"
scp -i "${IDENTITY}" "${DEB}" "root@${SERVER}:"

echo "INFO: installing deb package"
ssh -i "${IDENTITY}" "root@${SERVER}" "dpkg --force-confdef -i $(basename ${DEB})"

echo "INFO: install complete at $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
