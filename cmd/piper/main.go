package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/letung3105/piper/pkg/hub"
)

func main() {
	crt := flag.String("crt", "./keys/server/servercert.pem", "server certificate")
	key := flag.String("key", "./keys/server/serverkey.pem", "server key")
	usersCreds := flag.String("users", "./.creds.json", "Users' credentials")
	binary := flag.String("i", "python3", "name of interpreter")
	script := flag.String("f", "./scripts/main.py", "name of script file")
	backPort := flag.String("back-port", "4433", "backend port to listen")
	flag.Parse()

	// users login infomation
	creds, err := ioutil.ReadFile(*usersCreds)
	if err != nil {
		log.Fatalf("could not read users file; got %v", err)
	}
	var users map[string]*hub.UserInfo
	if err := json.Unmarshal(creds, &users); err != nil {
		log.Fatalf("could not parse users creds; got %v", err)
	}

	// Create and start broadcasting hub
	h := hub.New()
	go h.Run()
	go h.BroadcastScript(*binary, *script)

	// Routing for HTTP connection
	r := mux.NewRouter()
	// Serve index page on all unhandled routes
	r.Handle("/data", h)
	r.Handle("/control", h.Control()).Methods("POST")

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", *backPort),
		Handler:      r,
		TLSConfig:    nil, // TODO: configure TLS
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Printf("serving backend on port %s", *backPort)
	log.Fatal(srv.ListenAndServeTLS(*crt, *key))
}
