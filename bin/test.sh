#!/bin/bash
#
# run tests locally
#

set -o errexit
set -o pipefail
set -o nounset

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPO_DIR="$(dirname "$SCRIPT_DIR")"
TMP_DIR="${REPO_DIR}/tmp"

echo "INFO: testing starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

if [ ! -d "${TMP_DIR}" ]; then
	echo "INFO: creating tmp directory at ${TMP_DIR}"
	mkdir "${TMP_DIR}"
else
	echo "INFO: using existing tmp directory '${TMP_DIR}'"
fi

#
# build the binary
#
#echo "INFO: building binary"
go build \
	-ldflags "-X main.COMMIT=$(git rev-parse --short HEAD) -X main.LASTMOD=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.BUILTBY=test.sh " \
	-o ${TMP_DIR}/redirect2me \
	${REPO_DIR}/cmd/server

ACTIONS=(addwww removewww lookup api auto)

for ACTION in "${ACTIONS[@]}"; do
	TEST_CSV="${REPO_DIR}/test/${ACTION}.csv"
	if [ ! -f "${TEST_CSV}" ]; then
		echo "WARNING: no tests for ${ACTION} (file ${TEST_CSV} not found), skipping"
		continue
	fi
	echo "INFO: $ACTION tests"

	echo "INFO: starting server in background"

	LOG_FORMAT=text \
	LOG_LEVEL=info \
	${TMP_DIR}/redirect2me \
		--action=${ACTION} \
		--http_port=8080 \
		--log_level=info \
		--log_format=text \
		&
	SERVER_PID=$!
	echo "INFO: server started with PID ${SERVER_PID}"

	echo "INFO: waiting for server to be ready..."
	while ! nc -z localhost 8080; do
		sleep 0.1 
		echo -n "."
	done
	echo ""

	curl \
		--silent \
		--show-error \
		--max-time 15 \
		http://localhost:8080/status.json

	while IFS=, read -r URL EXPECTED; do
		if [[ "${URL}" == '#'* ]]; then
			echo "DEBUG: skipping comment line '${URL}'"
			continue
		fi
		echo "DEBUG: testing '${URL}' expecting '${EXPECTED}'"
		#curl --header "Host: ${ACTION}.redirect2me.com" localhost:4000/original?query=true
	done < "${TEST_CSV}"

	echo "INFO: killing server with PID ${SERVER_PID}"
	kill -SIGTERM ${SERVER_PID}

	echo "INFO: waiting for server to exit"
	wait ${SERVER_PID} 2>/dev/null || true

	echo "INFO: ${ACTION} tests complete"

done

echo "INFO: testing complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
