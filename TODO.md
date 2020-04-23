# To do

- [ ] cli option to take full address instead of just port
- [ ] support "data:" urls
- [ ] run.sh to reload on changes
- [ ] merge status.go changes
- [ ] check to name is IP address (net.ParseIP)
- [ ] getUrl(r)
- [ ] getHost(r)
- [ ] getClientIP (CF-Connecting-IP, X-Forwarded-For)
- [ ] json logging
- [ ] install.sh: adduser is prompting
- [ ] actually use the r2me user
- [ ] lookup caching
- [ ] init script
- [ ] lookup: parse TXT value and use http/s and/or path/querystring
- [ ] different http status codes: 301, 302, 307...
- [ ] papertrail
- [ ] check for short hostnames (less than len(www))
- [ ] log complete source url
- [ ] recursive lookup
- [ ] interstitial redirect pages
 
## https support

- [ ] keypair table
- [ ] generate csr & store
- [ ] generate private key & store
- [ ] save nonce from LetsEncrypt
- [ ] handler for .well-known
- [ ] shoelace.css for error page
- [ ] handle https

## Maybe, maybe not

- [ ] auto action that does lookup, then add/remove


```bash
./r2server --hostname=addwww.redirect2.me --port=80 --debug --verbose --action=addwww
./r2server --hostname=addwww-ssl.redirect2.me --port=80 --verbose --action=addwww
./r2server --hostname=removewww-ssl.redirect2.me --port=80 --verbose --action=removewww
./r2server --hostname=lookup-ssl.redirect2.me --port=80 --verbose --action=lookup
```
