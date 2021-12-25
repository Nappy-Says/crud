package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/nappy-says/crud/cmd/app"
	"github.com/nappy-says/crud/pkg/customer"
)


func main() {
	host := "0.0.0.0"
	port := "9999"
	pgxDNS := "postgres://app:pass@localhost/db"


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
	customerSvc := customer.NewService(db)
	
	server := app.NewServer()

}

