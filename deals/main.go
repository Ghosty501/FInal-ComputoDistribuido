package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	// NOTA: No se necesita "database/sql" ni código de conexión a DB aquí.

	"CRM-V1/deals/server"
	pb "CRM-V1/proto/deals"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("DEALS_PORT")
	if port == "" {
		port = "50052"
	}
	addr := fmt.Sprintf(":%s", port)

	contactsAddr := os.Getenv("CONTACTS_ADDR") // ej: "contacts-service:50051" en K8s
	var contactsConn *grpc.ClientConn
	if contactsAddr != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// El grpc.WithBlock() asegura que el Dial intente conectarse antes de continuar
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

	// CAMBIO CLAVE: Pasamos contactsConn (primero) y nil para la DB (segundo)
	// Ya que no estamos usando la DB real en este servicio por ahora.
	srv := server.NewServer(contactsConn, nil)

	pb.RegisterDealsServiceServer(grpcServer, srv)

	reflection.Register(grpcServer) // Habilitar la reflexión para depuración gRPC

	log.Printf("🚀 deals server listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
