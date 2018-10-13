#!/bin/bash
#
# install the r2me server on a vanilla Ubuntu box
#

SERVER=removewww.redirect2.me

echo "INFO: starting install on ${SERVER} at $(date)"

echo "INFO: adding user"
ssh root@${SERVER} "adduser --home /app --disabled-login r2me"

echo "INFO: copying files..."
scp r2server root@${SERVER}:/app

echo "INFO: changing ownership..."
ssh root@${SERVER} "chown --recursive r2me:r2me /app"


echo "INFO: install complete on ${SERVER} at $(date)"
