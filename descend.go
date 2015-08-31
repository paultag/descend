package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"pault.ag/go/debian/control"
	"pault.ag/go/fancytext"
)

func DputFile(client *http.Client, host, archive, fpath string) error {
	filename := path.Base(fpath)
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("https://%s/%s/%s", host, archive, filename),
		nil,
	)

	if err != nil {
		return err
	}
	fd, err := os.Open(fpath)
	if err != nil {
		return err
	}
	req.Body = fd
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func DoPutChanges(client *http.Client, changes *control.Changes, host, archive string) error {
	root := path.Dir(changes.Filename)
	for _, file := range changes.Files {
		done := fancytext.BooleanFormatSpinner(fmt.Sprintf("%%s   %s", file.Filename))
		err := DputFile(client, host, archive, path.Join(root, file.Filename))
		done(err == nil)
		if err != nil {
			return err
		}
	}
	return nil
}

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
		err = DoPutChanges(client, changes, "localhost:1984", "foo")
		if err != nil {
			panic(err)
		}
	}
}
