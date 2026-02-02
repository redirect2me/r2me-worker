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
TMP_DIR="${REPO_DIR}/tmp"

echo "INFO: server starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

ENV_FILE="${SCRIPT_DIR}/.env"
if [ -f "${ENV_FILE}" ]; then
	echo "INFO: loading environment from ${ENV_FILE}"
	export $(grep -v '^#' "${ENV_FILE}" | xargs)
else
	echo "WARN: no environment file found at ${ENV_FILE}, proceeding with existing environment"
fi

if [ "${TOKEN:-BAD}" == "BAD" ]; then
	echo "ERROR: TOKEN environment variable not set"
	exit 1
fi

if [ "${IDENTITY:-BAD}" == "BAD" ]; then
	echo "ERROR: IDENTITY environment variable not set"
	exit 1
fi

ACTION=${1:-BAD}
if [ "${ACTION}" == "BAD" ]; then
	echo "ERROR: action argument not set"
	echo "USAGE: $0 {lookup|addwww|removewww}"
	exit 1
fi
CONFIG_FILE="${REPO_DIR}/pkg/config/${ACTION}.yaml"
if [ ! -f "${CONFIG_FILE}" ]; then
	echo "ERROR: config file for ${ACTION} not found at ${CONFIG_FILE}"
	exit 1
fi

DEB_FILE=$(find "${BUILD_DIR}" -name "*.deb")
if [ -z "${DEB_FILE}" ]; then
	echo "ERROR: no deb file found in ${BUILD_DIR}"
	exit 1
fi
echo "INFO: deb file found at ${DEB_FILE}"

IDENTITY_PATH="${IDENTITY/#\~/$HOME}"
echo "INFO: getting ssh key id from '${IDENTITY_PATH}.pub'"
SSH_KEY_ID=$(ssh-keygen -l -E md5 -f "${IDENTITY_PATH}.pub" | awk -F ' ' '{print $2}' | cut -f 2- -d ':')
echo "INFO: ssh key id is ${SSH_KEY_ID}"

echo "INFO: calculating unique id for droplet"
UNIQUE_ID=$(date  -u +%Y%m%dT%H%M%SZ)
echo "INFO: unique id is ${UNIQUE_ID}"

echo "INFO: creating droplet"
CREATE_JSON=$(curl -X POST \
	--header 'Authorization: Bearer '${TOKEN}'' \
	--header 'Content-Type: application/json' \
	-d '{"name":"r2me-'${ACTION}'-'${UNIQUE_ID}'",
		"size":"s-1vcpu-512mb-10gb",
		"region":"nyc3",
		"image":"ubuntu-24-04-x64",
		"ssh_keys":["'${SSH_KEY_ID}'"],
		"monitoring":true
		}' \
	--silent \
	--show-error \
	--url "https://api.digitalocean.com/v2/droplets"
)
echo "INFO: created droplet"
echo "${CREATE_JSON}" | jq . >"${TMP_DIR}/create.json"
DROPLET_ID=$(echo "${CREATE_JSON}" | jq -r '.droplet.id')
echo "INFO: droplet id is ${DROPLET_ID}"

echo "INFO: checking droplet status"
STATUS="new"
IPADDRESS="not-set"
while [ "${STATUS}" != "active" ]; do
	STATUS_JSON=$(curl \
		--header "Authorization: Bearer ${TOKEN}" \
		--header 'Content-Type: application/json' \
		--show-error \
		--silent \
		--url "https://api.digitalocean.com/v2/droplets/${DROPLET_ID}"
	)
	echo "${STATUS_JSON}" | jq . >"${TMP_DIR}/status.json"
	STATUS=$(echo "${STATUS_JSON}" | jq -r '.droplet.status')
	IPADDRESS=$(echo "${STATUS_JSON}" | jq -r '.droplet.networks.v4[] | select(.type=="public") | .ip_address')
	echo "INFO: droplet status is ${STATUS} at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
	sleep 15
done

echo "INFO: droplet is active with IP address ${IPADDRESS}"

#
# add droplet to known hosts
#
echo "INFO: adding ssh identity for ${IPADDRESS} to known hosts"
ssh-keygen -R "${IPADDRESS}"
# not sure why this doesn't work, but it doesn't...
#ssh-keyscan -H "${IPADDRESS}" >> ~/.ssh/known_hosts
ssh -i "${IDENTITY_PATH}" -o "StrictHostKeyChecking=no" -o "ConnectionAttempts=10" "root@${IPADDRESS}" "whoami"

#
# copy the appropriate config file
#
echo "INFO: copying config file ${CONFIG_FILE}"
ssh -i "${IDENTITY_PATH}" "root@${IPADDRESS}" "mkdir -p /etc/redirect2me"
scp -i "${IDENTITY_PATH}" "${CONFIG_FILE}" "root@${IPADDRESS}:/etc/redirect2me/config.yaml"

#
# copy the deb file
#
echo "INFO: copying deb file ${DEB_FILE}"
scp -i "${IDENTITY}" "${DEB_FILE}" "root@${IPADDRESS}:"

echo "INFO: installing deb package"
ssh -i "${IDENTITY}" "root@${IPADDRESS}" "dpkg --force-confdef -i $(basename ${DEB_FILE})"
echo "INFO: server complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
