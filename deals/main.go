package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"CRM-V1/deals/server"
	pb "CRM-V1/proto/deals"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	port := os.Getenv("DEALS_PORT")
	if port == "" {
		port = "50052"
	}
	addr := fmt.Sprintf(":%s", port)

	contactsAddr := os.Getenv("CONTACTS_ADDR") // ej: "localhost:50051" o "contacts:50051" en K8s
	var contactsConn *grpc.ClientConn
	if contactsAddr != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, contactsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Printf("WARNING: no se pudo conectar a contacts en %s: %v (continuando en modo degradado)", contactsAddr, err)
		} else {
			contactsConn = conn
			log.Printf("Conectado a contacts en %s", contactsAddr)
		}
	} else {
		log.Println("CONTACTS_ADDR no definido; corriendo sin verificación")
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	srv := server.NewServer(contactsConn)
	pb.RegisterDealsServiceServer(grpcServer, srv)

	log.Printf("🚀 deals server listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
