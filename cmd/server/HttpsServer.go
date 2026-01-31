package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/caddyserver/certmagic"
	"go.uber.org/zap"
)

func HttpsServer(httpsAddress string, mux http.Handler) *http.Server {

	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = Config.AcmeEmail

	if Config.AcmeStaging {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
	} else {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
	}

	zlogger, _ := zap.NewDevelopment()
	certmagic.Default.Logger = zlogger

	certmagic.Default.OnDemand = &certmagic.OnDemandConfig{
		DecisionFunc: func(ctx context.Context, name string) error {
			//LATER: DNS check
			//LATER: algorithm check
			return nil
		},
	}

	magic := certmagic.NewDefault()
	httpsServer := &http.Server{
		Addr:      httpsAddress,
		Handler:   mux,
		TLSConfig: magic.TLSConfig(),
	}

	go func() {
		httpsListenErr := httpsServer.ListenAndServeTLS("", "")
		if httpsListenErr != nil && !errors.Is(httpsListenErr, http.ErrServerClosed) {
			Logger.Error("unable to listen", "address", httpsAddress, "error", httpsListenErr)
		}
	}()

	return httpsServer
}
