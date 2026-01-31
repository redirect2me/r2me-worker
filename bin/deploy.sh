#!/bin/bash
#
# install the r2me server on a vanilla Ubuntu box
#

set -o errexit
set -o pipefail
set -o nounset

SERVER=premium.redirect2.me

echo "INFO: starting deploy on ${SERVER} at $(date)"

echo "INFO: copying files..."
scp r2server root@${SERVER}:/app

echo "INFO: changing ownership..."
ssh root@${SERVER} "chown --recursive r2me:r2me /app"


echo "INFO: deploy complete on ${SERVER} at $(date)"
