package main

import (
    "context"
    "flag"
    "log"
    "net"
    "os"
)

var (
    verbose = flag.Bool("verbose", true, "verbose logging");
    debug = flag.Bool("debug", false, "print instead of redirect");
    http_port = flag.Int("port", 80, "port to listen on");

    logger = log.New(os.Stdout, "DNSCHECK: ", log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC);
)


func lookupLow(domain string) (string) {
    prefix := "redirect2me=";

        r := net.Resolver{
            PreferGo: true,
            Dial: CustomDNSDialer,
        }

    results, err := r.LookupTXT(context.Background(), domain);

    if (err != nil || len(results) == 0) {
        return "";
    }

    for rIndex := 0; rIndex < len(results); rIndex++ {
        logger.Printf("INFO: %s: %s (%d)\n", domain, results[rIndex], rIndex);
        if (results[rIndex][0:len(prefix)] == prefix) {
            return results[rIndex][len(prefix):]
        }
    }

    return ""
}

func CustomDNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
    d := net.Dialer{}
    return d.DialContext(ctx, "udp", "1.1.1.1:53")
}

func main() {

    flag.Parse();

    if (*debug) {
        logger.Printf("DEBUG: running in debug mode\n");
    }

    domains := flag.Args();

    if (len(domains) == 0) {
        logger.Printf("ERROR: list some domains to check\n");
        os.Exit(1);
    }
    ctx := context.Background()


    r := net.Resolver{
        PreferGo: true,
        Dial: CustomDNSDialer,
    }

    for index := 0; index < len(domains); index++ {
        domain := domains[index];
        logger.Printf("INFO: %s starting\n", domain);

        results, err := r.LookupTXT(ctx, domain);
        if (err != nil) {
            logger.Printf("ERROR: %s %s\n", domain, err.Error());
            continue;
        }

        for rIndex := 0; rIndex < len(results); rIndex++ {
            logger.Printf("INFO: %s: %s (%d)\n", domain, results[rIndex], rIndex);
        }

        logger.Printf("OUTPUT: lookupLow=%s\n", lookupLow(domain));
    }

    logger.Printf("INFO: complete\n");
}
