package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/nappy-says/crud/pkg/customers"
)


type Server struct {
	mux 		*http.ServeMux
	customerSvc	*customer.Service
}

func NewServer(mux *http.ServeMux, customerSvc *customer.Service) *Server {
	return &Server{mux: mux, customerSvc: customerSvc}
}


func (s *Server) ServeHTTP(write http.ResponseWriter, request *http.Request)  {
	s.mux.ServeHTTP(write, request)
}


func (s *Server) Init()  {
	s.mux.HandleFunc("/customers.getById", 		s.handleCustomerGetByID)
	s.mux.HandleFunc("/customers.getAll", 		s.handleCustomerGetAll)
	s.mux.HandleFunc("/customers.getAllActive", s.handleCustomerGetAllActive)

	s.mux.HandleFunc("/customers.save", 		s.handleCustomerSave)
	s.mux.HandleFunc("/customers.removeById", 	s.handleCustomerRemoveByID)
	s.mux.HandleFunc("/customers.blockById", 	s.handleCustomerBlockByID)
	s.mux.HandleFunc("/customers.unblockById", 	s.handleCustomerUnblockByID)
}




func (s *Server) handleCustomerGetAll(write http.ResponseWriter, request *http.Request)  {
	items, err := s.customerSvc.CustomerGetAll(request.Context())

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(422), 422)
	}

	data, err := json.Marshal(items)
	
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
	}

	respondJSON(write, data)
}


func (s *Server) handleCustomerGetAllActive(write http.ResponseWriter, request *http.Request)  {
	items, err := s.customerSvc.CustomerGetAllActive(request.Context())

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(422), 422)
	}

	data, err := json.Marshal(items)
	
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
	}

	respondJSON(write, data)
}


func (s *Server) handleCustomerGetByID(write http.ResponseWriter, request *http.Request)  {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
	}
	
	item, err := s.customerSvc.CustomerGetByID(request.Context(), id)

	if errors.Is(err, customer.ErrNotFound) {
		log.Println(err)
		http.Error(write, http.StatusText(404), 404)
	}

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
	}
	
	data, err := json.Marshal(item)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
	}

	respondJSON(write, data)
}


func (s *Server) handleCustomerSave(write http.ResponseWriter, request *http.Request)  {
	idParam := request.PostFormValue("id")
	name := request.Form.Get("name")
	phone := request.Form.Get("phone")

	id, err := strconv.ParseUint(idParam, 10, 64)
	
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	item, err := s.customerSvc.CustomerSave(request.Context(), id, name, phone)

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

func (s *Server) handleCustomerRemoveByID(write http.ResponseWriter, request *http.Request)  {
	idParam := request.URL.Query().Get("id")

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
	idParam := request.URL.Query().Get("id")

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


func (s *Server) handleCustomerUnblockByID(write http.ResponseWriter, request *http.Request)  {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
	}
	
	item, err := s.customerSvc.CustomerUnblockByID(request.Context(), id)

	if errors.Is(err, customer.ErrNotFound) {
		log.Println(err)
		http.Error(write, http.StatusText(404), 404)
	}

	if err != nil {
		http.Error(write, http.StatusText(500), 500)
	}
	
	data, err := json.Marshal(item)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
	}

	respondJSON(write, data)
}
