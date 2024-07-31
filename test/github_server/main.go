package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mysterion/avrp/internal/utils"
	"github.com/mysterion/avrp/web"
	"github.com/mysterion/avrp/web/dist"
)

func setupEnv() {
	os.Setenv("URL_LATEST_RELEASE", "http://localhost:5678/releases/latest")
}

// port = 0, for a random port
//
// mux can be nil
//
// tlsConfig can be nil
//
// returns [<-ready], [<-done], [*http.Server]
func new(port int, mux http.Handler, tlsConfig *tls.Config) (<-chan bool, <-chan bool, *http.Server) {

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

		err = server.Serve(listener)

		if err != nil && err != http.ErrServerClosed {
			fmt.Println(err.Error())
			panic(err)
		}
	}()

	return ready, done, &server
}

func main() {
	setupEnv()
	mux := http.NewServeMux()
	mux.HandleFunc("/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		var rel dist.Release
		var as dist.Asset
		as.BrowserDownloadUrl = "http://localhost:5678/file/dist-1.3.2.zip"
		rel.Tag = "1.3.2"
		rel.URL = "http://localhost:5678/file/dist-1.3.2.zip"
		rel.Assets = append(rel.Assets, as)
		b, err := json.Marshal(rel)
		log.Println(rel)
		if err != nil {
			w.Write([]byte("Failed to marshal release"))
			return
		}
		w.Write(b)
	})

	filePath := "/file/"
	fileServer := http.FileServer(http.Dir("test/github_server"))
	mux.Handle(filePath, http.StripPrefix(filePath, fileServer))

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	port := 5678
	ready, done, s := new(port, mux, nil)
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
