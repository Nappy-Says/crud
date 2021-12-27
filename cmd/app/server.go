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
	GET = "GET"
	POST = "POST"
	PATCH = "PATCH"
	DELETE = "DELETE"
)


type Server struct {
	mux 		*mux.Router
	// mux 		*http.ServeMux
	customerSvc	*customer.Service
}

func NewServer(mux *mux.Router, customerSvc *customer.Service) *Server {
	return &Server{mux: mux, customerSvc: customerSvc}
}


func (s *Server) ServeHTTP(write http.ResponseWriter, request *http.Request)  {
	s.mux.ServeHTTP(write, request)
}


func (s *Server) Init()  {
	s.mux.HandleFunc("/customers/{id}", 	  s.handleCustomerGetByID).Methods(GET)
	s.mux.HandleFunc("/customers", 			  s.handleCustomerGetAll).Methods(GET)
	s.mux.HandleFunc("/customers.active", 	  s.handleCustomerGetAllActive).Methods(GET)

	s.mux.HandleFunc("/customers", 			  s.handleCustomerSave).Methods(POST)
	s.mux.HandleFunc("/customers/{id}", 	  s.handleCustomerRemoveByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", s.handleCustomerBlockByID).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleCustomerUnblockByID).Methods(DELETE)
}




func (s *Server) handleCustomerGetAll(write http.ResponseWriter, request *http.Request)  {
	items, err := s.customerSvc.CustomerGetAll(request.Context())

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(422), 422)
		return
	}

	data, err := json.Marshal(items)
	
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	respondJSON(write, data)
}


func (s *Server) handleCustomerGetAllActive(write http.ResponseWriter, request *http.Request)  {
	items, err := s.customerSvc.CustomerGetAllActive(request.Context())

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(422), 422)
		return
	}

	data, err := json.Marshal(items)
	
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	respondJSON(write, data)
}


func (s *Server) handleCustomerGetByID(write http.ResponseWriter, request *http.Request)  {
	idParam, ok := mux.Vars(request)["id"]

	if !ok {
		http.Error(write, http.StatusText(400), 400)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}
	
	item, err := s.customerSvc.CustomerGetByID(request.Context(), id)

	if errors.Is(err, customer.ErrNotFound) {
		log.Println(err)
		http.Error(write, http.StatusText(404), 404)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}
	
	data, err := json.Marshal(item)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	respondJSON(write, data)
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
	// item, err = s.customerSvc.CustomerSave(request.Context(), item)

	// if errors.Is(err, customer.ErrNotFound) {
	// 	log.Println(err)
	// 	http.Error(write, http.StatusText(404), 404)
	// 	return
	// }

	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(write, http.StatusText(500), 500)
	// 	return
	// }

	// data, err := json.Marshal(item)

	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(write, http.StatusText(500), 500)
	// 	return
	// }

	// respondJSON(write, data)
}

func (s *Server) handleCustomerRemoveByID(write http.ResponseWriter, request *http.Request)  {
	idParam, ok := mux.Vars(request)["id"]

	if !ok {
		http.Error(write, http.StatusText(400), 400)
		return
	}
	
	id, err := strconv.ParseUint(idParam, 10, 64)
	
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	item, err := s.customerSvc.CustomerRemoveByID(request.Context(), id)

	if errors.Is(err, customer.ErrNotFound) {
		log.Println(err)
		http.Error(write, http.StatusText(404), 404)
		return
	}

	if err != nil {
		http.Error(write, http.StatusText(500), 500)
		return
	}
	
	data, err := json.Marshal(item)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	respondJSON(write, data)
}

func (s *Server) handleCustomerBlockByID(write http.ResponseWriter, request *http.Request)  {
	idParam, ok := mux.Vars(request)["id"]

	if !ok {
		http.Error(write, http.StatusText(400), 400)
		return
	}
	
	id, err := strconv.ParseUint(idParam, 10, 64)
	
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}
	item, err := s.customerSvc.CustomerBlockByID(request.Context(), id)

	if errors.Is(err, customer.ErrNotFound) {
		log.Println(err)
		http.Error(write, http.StatusText(404), 404)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}
	
	data, err := json.Marshal(item)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	respondJSON(write, data)
}


func (s *Server) handleCustomerUnblockByID(write http.ResponseWriter, request *http.Request)  {
	idParam, ok := mux.Vars(request)["id"]

	if !ok {
		http.Error(write, http.StatusText(400), 400)
		return
	}
	
	id, err := strconv.ParseUint(idParam, 10, 64)
	
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}
	
	item, err := s.customerSvc.CustomerUnblockByID(request.Context(), id)

	if errors.Is(err, customer.ErrNotFound) {
		log.Println(err)
		http.Error(write, http.StatusText(404), 404)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}
	
	data, err := json.Marshal(item)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	respondJSON(write, data)
}
