package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	contactsPb "CRM-V1/proto/contacts"
	dealsPb "CRM-V1/proto/deals"
)

// RegisterRoutes registra rutas HTTP en el mux
func RegisterRoutes(mux *http.ServeMux, contactsClient contactsPb.ContactsServiceClient, dealsClient dealsPb.DealsServiceClient) {
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/contacts/create", makeCreateContactHandler(contactsClient))
	mux.HandleFunc("/deals/create", makeCreateDealHandler(dealsClient))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func makeCreateContactHandler(contactsClient contactsPb.ContactsServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

func makeCreateDealHandler(dealsClient dealsPb.DealsServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
