package config

import (
	"os"
)

var (
	//BindAddress where to listen
	BindAddress string

	//TLSEnabled is TLS enabled
	TLSEnabled bool

	//TLSCert TLS cert file to use
	TLSCert string

	//TLSKey TLS key file to use
	TLSKey string
)

func init() {
	TLSEnabled = os.Getenv("TLS_ENABLED") == "true"
	TLSCert = os.Getenv("TLS_CRT")
	TLSKey = os.Getenv("TLS_KEY")
	Rules = []Rule{}

	if TLSEnabled {
		BindAddress = ":8443"
	} else {
		BindAddress = ":8080"
	}
}
