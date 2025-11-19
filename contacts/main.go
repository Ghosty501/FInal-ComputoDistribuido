package main

import (
	"fmt"
	"log"
	"net"

	pb "CRM-V1/proto/contacts"

	"CRM-V1/contacts/server"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	grpcServer := grpc.NewServer()
	contactServer := server.NewServer()

	pb.RegisterContactsServiceServer(grpcServer, contactServer)

	fmt.Println("🚀 Servidor de contactos escuchando en el puerto 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
