

 - [ ] install.sh: adduser is prompting
 - [ ] actually use the r2me user
 - [ ] init script
 - [ ] lookup: parse TXT value and use http/s and/or path/querystring
 - [ ] papertrail
 - [ ] check for short hostnames
 - [ ] log complete source url
 - [ ] recursive lookup
 - [ ] auto action that does lookup, then add/remove

```bash
./r2server --hostname=addwww.redirect2.me --port=80 --debug --verbose --action=addwww
./r2server --hostname=addwww-ssl.redirect2.me --port=80 --verbose --action=addwww
./r2server --hostname=removewww-ssl.redirect2.me --port=80 --verbose --action=removewww
./r2server --hostname=lookup-ssl.redirect2.me --port=80 --verbose --action=lookup
```