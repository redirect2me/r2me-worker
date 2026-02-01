#!/bin/bash
#
# preremove script for redirect2me package
#

# Stop the service if it is running
systemctl stop redirect2me.service || true

# Disable the service to remove startup symlinks
systemctl disable redirect2me.service || true

# Optional: Reload daemon to recognize the unit is gone
systemctl daemon-reload