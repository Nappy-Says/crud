package app

import (
	"github.com/Nappy-Says/crud/pkg/security"
	"github.com/Nappy-Says/crud/cmd/app/middleware"
	"encoding/json"
	"errors"
	"log"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/Nappy-Says/crud/pkg/customers"
)
const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)
type Server struct {
	mux         *mux.Router
	customerSvc *customers.Service
	securitySvc *security.Service
}

func NewServer(m *mux.Router, cSvc *customers.Service, sSvc *security.Service) *Server {
	return &Server{
		mux:         m,
		customerSvc: cSvc,
		securitySvc: sSvc,
	}
}
func (s *Server) Init() {
	
	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnBlockByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}", s.handleDelete).Methods(DELETE)
	s.mux.HandleFunc("/api/customers", s.handleSave).Methods(POST)
	s.mux.HandleFunc("/api/customers/token", s.handleCreateToken).Methods(POST)
	s.mux.HandleFunc("/api/customers/token/validate", s.handleValidateToken).Methods(POST)
	s.mux.Use(middleware.Basic(s.securitySvc.Auth))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	items, err := s.customerSvc.All(r.Context())
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, items)
}

func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request) {
	items, err := s.customerSvc.AllActive(r.Context())
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, items)
}

func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.ByID(r.Context(), id)

	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, item)
}

func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.ChangeActive(r.Context(), id, false)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, item)
}

func (s *Server) handleUnBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.ChangeActive(r.Context(), id, true)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, item)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.Delete(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, item)
}

func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	var item *customers.Customer

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	customer, err := s.customerSvc.Save(r.Context(), item)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, customer)
}

func (s *Server) handleCreateToken(w http.ResponseWriter, r *http.Request) {
	var item *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	token, err := s.securitySvc.TokenForCustomer(r.Context(), item.Login, item.Password)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	respondJSON(w, map[string]interface{}{"status": "ok", "token": token})
}

func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	var item *struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	id, err := s.securitySvc.AuthenticateCustomer(r.Context(), item.Token)

	if err != nil {
		respondJSON(w, map[string]interface{}{"status": "fail", "reason": fmt.Sprintf("%v", err)})
		return
	}
	respondJSON(w, map[string]interface{}{"status": "ok", "customerId": id})
}

func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	log.Print(err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}

func respondJSON(w http.ResponseWriter, iData interface{}) {
	data, err := json.Marshal(iData)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func respondJSONWithCode(w http.ResponseWriter, sts int, iData interface{}) {
	data, err := json.Marshal(iData)

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(sts)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}