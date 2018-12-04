#!/bin/bash

echo "DEBUG: root"
curl --header "Host: redirect2me.com" localhost:4000
echo "DEBUG: www"
curl --header "Host: www.redirect2me.com" localhost:4000
echo "DEBUG: noscheme"
curl --header "Host: noscheme.redirect2me.com" localhost:4000
echo "DEBUG: http"
curl --header "Host: http.redirect2me.com" localhost:4000
echo "DEBUG: https"
curl --header "Host: https.redirect2me.com" localhost:4000
echo "DEBUG: path"
curl --header "Host: path.redirect2me.com" localhost:4000
echo "DEBUG: schemepath"
curl --header "Host: schemepath.redirect2me.com" localhost:4000
echo "DEBUG: test1"
curl --header "Host: test1.redirect2me.com" localhost:4000

curl --header "Host: redirect2me.com" localhost:4000/original?query=true
echo "DEBUG: www"
curl --header "Host: www.redirect2me.com" localhost:4000/original?query=true
echo "DEBUG: noscheme"
curl --header "Host: noscheme.redirect2me.com" localhost:4000/original?query=true
echo "DEBUG: http"
curl --header "Host: http.redirect2me.com" localhost:4000/original?query=true
echo "DEBUG: https"
curl --header "Host: https.redirect2me.com" localhost:4000/original?query=true
echo "DEBUG: path"
curl --header "Host: path.redirect2me.com" localhost:4000/original?query=true
echo "DEBUG: schemepath"
curl --header "Host: schemepath.redirect2me.com" localhost:4000/original?query=true
echo "DEBUG: test1"
curl --header "Host: test1.redirect2me.com" localhost:4000/original?query=true

