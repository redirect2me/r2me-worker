package main

import (
    "context"
    "fmt"
    "encoding/json"
    "log"
    "flag"
    "net"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
)

var (
    verbose = flag.Bool("verbose", true, "verbose logging");
    debug = flag.Bool("debug", false, "print instead of redirect");
    port = flag.Int("port", 80, "port to listen on");
    hostname = flag.String("hostname", "localhost", "hostname of this server");
    action = flag.String("action", "lookup", "action [lookup|addwww|removewww]");

    supportUrl = "https://www.redirect2.me/support/faq.html?page=";

    logger = log.New(os.Stdout, "R2W: ", log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC);
    mapFunc func(*http.Request) string;
)

func customDNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
    d := net.Dialer{}
    return d.DialContext(ctx, "udp", "1.1.1.1:53")
}

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

func lookupLow(domain string) (string) {
    prefix := "redirect2me=";

    r := net.Resolver{
        PreferGo: true,
        Dial: customDNSDialer,
    }

    results, err := r.LookupTXT(context.Background(), domain);

    if (err != nil || len(results) == 0) {
        return "";
    }

    for rIndex := 0; rIndex < len(results); rIndex++ {
        result := results[rIndex];
        logger.Printf("INFO: %s(%d): %s\n", domain, rIndex, result);
        if (len(result) > len(prefix) + 1 && result[0:len(prefix)] == prefix) {
            return results[rIndex][len(prefix):]
        }
    }

    return ""
}

func mapAddWww(r *http.Request) (string) {

    host := r.Host;

    if (strings.HasPrefix(host, "www.")) {
        return supportUrl + "error-already-has-www.html";
    }

    loc := url.URL(*r.URL);
    loc.Scheme = getScheme(r);
    loc.Host = "www." + host;
    return loc.String();
}

func mapRemoveWww(r *http.Request) (string) {

    host := r.Host;

    if (!strings.HasPrefix(host, "www.")) {
        return supportUrl + "error-does-not-have-www.html";
    }

    loc := url.URL(*r.URL);
    loc.Scheme = getScheme(r);
    loc.Host = host[4:];
    return loc.String();
}

func mapLookup(r *http.Request) (string) {

    host := r.Host;

    newHost := lookupLow(host);

    if (newHost == "") {
        return supportUrl + "error-lookup-not-found.html";
    }

    loc := url.URL(*r.URL);
    loc.Scheme = getScheme(r);
    loc.Host = newHost;
    return loc.String();
}

func redirect_handler(w http.ResponseWriter, r *http.Request) {

    /* LATER:
    if (isAddress(r)) {
        return supportUrl + "error-host-is-ip-address.html";
    }
    */

    destination := mapFunc(r);

    if (*verbose == true) {
        logger.Printf("INFO: redirecting from '%s' to '%s'\n", r.URL.String(), destination);
    }
    if (*debug == true) {
        fmt.Fprintf(w, "DEBUG: redirecting to '%s'\n", destination);
    } else {
        http.Redirect(w, r, destination, http.StatusTemporaryRedirect);
    }

    data := make(map[string]string)
    data["source_host"] = r.Host;
    data["source_url"] = r.URL.String();
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

func www_handler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://www.redirect2.me/", http.StatusTemporaryRedirect);
}

func main() {

    flag.Parse();

    switch *action {
        case "lookup":
            mapFunc = mapLookup;
        case "addwww":
            mapFunc = mapAddWww;
        case "removewww":
            mapFunc = mapRemoveWww;
        default:
            logger.Panicf("ERROR: invalid action '%s'\n", action);
    }

    if (*debug) {
        logger.Printf("DEBUG: running in debug mode\n");
    }

    http.HandleFunc("/", redirect_handler);
	http.HandleFunc(*hostname + "/status.json", Status_handler);
	http.HandleFunc(*hostname + "/", www_handler);

    if (*verbose) {
        logger.Printf("INFO: running on port %d\n", *port);
        logger.Printf("INFO: hostname is %s\n", *hostname);
        logger.Printf("INFO: action is %s\n", *action);
    }
    err := http.ListenAndServe(":" + strconv.Itoa(*port), nil);
    if (err != nil) {
        logger.Panicf("ERROR: unable to listen on port %d: %s\n", *port, err);
    }
}
