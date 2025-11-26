package handlers

import (
	"context"
	"encoding/json"
	"log" // Añadido para logging en listContactsHandler
	"net/http"
	"time"

	contactsPb "CRM-V1/proto/contacts"
	dealsPb "CRM-V1/proto/deals"
)

const requestTimeout = 5 * time.Second

// RegisterRoutes registra rutas HTTP en el mux
func RegisterRoutes(mux *http.ServeMux, contactsClient contactsPb.ContactsServiceClient, dealsClient dealsPb.DealsServiceClient) {
	mux.HandleFunc("/health", healthHandler)
	// Contactos
	mux.HandleFunc("/contacts/create", makeCreateContactHandler(contactsClient))
	mux.HandleFunc("/contacts/list", makeListContactsHandler(contactsClient)) // <-- NUEVA RUTA
	// Deals
	mux.HandleFunc("/deals/create", makeCreateDealHandler(dealsClient))
	mux.HandleFunc("/deals/list", makeListDealsHandler(dealsClient))     // <--- ¡ASEGÚRATE DE QUE ESTA LÍNEA ESTÉ PRESENTE!
	mux.HandleFunc("/deals/health", makeDealsHealthHandler(dealsClient)) // <-- NUEVA RUTA para Health Check
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// --- HANDLERS DE CONTACTS ---

func makeCreateContactHandler(contactsClient contactsPb.ContactsServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { // Asegurar el método POST
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if contactsClient == nil {
			http.Error(w, "contacts client not available", http.StatusServiceUnavailable)
			return
		}
		var req struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Phone string `json:"phone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		resp, err := contactsClient.CreateContact(ctx, &contactsPb.CreateContactRequest{
			Name:  req.Name,
			Email: req.Email,
			Phone: req.Phone,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// NUEVO: Handler para Listar Contactos
func makeListContactsHandler(contactsClient contactsPb.ContactsServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet { // Asegurar el método GET
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if contactsClient == nil {
			http.Error(w, "contacts client not available", http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Asumiendo que ListContacts no requiere parámetros en la solicitud gRPC
		resp, err := contactsClient.ListContacts(ctx, &contactsPb.Empty{})
		if err != nil {
			log.Printf("Error al llamar ListContacts gRPC: %v", err)
			http.Error(w, "Error interno del servidor al listar contactos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Devolvemos la lista de contactos
		json.NewEncoder(w).Encode(resp.Contacts)
	}
}

// --- HANDLERS DE DEALS ---

func makeCreateDealHandler(dealsClient dealsPb.DealsServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { // Asegurar el método POST
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if dealsClient == nil {
			http.Error(w, "deals client not available", http.StatusServiceUnavailable)
			return
		}
		var req struct {
			ContactId string  `json:"contact_id"`
			Title     string  `json:"title"`
			Amount    float64 `json:"amount"`
			Currency  string  `json:"currency"`
			Status    string  `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		resp, err := dealsClient.CreateDeal(ctx, &dealsPb.CreateDealRequest{
			ContactId: req.ContactId,
			Title:     req.Title,
			Amount:    req.Amount,
			Currency:  req.Currency,
			Status:    req.Status,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func makeListDealsHandler(client dealsPb.DealsServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		// Llama al método ListDeals en el microservicio de Deals
		resp, err := client.ListDeals(ctx, &dealsPb.Empty{})
		if err != nil {
			log.Printf("Error al listar deals: %v", err)
			http.Error(w, "Error interno del servicio de Deals", http.StatusInternalServerError)
			return
		}

		// Envía la lista de deals como respuesta JSON
		w.Header().Set("Content-Type", "application/json")

		// NOTA: 'resp.Deals' contiene la estructura de array que queremos
		if err := json.NewEncoder(w).Encode(resp.Deals); err != nil {
			log.Printf("Error al codificar respuesta JSON: %v", err)
			// Ya se escribió un encabezado, es difícil recuperarse aquí
		}
	}
}

// NUEVO: Handler para Health Check de Deals
func makeDealsHealthHandler(dealsClient dealsPb.DealsServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet { // Asegurar el método GET
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if dealsClient == nil {
			http.Error(w, "deals client not available", http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// 1. Llamar al método gRPC HealthCheck (asumiendo que está en el .proto de deals)
		_, err := dealsClient.HealthCheck(ctx, &dealsPb.HealthCheckRequest{})

		// Si hay error, el servicio no está sano
		if err != nil {
			log.Printf("HealthCheck de Deals falló: %v", err)
			http.Error(w, "Deals Service no está sano", http.StatusServiceUnavailable)
			return
		}

		// 2. Éxito: Escribir una respuesta simple
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "deals"})
	}
}
