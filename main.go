package main

import (
	"net/http"
	"time"

	"github.com/chickenzord/kube-annotate/annotator"
	"github.com/chickenzord/kube-annotate/config"
	"github.com/chickenzord/kube-annotate/web"
	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"
)

var log = config.AppLogger

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/health", web.HealthHandler)
	router.HandleFunc("/metrics", web.MetricsHandler)
	router.HandleFunc("/rules", web.RulesHandler)
	router.Handle("/mutate", &annotator.Annotator{})

	tlsConfig, err := config.TLSConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Starting server at %s", config.BindAddress)

	n := negroni.New(negronilogrus.NewMiddlewareFromLogger(log, "annotate"))
	n.UseHandler(router)
	srv := &http.Server{
		Handler:      n,
		Addr:         config.BindAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if tlsConfig == nil {
		log.Debug("TLS is disabled")
		log.Fatal(srv.ListenAndServe())
	} else {
		log.Debug("TLS is enabled")
		log.Fatal(srv.ListenAndServeTLS(config.TLSCert, config.TLSKey))
	}

}
