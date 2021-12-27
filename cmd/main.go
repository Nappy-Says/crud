package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nappy-says/crud/cmd/app"
	"github.com/nappy-says/crud/pkg/customers"
	"go.uber.org/dig"
)


func main() {
	host := "0.0.0.0"
	port := "9991"
	pgxDNS := "postgres://postgres:passs@localhost/db"


	if err := execute(host, port, pgxDNS); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}


func execute(host string, port string, dns string) (err error) {

	deps := []interface{}{
		app.NewServer,
		mux.NewRouter,
		// http.NewServeMux,
 		func() (*pgxpool.Pool, error) {
			ctx, _ := context.WithTimeout(context.Background(), time.Second * 5)
			return pgxpool.Connect(ctx, dns)
		},
		customer.NewService,
		func (server *app.Server) *http.Server {
			return &http.Server{
				Addr: net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}
	
	container := dig.New()

	for _, dep := range deps {
		err := container.Provide(dep)

		if err != nil {
			return err
		}
	}

	err = container.Invoke(func (server *app.Server)  {
		server.Init()
	})

	if err != nil {
		return err
	}

	log.Println("============| Start server |============")

	return container.Invoke(func (server *http.Server) error {
		return server.ListenAndServe()
	})
}

