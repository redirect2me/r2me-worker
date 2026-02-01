#!/bin/bash
#
# postinstall script for redirect2me package
#

# Create the user
if ! getent user redirect2me > /dev/null 2>&1; then
    useradd --system --no-create-home --shell /usr/sbin/nologin redirect2me
fi

# Fix log directory ownership if it wasn't set by the package manager
mkdir -p /var/log/redirect2me
chown redirect2me:redirect2me /var/log/redirect2me
chmod 750 /var/log/redirect2me

# Fix the config directory and file ownership if it wasn't set by the package manager
mkdir -p /etc/redirect2me
chown -R redirect2me:redirect2me /etc/redirect2me
chmod 640 /etc/redirect2me/config.yaml

# Enable and start the server
systemctl daemon-reload
systemctl enable redirect2me
systemctl start redirect2me