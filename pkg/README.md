# Packaging and Deploying

These are scripts to release a new version of the app

* [`build.sh`](build.sh) - builds the binary and packages it as a `.deb` file
* [`install.sh`](install.sh) - copies the `.deb` on a server and installs it

## Notes

* Config is loaded from `/etc/redirect2me/config.yaml`.
* Logs are written to `/var/log/redirect2me/server.log`.
* Certificates are stored in `/var/lib/redirect2me/certs/` if you use the sample config file.

## Troubleshooting

Common commands to run for troubleshooting:

```
# view the server logs (if the service has started)
cat /var/log/redirect2me/server.log

# view the current status of the service
systemctl status redirect2me.service

# if the service has not started
journalctl -u redirect2me.service -ex
```
