package descend

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func NewClient(caCert, cert, key string) (*http.Client, error) {
	caPool := x509.NewCertPool()
	x509CaCrt, err := ioutil.ReadFile(caCert)
	if ok := caPool.AppendCertsFromPEM(x509CaCrt); !ok {
		return nil, fmt.Errorf("Error appending CA cert from PEM!")
	}
	xcert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{xcert},
			ClientAuth:   tls.RequireAnyClientCert,
			RootCAs:      caPool,
		},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	return client, nil
}
