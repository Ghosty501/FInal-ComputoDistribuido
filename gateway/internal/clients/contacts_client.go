package clients

import (
	"log"
	"os"
	"time"

	contactsPb "CRM-V1/proto/contacts"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewContactsClient() (contactsPb.ContactsServiceClient, *grpc.ClientConn) {
	addr := os.Getenv("CONTACTS_ADDR")
	if addr == "" {
		addr = "localhost:50051" // valor por defecto para pruebas locales
	}
	ctxDialTimeout := 5 * time.Second
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(ctxDialTimeout))
	if err != nil {
		log.Printf("WARNING: no se pudo conectar a contacts (%s): %v", addr, err)
		return nil, nil
	}
	client := contactsPb.NewContactsServiceClient(conn)
	return client, conn
}
