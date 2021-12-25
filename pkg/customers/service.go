package customer

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

var ErrNotFound = errors.New("customer not found")
var ErrInternalServer = errors.New("internal server error")


type Service struct {
	db *sql.DB
}

type Customer struct {
	ID		int64
	Name	string
	Phone	string
	Active	bool
	Created time.Time
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}


func (s *Service) CustomerGetByID(ctx context.Context, customerID int64) (*Customer, error) {
	item := &Customer{}

	err := s.db.QueryRowContext(ctx, `
		SELECT c.id, c.name, c.phone, c.active, c.created
		FROM customers c
		WHERE c.id = $1
	`, customerID).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	

	if err != nil {
		log.Println(err)	
		return nil, ErrInternalServer
	}

	return item, nil
}


// func (s *Service) CustomerGetAll(ctx context.Context) (*[]Customer, error) {
// 	items := &Customer{}

	

	
// }
