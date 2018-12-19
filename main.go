package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chickenzord/kube-annotate/annotator"
	"github.com/chickenzord/kube-annotate/config"
	"github.com/chickenzord/kube-annotate/web"
	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	prommiddleware "github.com/slok/go-prometheus-middleware"
	promnegroni "github.com/slok/go-prometheus-middleware/negroni"
	"github.com/urfave/negroni"
)

var log = config.AppLogger

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	if rulesFile, ok := config.LoadRules(); ok {
		log.Infof("loaded rules from %s", rulesFile)
	} else {
		log.Warn("no rules file defined")
	}

	tlsConfig, err := config.TLSConfig()
	if err != nil {
		log.WithError(err).Fatal("invalid TLS config")
	}

	rInternal := mux.NewRouter()
	rInternal.HandleFunc("/health", web.HealthHandler)
	rInternal.Handle("/metrics", promhttp.Handler())
	nInternal := negroni.New()
	nInternal.UseHandler(rInternal)
	internal := &http.Server{
		Handler:      nInternal,
		Addr:         config.BindAddressInternal,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	mLogger := negronilogrus.NewMiddlewareFromLogger(log.Logger, config.AppName)
	mProm := promnegroni.Handler("", prommiddleware.NewDefault())

	rServer := mux.NewRouter()
	rServer.HandleFunc("/mutate", annotator.MutateHandler)
	rServer.HandleFunc("/rules", annotator.RulesHandler)
	nServer := negroni.New(mLogger, mProm)
	nServer.UseHandler(rServer)
	server := &http.Server{
		Handler:      nServer,
		Addr:         config.BindAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Infof("API server is listening on %s", server.Addr)
		var err error
		if tlsConfig == nil {
			log.Debug("API server TLS is disabled")
			err = server.ListenAndServe()
		} else {
			log.Debug("API server TLS is enabled")
			err = server.ListenAndServeTLS(config.TLSCert, config.TLSKey)
		}
		if err != http.ErrServerClosed {
			log.WithError(err).Fatalf("API server failed to listen on %s", server.Addr)
		}
	}()

	go func() {
		log.Infof("internal server is listening on %s", internal.Addr)
		if err := internal.ListenAndServe(); err != http.ErrServerClosed {
			log.WithError(err).Fatalf("internal server failed to listen on %s", internal.Addr)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Infof("stopping gracefully")
	if err := internal.Shutdown(ctx); err != nil {
		log.WithError(err).
			Infof("failed to stop internal server gracefully")
	} else {
		log.Infof("internal server gracefully stopped")
	}
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).
			Infof("failed to stop API server gracefully")
	} else {
		log.Infof("API server gracefully stopped")
	}

}
