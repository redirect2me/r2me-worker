Redirect2www
============

Server to redirect all requests from a naked domain to the `www.` version.

To Do
-----
 - [x] go web server
 - [x] test redirects
 - [x] command line flags
 - [ ] handle IP address hostnames
 - [x] stdout logging
 - [ ] error page (on www)
 - [x] use port from args/env
 - [x] log referrer
 - [ ] getHost
 - [ ] getClientIP (CF-Connecting-IP, X-Forwarded-For)
 - [ ] keypair table
 - [ ] generate csr & store
 - [ ] generate private key & store
 - [ ] save nonce from LetsEncrypt
 - [ ] handler for .well-known
 - [ ] shoelace.css for error page
 - [ ] handle https
 - [ ] errors should have link to FAQ page on main website
 - [ ] flag for testing db on startup
 - [ ] arg for db connection string
 - [ ] server-side google analytics (https://developers.google.com/analytics/devguides/collection/protocol/v1/)

Later
-----
 - [ ] more data in dblog detail: schema, querystring
 - [ ] check for oversize field lengths in dblog
 - [ ] update hostlog (in trigger, maybe after PGSQL 9.6 upgrade)
 - [ ] CORS headers (why is this any use?)

DNS library: https://github.com/miekg/dns https://miek.nl/2014/august/16/go-dns-package/

DNS cache: https://github.com/rs/dnscache

Logging with timer: https://github.com/hashicorp/http-echo/blob/master/handlers.go

JSON logging: https://github.com/rs/zerolog

DNS debugging: https://github.com/rs/dnstrace

Public Suffix List: https://godoc.org/golang.org/x/net/publicsuffix

https://github.com/gorilla/reverse/blob/master/matchers.go#L238-L249
```go
func getHost(r *http.Request) string {
    if r.URL.IsAbs() {
        host := r.Host
        // Slice off any port information.
        if i := strings.Index(host, ":"); i != -1 {
            host = host[:i]
        }
        return host
    }
    return r.URL.Host
}
```

Database
--------

RequestLog
 - date
 - time
 - scheme
 - host
 - port
 - path
 - querystring
 - user-agent
 - client ip
 - headers


HostLog
 - host
 - firsthit
 - lasthit

KeyPair
 - host
 - acme_account
 - private_key
 - csr
 - validation_token
 - certificate
 - expires_at
 - environment (dev/prod)

Cert flow
---------
 * gen private key
 * create acme account
 * generate pkcs10 csr
 * submit csr
 * receive challenges
 * validation of challenge
 * get certificate


Build Notes
-----------
export GOPATH=/home/amarcuse/gocode
go get gopkg.in/alecthomas/kingpin.v2

https://godoc.org/golang.org/x/crypto/acme
https://blog.golang.org/context

https://github.com/golang/crypto/blob/master/acme/autocert/autocert.go#L470

key, err := rsa.GenerateKey(rand.Reader, 2048)
if err != nil {
	log.Fatal(err)
}
client := &Client{Key: key}

if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}

https://godoc.org/golang.org/x/crypto/acme#Client.CreateCert

FetchCert

background task to do stuff:
https://github.com/nf/webfront/blob/master/main.go