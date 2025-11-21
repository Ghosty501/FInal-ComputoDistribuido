package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"CRM-V1/gateway/internal/clients"
	"CRM-V1/gateway/internal/handlers"
)

func main() {
	// crear clientes gRPC
	contactsClient, contactsConn := clients.NewContactsClient()
	if contactsConn != nil {
		defer contactsConn.Close()
	}
	dealsClient, dealsConn := clients.NewDealsClient()
	if dealsConn != nil {
		defer dealsConn.Close()
	}

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, contactsClient, dealsClient)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Gateway HTTP escuchando en :%s (CONTACTS_ADDR=%s, DEALS_ADDR=%s)", port, os.Getenv("CONTACTS_ADDR"), os.Getenv("DEALS_ADDR"))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("gateway ListenAndServe: %v", err)
	}
}
