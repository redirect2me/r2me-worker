package main

import (
	"errors"
	"net/http"

	"github.com/caddyserver/certmagic"
	"go.uber.org/zap"
)

func HttpsServer(httpsAddress string, mux http.Handler) (*http.Server, *certmagic.ACMEIssuer) {

	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = Config.AcmeEmail
	certmagic.DefaultACME.Profile = "shortlived"

	if Config.AcmeStaging {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
	} else {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
	}

	zlogger, _ := zap.NewDevelopment()
	certmagic.Default.Logger = zlogger

	certmagic.Default.Storage = &certmagic.FileStorage{Path: Config.CertDir}

	certmagic.Default.OnDemand = &certmagic.OnDemandConfig{
		DecisionFunc: IsCertEligible,
	}

	magic := certmagic.NewDefault()
	httpsServer := &http.Server{
		Addr:      httpsAddress,
		Handler:   mux,
		TLSConfig: magic.TLSConfig(),
	}

	myACME := certmagic.NewACMEIssuer(magic, certmagic.DefaultACME)

	go func() {
		httpsListenErr := httpsServer.ListenAndServeTLS("", "")
		if httpsListenErr != nil && !errors.Is(httpsListenErr, http.ErrServerClosed) {
			Logger.Error("unable to listen", "address", httpsAddress, "error", httpsListenErr)
		}
	}()

	return httpsServer, myACME
}
