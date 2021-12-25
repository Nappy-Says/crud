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
	// s.mux.HandleFunc("/customers.getAll", 	s.handleCustomerGetAll)
	s.mux.HandleFunc("/customers.getById", 	s.handleCustomerGetByID)
}




// func (s *Server) handleCustomerGetAll(write http.ResponseWriter, request *http.Request)  {
	
// }


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
		http.Error(write, http.StatusText(500), 500)
	}
	
	data, err := json.Marshal(item)

	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(500), 500)
	}

	respondJSON(write, data)
}
