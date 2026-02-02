#!/usr/bin/env bash
#
# update an existing redirect2me server
#

set -o errexit
set -o pipefail
set -o nounset

SCRIPT_DIR=$(dirname "$0")
REPO_DIR=$(realpath "${SCRIPT_DIR}/..")
BUILD_DIR="${REPO_DIR}/tmp/build"

echo "INFO: update starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

ENV_FILE="${SCRIPT_DIR}/.env"
if [ -f "${ENV_FILE}" ]; then
	echo "INFO: loading environment from ${ENV_FILE}"
	export $(grep -v '^#' "${ENV_FILE}" | xargs)
else
	echo "WARN: no environment file found at ${ENV_FILE}, proceeding with existing environment"
fi

if [ "${IDENTITY:-BAD}" == "BAD" ]; then
	echo "ERROR: IDENTITY environment variable not set"
	exit 1
fi

SERVER=${1:-BAD}
if [ "${SERVER}" == "BAD" ]; then
	echo "ERROR: server argument not set"
	echo "USAGE: $0 {server_ip_or_hostname}"
	exit 1
fi
echo "INFO: target server is ${SERVER}"

DEB_FILE=$(find "${BUILD_DIR}" -name "*.deb")
if [ -z "${DEB_FILE}" ]; then
	echo "ERROR: no deb file found in ${BUILD_DIR}"
	exit 1
fi
echo "INFO: deb file found at ${DEB_FILE}"

echo "INFO: copying deb file"
scp -i "${IDENTITY}" "${DEB_FILE}" "root@${SERVER}:"

echo "INFO: installing deb package"
ssh -i "${IDENTITY}" "root@${SERVER}" "dpkg --force-confdef -i $(basename ${DEB_FILE})"

echo "INFO: update complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
