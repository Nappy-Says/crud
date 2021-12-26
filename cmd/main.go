package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/nappy-says/crud/cmd/server"
	"github.com/nappy-says/crud/pkg/customers"
)


func main() {
	host := "0.0.0.0"
	port := "9999"
	pgxDNS := "postgres://postgres:passs@localhost/db"


	if err := execute(host, port, pgxDNS); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}


func execute(host string, port string, dns string) (err error) {
	db, err := sql.Open("pgx", dns)

	if err != nil {
		log.Println(err)
		return err
	}

	defer func ()  {
		if cerr := db.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
		}

		log.Println(err)
	} ()


	mux := http.NewServeMux()
	customerSvc := customers.NewService(db)

	server := app.NewServer(mux, customerSvc)
	server.Init()

	srv := &http.Server{
		Addr: net.JoinHostPort(host, port),
		Handler: server,
	}

	log.Println("============| Start server |============")

	return srv.ListenAndServe()
}

