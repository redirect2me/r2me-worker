#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

echo "INFO: creating group"
groupadd --system redirect2me

echo "INFO: creating user"
useradd --system \
			--gid redirect2me \
			--create-home \
			--home-dir /var/lib/redirect2me \
			--shell /usr/sbin/nologin \
			--comment "redirect2me server" \
			redirect2me

echo "INFO: linking service"
ln -s /opt/redirect2me/r2m-server.service /etc/systemd/system/r2m-server.service

echo "INFO: linking socket"
ln -s /opt/redirect2me/r2m-server.socket /etc/systemd/system/r2m-server.socket

echo "INFO: creating log directory"
mkdir /var/log/r2m-server
chown syslog:syslog /var/log/r2m-server

echo "INFO: linking rsyslog conf"
ln -s /opt/redirect2me/rsyslog.conf /etc/rsyslog.d/r2m-server.conf

echo "INFO: restarted rsyslog"
systemctl restart rsyslog.service

echo "INFO: reloading systemd"
systemctl daemon-reload

echo "INFO: enabling service"
systemctl enable r2m-server
systemctl start r2m-server
