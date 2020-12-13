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
func (s *Service) Auth(login, password string) bool{
	stateSql := `select login * password from managers where login=$1 and password=$2`

	//Req 
	err := s.db.QueryRow(context.Background(), stateSql, login, password).Scan(&loggin, &password)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}