package main

import (
	"context"
	"log"
	"net/http"

	"github.com/caddyserver/certmagic"
)

func HttpsServer(bind string, mux http.Handler) {

	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = Config.AcmeEmail

	if Config.AcmeStaging {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
	} else {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
	}
	certmagic.Default.OnDemand = &certmagic.OnDemandConfig{
		DecisionFunc: func(ctx context.Context, name string) error {
			//LATER: DNS check
			//LATER: algorithm check
			return nil
		},
	}

	magic := certmagic.NewDefault()
	httpsServer := &http.Server{
		Addr:      bind,
		Handler:   mux,
		TLSConfig: magic.TLSConfig(),
	}

	log.Fatal(httpsServer.ListenAndServeTLS("", ""))
}
