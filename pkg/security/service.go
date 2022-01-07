package security

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)


type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}


func (s *Service) Auth(login string, password string) (ok bool) {
	err := s.db.QueryRow(context.Background(), `
		select login, password
		from managers
		where login = $1 and password = $2;
	`, login, password).Scan(&login, &password)


	log.Println(err)

	if err != nil {
		return false
	}
	
	return true
}
