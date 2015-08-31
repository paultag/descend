package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"pault.ag/go/debian/control"
	"pault.ag/go/descend/descend"
)

func Missing(values ...*string) {
	for _, value := range values {
		if *value != "" {
			continue
		}
		flag.Usage()
		os.Exit(0)
	}
}

func main() {

	caCert := flag.String("ca", "", "CA Cert")
	clientCrt := flag.String("cert", "", "Client Cert")
	clientKey := flag.String("key", "", "Client Key")

	host := flag.String("host", "localhost", "Host to PUT to")
	port := flag.Int("port", 80, "Port to PUT on")
	archive := flag.String("archive", "/", "Archive to PUT to")

	flag.Parse()

	Missing(caCert, clientCrt, clientKey)

	caPool := x509.NewCertPool()
	x509CaCrt, err := ioutil.ReadFile(*caCert)
	if err != nil {
		panic(err)
	}
	if ok := caPool.AppendCertsFromPEM(x509CaCrt); !ok {
		panic(fmt.Errorf("Error appending CA cert from PEM!"))
	}
	cert, err := tls.LoadX509KeyPair(*clientCrt, *clientKey)
	if err != nil {
		panic(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAnyClientCert,
			RootCAs:      caPool,
		},
		DisableCompression: true,
	}

	for _, changesPath := range flag.Args() {
		changes, err := control.ParseChangesFile(changesPath)
		if err != nil {
			panic(err)
		}

		client := &http.Client{Transport: tr}
		err = descend.DoPutChanges(
			client, changes,
			fmt.Sprintf("%s:%d", *host, *port),
			*archive,
		)
		if err != nil {
			panic(err)
		}
	}
}
