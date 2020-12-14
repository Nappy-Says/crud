package security

import (
	"log"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)
type Service struct {
	db *pgxpool.Pool
}
func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

//req
func (s *Service) Auth(login, password string) bool {
	sqlStatement := `SELECT login, password FROM managers WHERE login=$1 AND password=$2`
	err := s.db.QueryRow(context.Background(), sqlStatement, login, password).Scan(&login, &password)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}
