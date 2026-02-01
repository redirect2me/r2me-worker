#!/bin/bash
#
# postremove script for redirect2me package
#

# Only run on 'remove' (not upgrade) to avoid deleting user during updates
if [ "$1" = "remove" ] || [ "$1" = "purge" ]; then
    userdel redirect2me || true
fi