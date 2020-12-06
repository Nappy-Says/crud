package main

import (
	"go.uber.org/dig"
	"time"
	"context"
	"github.com/Nappy-Says/crud/cmd/app"
	"github.com/Nappy-Says/crud/pkg/customers"
	"net/http"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

func main() {
	host :="0.0.0.0"
	port := "9999"
	dbConnectionString :="postgres://app:pass@localhost:5432/db"
	if err := execute(host, port, dbConnectionString); err != nil{
		log.Print(err)
		os.Exit(1)
	}
}
func execute(host, port, dbConnectionString string) (err error){
	dependencies := []interface{}{
		app.NewServer,
		http.NewServeMux,
		func() (pgxpool.Pool, error){
			connCtx, _ := context.WithTimeout(context.Background(), time.Second)
			return pgxpool.Connect(connCtx, dbConnectionString)
		},
		customers.NewService,
		func(server *app.Server)*http.Server{
			return &http.Server{
				Addr:host+":"+port,
				Handler: server,
			}
		},
	}
	container := dig.New()
	for _, v := range dependencies {
		err = container.Provide(v)
		if err !=nil{
			return err
		}
	}
	err = container.Invok(func(server *app.Server){
		server.Init()
	}
	if err != nil{
		return err
	}
	return container.Invoke(func(server *http.Server) error{
		return server.ListenAndServe()
	})
}
