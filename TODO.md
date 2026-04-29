# To do

## Basics
- [ ] metrics: free disk/memory
- [ ] certmagic logging https://pkg.go.dev/go.uber.org/zap#section-readme
- [ ] source file/line value in log statements needs to be one higher
- [ ] timeouts ([context with timeout](https://blog.golang.org/context))
- [ ] check for short hostnames (less than len(www))
- [ ] Map: `api` action
- [ ] Map: `auto` action that tries lookup, then either add or remove www

## Future 
- [ ] lookup: parse TXT value to get key/value pairs (space separated)
- [ ] interstitial redirect pages
- [ ] make sure client IP is correct (CF-Connecting-IP, X-Forwarded-For, etc)
- [ ] lookup caching
- [ ] support "data:" urls
- [ ] [DNS caching](https://github.com/rs/dnscache) (or just regular result caching?)

## Maybe
- [ ] different http status codes: 301, 302, 307... (or maybe not: too dangerous)
- [ ] check names vs [public suffix list](https://pkg.go.dev/golang.org/x/net/publicsuffix) (non-public could be premium)
- [ ] flag (and alternate destination) for bots
- [ ] coming soon page (=interstital with no link)
- [ ] signup page (but could just be link to Google Forms)
- [ ] whitelist of allowed domains
- [ ] log response size ([example](https://github.com/hashicorp/http-echo/blob/master/handlers.go))