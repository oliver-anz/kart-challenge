package main

import (
	"backend-challenge/api"
	"backend-challenge/db"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	dbPath := flag.String("db", "data/store.db", "Path to SQLite database")
	flag.Parse()

	database, err := db.New(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	handler := api.NewHandler(database)
	router := handler.SetupRoutes()

	addr := ":" + *port
	fmt.Printf("Server starting on port %s...\n", *port)
	fmt.Printf("API available at http://localhost:%s/api\n", *port)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
