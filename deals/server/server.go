package server

import (
	"context"
	"database/sql" // Necesario para definir el tipo de conexión a DB
	"fmt"
	"log"
	"sync"
	"time"

	contactsPb "CRM-V1/proto/contacts"
	pb "CRM-V1/proto/deals"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DealServer struct {
	pb.UnimplementedDealsServiceServer
	mu             sync.Mutex
	deals          map[string]*pb.DealResponse
	nextID         int
	contactsClient contactsPb.ContactsServiceClient
	DB             *sql.DB // AÑADIDO: Campo para la conexión a la DB
}

// NewServer ahora recibe la conexión a la DB como segundo argumento
func NewServer(contactsConn *grpc.ClientConn, dbConn *sql.DB) *DealServer {
	var c contactsPb.ContactsServiceClient
	if contactsConn != nil {
		c = contactsPb.NewContactsServiceClient(contactsConn)
	}
	return &DealServer{
		deals:          make(map[string]*pb.DealResponse),
		nextID:         1,
		contactsClient: c,
		DB:             dbConn, // ASIGNADO: Se guarda la conexión (o nil, si no se usa DB)
	}
}

func nowString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func (s *DealServer) CreateDeal(ctx context.Context, req *pb.CreateDealRequest) (*pb.DealResponse, error) {
	if s.contactsClient != nil {
		// Validación de contacto existente
		_, err := s.contactsClient.GetContact(ctx, &contactsPb.GetContactRequest{Id: req.ContactId})
		if err != nil {
			return nil, fmt.Errorf("contacto %s no encontrado: %w", req.ContactId, err)
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("%d", s.nextID)
	s.nextID++

	deal := &pb.DealResponse{
		Id:        id,
		ContactId: req.ContactId,
		Title:     req.Title,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Status:    req.Status,
		CreatedAt: nowString(),
		UpdatedAt: nowString(),
	}

	s.deals[id] = deal
	fmt.Printf("Deal creado: %+v\n", deal)
	return deal, nil
}

func (s *DealServer) GetDeal(ctx context.Context, req *pb.GetDealRequest) (*pb.DealResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	d, ok := s.deals[req.Id]
	if !ok {
		return nil, fmt.Errorf("deal %s no encontrado", req.Id)
	}
	return d, nil
}

func (s *DealServer) ListDeals(ctx context.Context, _ *pb.Empty) (*pb.DealsList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	resp := &pb.DealsList{}
	for _, d := range s.deals {
		resp.Deals = append(resp.Deals, d)
	}
	return resp, nil
}

func (s *DealServer) UpdateDealStatus(ctx context.Context, req *pb.UpdateDealStatusRequest) (*pb.DealResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	d, ok := s.deals[req.Id]
	if !ok {
		return nil, fmt.Errorf("deal %s no encontrado", req.Id)
	}
	d.Status = req.Status
	d.UpdatedAt = nowString()
	return d, nil
}

func (s *DealServer) DeleteDeal(ctx context.Context, req *pb.DeleteDealRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.deals[req.Id]
	if !ok {
		return &pb.DeleteResponse{Success: false}, fmt.Errorf("deal %s no encontrado", req.Id)
	}
	delete(s.deals, req.Id)
	return &pb.DeleteResponse{Success: true}, nil
}

// Implementación de HealthCheck
func (s *DealServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	// Si la DB es nil (como será en tu main.go), el servidor siempre reporta estar sano.
	if s.DB != nil {
		// Si el DB no es nil, hacemos la verificación real de la conexión
		if err := s.DB.PingContext(ctx); err != nil {
			log.Printf("DB Ping failed for Deals: %v", err)
			// Devolvemos un error gRPC que el Gateway puede capturar (codes.Unavailable)
			return nil, status.Errorf(codes.Unavailable, "database not ready: %v", err)
		}
	}

	// Si la DB es nil o el ping fue exitoso, el servicio gRPC está vivo.
	return &pb.HealthCheckResponse{Success: true}, nil
}
