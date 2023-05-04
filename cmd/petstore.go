package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ashwinrs/gos/internal/models"
	"github.com/ashwinrs/gos/pkg/petstore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	// create database connection
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=ashwin dbname=mango port=5432 sslmode=disable",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// i1 := models.InsuranceEntity{Name: "oscar"}
	// i2 := models.InsuranceEntity{Name: "cigna"}
	// i3 := models.InsuranceEntity{Name: "Aetna"}
	// insur := []models.InsuranceEntity{i1, i2, i3}
	// db.Create(&insur)

	// Migrate the schema
	db.AutoMigrate(&models.InsuranceEntity{})

	// Create an instance of our handler which satisfies the generated interface
	petStore := petstore.NewPetStoreHandler(db)

	// This is how you set up a basic chi router
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
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
