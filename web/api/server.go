package api

import (
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/mysterion/avrp/web"
)

//go:embed certs/*
var Certs embed.FS

func getTlsConfig() (*tls.Config, error) {
	config := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		Certificates:             make([]tls.Certificate, 1),
	}

	certBlock, err := Certs.ReadFile("certs/cert.pem")
	if err != nil {
		return nil, err
	}

	keyBlock, err := Certs.ReadFile("certs/key.pem")
	if err != nil {
		return nil, err
	}

	config.Certificates[0], err = tls.X509KeyPair(certBlock, keyBlock)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func Start(port int) {
	tlsConfig, err := getTlsConfig()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc(listPath, listHandler)

	fileServer := http.FileServer(http.Dir(servDir))
	mux.Handle(filePath, http.StripPrefix(filePath, fileServer))

	server := &http.Server{
		Addr:      fmt.Sprintf(":%v", port),
		TLSConfig: tlsConfig,
		Handler:   mux,
	}

	go func() {
		fmt.Println("Server listening on: ")
		for _, ip := range web.GetIps() {
			fmt.Printf("https://%v:%v\n", ip, port)
		}
		err = server.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatal(err)
		}
	}()

	// todo : graceful shutdown/startup
	select {}
}
