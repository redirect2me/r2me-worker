# redirect2me server [<img alt="redirect2me Logo" src="https://www.redirect2.me/favicon.svg" height="90" align="right" />](https://www.redirect2.me/)

Server that actually handle redirects.

## License

Copyright (c) 2018 by Andrew Marcuse.  All Rights Reserved.

To Do
-----

Later
-----

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