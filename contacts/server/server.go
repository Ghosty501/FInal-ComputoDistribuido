package server

import (
	"context"
	"fmt"
	"sync"

	pb "CRM-V1/proto/contacts"
)

// ContactServer implementa el servicio gRPC definido en contacts.proto
type ContactServer struct {
	pb.UnimplementedContactsServiceServer
	mu       sync.Mutex
	contacts map[string]*pb.ContactResponse
	nextID   int
}

// NewServer crea una nueva instancia del servidor
func NewServer() *ContactServer {
	return &ContactServer{
		contacts: make(map[string]*pb.ContactResponse),
		nextID:   1,
	}
}

// CreateContact crea un nuevo contacto
func (s *ContactServer) CreateContact(ctx context.Context, req *pb.CreateContactRequest) (*pb.ContactResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("%d", s.nextID)
	s.nextID++

	contact := &pb.ContactResponse{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	s.contacts[id] = contact

	fmt.Printf("Contacto creado: %+v\n", contact)
	return contact, nil
}

// GetContact devuelve un contacto por ID
func (s *ContactServer) GetContact(ctx context.Context, req *pb.GetContactRequest) (*pb.ContactResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	contact, exists := s.contacts[req.Id]
	if !exists {
		return nil, fmt.Errorf("contacto no encontrado")
	}

	return contact, nil
}

// ListContacts devuelve todos los contactos
func (s *ContactServer) ListContacts(ctx context.Context, _ *pb.Empty) (*pb.ContactsList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := &pb.ContactsList{}
	for _, c := range s.contacts {
		list.Contacts = append(list.Contacts, c)
	}

	return list, nil
}

// DeleteContact elimina un contacto por ID
func (s *ContactServer) DeleteContact(ctx context.Context, req *pb.DeleteContactRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.contacts[req.Id]
	if !exists {
		return &pb.DeleteResponse{Success: false}, fmt.Errorf("contacto no encontrado")
	}

	delete(s.contacts, req.Id)
	return &pb.DeleteResponse{Success: true}, nil
}
