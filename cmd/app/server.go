package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nappy-says/crud/pkg/customers"
)

const (
	GET		= "GET"
	POST	= "POST"
	DELETE 	= "DELETE"
)


type Server struct {
	mux     	*mux.Router
	customerSvc	*customer.Service
}

func NewServer(mux *mux.Router, customerSvc *customer.Service) *Server {
	return &Server{mux: mux, customerSvc: customerSvc}
}


func (s *Server) ServeHTTP(write http.ResponseWriter, request *http.Request)  {
	s.mux.ServeHTTP(write, request)
}


func (s *Server) Init() {
	s.mux.HandleFunc("/customers", 				s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/active", 		s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", 		s.handleGetCustomerById).Methods(GET)

	s.mux.HandleFunc("/customers/{id}", 		s.handleDelete).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", 	s.handleUnBlockByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", 	s.handleBlockByID).Methods(POST)
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

func (s *Server) handleCustomerSave(write http.ResponseWriter, request *http.Request)  {
	var item *customer.Customer

	err := json.NewDecoder(request.Body).Decode(&item)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	s.customerSvc.CustomerSave(request.Context(), item)
}
