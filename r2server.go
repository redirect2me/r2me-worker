package main

import (
    "fmt"
    "encoding/json"
    "log"
    "flag"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
)

var (
    verbose = flag.Bool("verbose", true, "verbose logging");
    debug = flag.Bool("debug", false, "print instead of redirect");
    http_port = flag.Int("port", 80, "port to listen on");

    supportUrl = "https://www.redirect2.me/support/faq.html?page=";

    logger = log.New(os.Stdout, "R2W: ", log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC);
)


func getScheme(r *http.Request) (string) {

    proxyValue := r.Header.Get("X-Forwarded-Proto")
    if (proxyValue != "") {
        return proxyValue;
    }

    urlValue := r.URL.Scheme;
    if (urlValue != "") {
        return urlValue;
    }

    return "http";
}

func mapAddWww(r *http.Request) (string) {

    host := r.Host;

    if (strings.HasPrefix(host, "www.")) {
        return supportUrl + "already-has-www.html";
    }

    loc := url.URL(*r.URL);
    loc.Scheme = getScheme(r);
    loc.Host = "www." + host;
    return loc.String();
}

func handler(w http.ResponseWriter, r *http.Request) {

    destination := mapAddWww(r)

    if (*verbose == true) {
        logger.Printf("INFO: redirecting from '%s' to '%s'\n", r.URL.String(), destination);
    }
    if (*debug == true) {
        fmt.Fprintf(w, "DEBUG: redirecting to '%s'", destination);
    } else {
        http.Redirect(w, r, destination, http.StatusTemporaryRedirect);
    }

    data := make(map[string]string)
    data["source"] = r.URL.String();
    data["uri"] = r.RequestURI;
    data["destination"] = destination
    data["user-agent"] = r.Header.Get("user-agent");
    data["referrer"] = r.Header.Get("referer");
    data["path"] = r.URL.Path;
    data["remote-addr"] = r.RemoteAddr;

    extraJson, err := json.Marshal(data)
    if (err != nil) {
        logger.Panicf("ERROR: unable convert extra data to JSON %s\n", err);
        return;
    }

    logger.Printf("INFO: %s\n", extraJson);
}

func main() {

    flag.Parse();

    if (*debug) {
        logger.Printf("DEBUG: running in debug mode\n");
    }

    http.HandleFunc("/", handler)

    if (*verbose) {
        logger.Printf( "INFO: running on port %d\n", *http_port);
    }
    err := http.ListenAndServe(":" + strconv.Itoa(*http_port), nil);
    if (err != nil) {
        logger.Panicf("ERROR: unable to listen on port %d: %s\n", *http_port, err);
    }
}
