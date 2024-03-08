package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	var (
		err              error
		cert             tls.Certificate
		serverCert, body []byte
		pool             *x509.CertPool
		tlsConf          *tls.Config
		transport        *http.Transport
		client           *http.Client
		resp             *http.Response
	)
	if cert, err = tls.LoadX509KeyPair("clientCrt.pem", "clientKey.pem"); err != nil {
		log.Fatalln(err)
	}

	if serverCert, err = os.ReadFile("../server/serverCrt.pem"); err != nil {
		log.Fatalln(err)
	}

	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(serverCert)

	tlsConf = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	tlsConf.BuildNameToCertificate()

	transport = &http.Transport{
		TLSClientConfig: tlsConf,
	}
	client = &http.Client{
		Transport: transport,
	}

	if resp, err = client.Get("https://localhost:9443/hello"); err != nil {
		log.Fatalln(err)
	}

	if body, err = io.ReadAll(resp.Body); err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	fmt.Printf("Success: %s\n", body)
}
