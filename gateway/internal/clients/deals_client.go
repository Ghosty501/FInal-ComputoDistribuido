package clients

import (
	"log"
	"os"
	"time"

	dealsPb "CRM-V1/proto/deals"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewDealsClient() (dealsPb.DealsServiceClient, *grpc.ClientConn) {
	addr := os.Getenv("DEALS_ADDR")
	if addr == "" {
		addr = "localhost:50052" // cambiar si en local deals usa otro puerto
	}
	ctxDialTimeout := 5 * time.Second
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(ctxDialTimeout))
	if err != nil {
		log.Printf("WARNING: no se pudo conectar a deals (%s): %v", addr, err)
		return nil, nil
	}
	client := dealsPb.NewDealsServiceClient(conn)
	return client, conn
}
