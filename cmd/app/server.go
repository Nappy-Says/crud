package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nappy-says/crud/cmd/app/middleware"
	"github.com/nappy-says/crud/pkg/customers"
	"github.com/nappy-says/crud/pkg/security"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

type Server struct {
	mux         *mux.Router
	customerSvc *customer.Service
	securitySvc *security.Service
}

func NewServer(mux *mux.Router, customerSvc *customer.Service, securitySvc *security.Service) *Server {
	return &Server{mux: mux, customerSvc: customerSvc, securitySvc: securitySvc}
}

func (s *Server) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(write, request)
}

func (s *Server) Init() {
	s.mux.Use(middleware.Basic(s.securitySvc.Auth))
	s.mux.Use(middleware.Logger)

	s.mux.HandleFunc("/customers", 				s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/active", 		s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", 		s.handleGetCustomerById).Methods(GET)
	s.mux.HandleFunc("/customers", 				s.handleCustomerSave).Methods(POST)

	s.mux.HandleFunc("/customers/{id}", 		s.handleDelete).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", 	s.handleUnBlockByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", 	s.handleBlockByID).Methods(POST)

	// api
	s.mux.HandleFunc("/api/customers", 					s.handleCustomerRegistration).Methods(POST)
	s.mux.HandleFunc("/api/customers/token", 			s.handleGenerateToken).Methods(POST)
	s.mux.HandleFunc("/api/customers/token/validate", 	s.handleCustomerValidateToken).Methods(POST)

}

func (s *Server) handleGetCustomerById(writer http.ResponseWriter, request *http.Request) {
	idParam := mux.Vars(request)["id"]
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		log.Println(err)
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.CustomerGetByID(request.Context(), id)
	log.Println(item)

	if errors.Is(err, customer.ErrNotFound) {
		errorWriter(writer, http.StatusNotFound, err)
		return
	}

	if err != nil {
		log.Println(err)
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, item)
}

func (s *Server) handleGetAllCustomers(write http.ResponseWriter, request *http.Request) {
	items, err := s.customerSvc.CustomerGetAll(request.Context())

	if err != nil {
		errorWriter(write, http.StatusInternalServerError, err)
		return
	}

	respondJSON(write, items)
}

func (s *Server) handleGetAllActiveCustomers(write http.ResponseWriter, request *http.Request) {
	items, err := s.customerSvc.CustomerGetAllActive(request.Context())

	if err != nil {
		errorWriter(write, http.StatusInternalServerError, err)
		return
	}

	respondJSON(write, items)
}

func (s *Server) handleBlockByID(write http.ResponseWriter, request *http.Request) {
	idP := mux.Vars(request)["id"]
	id, err := strconv.ParseUint(idP, 10, 64)

	if err != nil {
		errorWriter(write, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.CustomerBlockByID(request.Context(), id)

	if errors.Is(err, customer.ErrNotFound) {
		errorWriter(write, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(write, http.StatusInternalServerError, err)
		return
	}

	respondJSON(write, item)
}

func (s *Server) handleUnBlockByID(write http.ResponseWriter, request *http.Request) {
	idP := mux.Vars(request)["id"]
	id, err := strconv.ParseUint(idP, 10, 64)

	if err != nil {
		errorWriter(write, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.CustomerUnblockByID(request.Context(), id)
	if errors.Is(err, customer.ErrNotFound) {
		errorWriter(write, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(write, http.StatusInternalServerError, err)
		return
	}
	respondJSON(write, item)
}

func (s *Server) handleDelete(write http.ResponseWriter, request *http.Request) {
	idP := mux.Vars(request)["id"]
	id, err := strconv.ParseUint(idP, 10, 64)

	if err != nil {
		errorWriter(write, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.CustomerRemoveByID(request.Context(), id)
	if errors.Is(err, customer.ErrNotFound) {
		errorWriter(write, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(write, http.StatusInternalServerError, err)
		return
	}

	respondJSON(write, item)
}

func (s *Server) handleCustomerSave(write http.ResponseWriter, request *http.Request) {
	var item *customer.Customer

	err := json.NewDecoder(request.Body).Decode(&item)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	s.customerSvc.CustomerSave(request.Context(), item)

	respondJSON(write, item)
}

func (s *Server) handleCustomerRegistration(write http.ResponseWriter, request *http.Request) {
	var userData *customer.Customer

	err := json.NewDecoder(request.Body).Decode(&userData)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	err, id := s.customerSvc.CustomerRegistration(request.Context(), userData)

	if err != nil {
		http.Error(write, http.StatusText(409), 409)
		return
	}

	respondJSON(write, id)
}

func (s *Server) handleCustomerValidateToken(write http.ResponseWriter, request *http.Request) {
	var status 	int
	var text  	string

	var tempToken *struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(request.Body).Decode(&tempToken)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	id, err := s.securitySvc.AuthenticateCustomer(request.Context(), tempToken.Token)

	if err != nil {
		switch err {
			case security.ErrNoSuchUser:
				status = 404
				text = "not Found"
				
			case security.ErrExpireToken:
				status = 400
				text = "expired"	
				
			default:
				status = 500
				text = "Internal Server Error"
		}

		data, err := json.Marshal(
			map[string]interface{}{"status": "fail", "reason": text})

		if err != nil {
			http.Error(write, text, status)
			return
		}

		respondJSON(write, data)
	}

	respond:= make(map[string]interface{})
	respond["status"] = "ok"
	respond["customerId"] = id
	
	respondJSON(write, respond)
}

func (s *Server) handleGenerateToken(write http.ResponseWriter, request *http.Request) {
	var tempAuth *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(request.Body).Decode(&tempAuth); err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	token, err := s.securitySvc.TokenForCustomer(request.Context(), tempAuth.Login, tempAuth.Password)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	respond := map[string]interface{}{"status": http.StatusText(http.StatusOK), "token": token}

	respondJSON(write, respond)
}
