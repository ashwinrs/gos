package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ashwinrs/gos/pkg/petstore"
	"github.com/go-chi/chi/v5"
)

func main() {
	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	// Create an instance of our handler which satisfies the generated interface
	petStore := petstore.NewPetStore()

	// This is how you set up a basic chi router
	r := chi.NewRouter()

	// We now register our petStore above as the handler for the interface
	petstore.HandlerFromMux(petStore, r)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
	}

	fmt.Println("Server is listening on port:", *port)
	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
