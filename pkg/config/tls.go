package config

import (
	"crypto/tls"
	"fmt"
)

//TLSConfig returns HTTP TLS config
func TLSConfig() (*tls.Config, error) {
	if !TLSEnabled {
		return nil, nil
	}

	pair, err := tls.LoadX509KeyPair(TLSCert, TLSKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %v", err)
	}

	return &tls.Config{Certificates: []tls.Certificate{pair}}, nil
}
