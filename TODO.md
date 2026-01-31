# To do

## Basics


- [ ] Configurable directory for certs
- [ ] certmagic logging https://pkg.go.dev/go.uber.org/zap#section-readme
- [ ] source value in log statements needs to be one higher

## Deploy

- [ ] init script
- [ ] logrotate
- [ ] build deb


- [ ] check for short hostnames (less than len(www))


## Future 

- [ ] lookup: parse TXT value and use http/s and/or path/querystring
- [ ] interstitial redirect pages
- [ ] check to name is IP address (net.ParseIP)
- [ ] make sure client IP is correct (CF-Connecting-IP, X-Forwarded-For, etc)
- [ ] lookup caching
- [ ] different http status codes: 301, 302, 307... (maybe not: too dangerous)
- [ ] support "data:" urls

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
