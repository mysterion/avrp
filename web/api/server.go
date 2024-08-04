package api

import (
	"context"
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mysterion/avrp/internal/utils"
	"github.com/mysterion/avrp/web"
	"github.com/rs/cors"
)

//go:embed certs/*
var Certs embed.FS

func GetTlsConfig() (*tls.Config, error) {
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

// port = 0, for a random port
//
// mux can be nil
//
// tlsConfig can be nil
//
// returns [<-ready], [<-done], [*http.Server]
func New(port int, mux http.Handler, tlsConfig *tls.Config) (<-chan bool, <-chan bool, *http.Server) {

	ready := make(chan bool, 1)
	done := make(chan bool, 1)

	server := http.Server{
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	go func() {
		defer close(ready)
		defer close(done)

		listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", "0.0.0.0", port))
		utils.Panic(err)

		server.Addr = listener.Addr().String()

		ready <- true

		err = server.ServeTLS(listener, "", "")

		if err != nil && err != http.ErrServerClosed {
			fmt.Println(err.Error())
			panic(err)
		}
	}()

	return ready, done, &server
}

func Start(port int) {
	tlsConfig, err := GetTlsConfig()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	const filePath = "/file/"
	var fileHandler = http.StripPrefix(filePath, http.FileServer(http.Dir(servDir)))

	listH := listHandler
	thumbH := thumbHandler

	distH := distHandler
	fileH := fileHandler

	if utils.DEV {
		log.Println("Enabling CORS")
		listH = wrapCors(listHandler)
		thumbH = wrapCors(thumbHandler)

		distH = cors.AllowAll().Handler(distHandler)
		fileH = cors.AllowAll().Handler(fileHandler)
	}

	mux.Handle(distPath, distH)
	mux.HandleFunc(listPath, listH)
	mux.HandleFunc(thumbPath, thumbH)
	mux.Handle(filePath, fileH)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	ready, done, s := New(port, mux, tlsConfig)
	<-ready

	fmt.Println("Server listening on: ")
	for _, ip := range web.GetIps() {
		fmt.Printf("https://%v:%v\n", ip, port)
	}
	<-sigint

	ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(time.Second*3))
	defer cancel()

	s.Shutdown(ctx)

	log.Println("Shutting down..")

	select {
	case <-done:
		log.Println("bye")
	case <-ctx.Done():
		log.Println("shutdown request timedout..")
		log.Println("ok..")
	}

}
