package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	contactsPb "CRM-V1/proto/contacts"
	pb "CRM-V1/proto/deals"

	"google.golang.org/grpc"
)

type DealServer struct {
	pb.UnimplementedDealsServiceServer
	mu             sync.Mutex
	deals          map[string]*pb.DealResponse
	nextID         int
	contactsClient contactsPb.ContactsServiceClient
}

func NewServer(contactsConn *grpc.ClientConn) *DealServer {
	var c contactsPb.ContactsServiceClient
	if contactsConn != nil {
		c = contactsPb.NewContactsServiceClient(contactsConn)
	}
	return &DealServer{
		deals:          make(map[string]*pb.DealResponse),
		nextID:         1,
		contactsClient: c,
	}
}

func nowString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func (s *DealServer) CreateDeal(ctx context.Context, req *pb.CreateDealRequest) (*pb.DealResponse, error) {
	if s.contactsClient != nil {
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
