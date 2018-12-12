package config

import (
	"os"
)

var (
	//AppName name of the app
	AppName string

	//BindAddress where app server should listen
	BindAddress string

	//BindAddressInternal where internal server should listen
	BindAddressInternal string

	//TLSEnabled is TLS enabled
	TLSEnabled bool

	//TLSCert TLS cert file to use
	TLSCert string

	//TLSKey TLS key file to use
	TLSKey string
)

func init() {
	if val, ok := os.LookupEnv("APP_NAME"); ok {
		AppName = val
	} else {
		AppName = "kube-annotate"
	}

	TLSEnabled = os.Getenv("TLS_ENABLED") == "true"
	TLSCert = os.Getenv("TLS_CRT")
	TLSKey = os.Getenv("TLS_KEY")
	Rules = []Rule{}

	if TLSEnabled {
		BindAddress = ":8443"
	} else {
		BindAddress = ":8080"
	}
	BindAddressInternal = ":8081"
}
