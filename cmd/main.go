package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/Nappy-Says/crud/cmd/app"
	"github.com/Nappy-Says/crud/pkg/customers"
	"github.com/Nappy-Says/crud/pkg/security"
	"go.uber.org/dig"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dbConnectionString := "postgres://app:pass@localhost:5432/db"
	if err := execute(host, port, dbConnectionString); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(host, port, dbConnectionString string) (err error) {

	dependencies := []interface{}{
			connCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(connCtx, dbConnectionString)
		},
			return &http.Server{
				Addr:    host + ":" + port,
				Handler: server,
			}
		},
	}

	container := dig.New()
	for _, v := range dependencies {
		err = container.Provide(v)
		if err != nil {
			return err
		}
	}

	err = container.Invoke(func(server *app.Server) {
		server.Init()
	})
	if err != nil {
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}
