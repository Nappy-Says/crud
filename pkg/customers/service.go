package customer

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var ErrNotFound = errors.New("customer not found")
var ErrInternalServer = errors.New("internal server error")


type Service struct {
	db *sql.DB
}

type Customer struct {
						// shit 
	ID		int64		`json:"id"`
	Name	string		`json:"name"`
	Phone	string		`json:"phone"`
	Active	bool		`json:"active"`
	Created time.Time	`json:"created"`
}
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}


func (s *Service) CustomerGetByID(ctx context.Context, customerID int64) (*Customer, error) {
	item := &Customer{}

	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, phone, active, created
		FROM customers 
		WHERE id = $1;
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


func (s *Service) CustomerGetAll(ctx context.Context) ([]*Customer, error) {
	// create an empty instance for further filling
	items := make([]*Customer, 0)

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, phone, active, created
		FROM customers;
	`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer func ()  {
		// Check close status
		if cerr := rows.Close(); cerr != nil {
			log.Println(err)
		}
	} ()

	// fill instance 
	for rows.Next() {
		tempItem := &Customer{}

		// scan each row
		err := rows.Scan(&tempItem.ID, &tempItem.Name, &tempItem.Phone, &tempItem.Active, &tempItem.Created)

		if err != nil {
			log.Println(err)
			continue
		}
		items = append(items, tempItem)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return items, nil
}

func (s *Service) CustomerGetAllActive(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, phone, active, created
		FROM customers
		WHERE active = true;
	`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer func ()  {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	} ()

	for rows.Next() {
		tempItem := &Customer{}

		err := rows.Scan(&tempItem.ID, &tempItem.Name, &tempItem.Phone, &tempItem.Active, &tempItem.Created)

		if err != nil {
			log.Println(err)
			continue
		}

		items = append(items, tempItem)
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return items, nil
}

func (s *Service) CustomerSave(ctx context.Context, id uint64, name string, phone string) (int64, error) {
	var result 	sql.Result
	var err 	error

	if id == 0 {
		result, err = s.db.ExecContext(ctx, `
			INSERT INTO customers(name, phone)
			VALUES ($1, $2)
		`, name, phone)
	} else {
		result, err = s.db.ExecContext(ctx, `
			UPDATE customers 
			SET name = $1, phone = $2
			WHERE id = $3; 
		`, name, phone, id)
	}

	return checkRowsAffetcted(result, err)
}


func (s *Service) CustomerRemoveByID(ctx context.Context, id uint64) (int64, error) {
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM customers
		WHERE id = $1;
	`, id)

	return checkRowsAffetcted(result, err)
}


func (s *Service) CustomerBlockByID(ctx context.Context, id uint64) (int64, error) {
	result, err := s.db.ExecContext(ctx, `
		UPDATE customers
		SET active = false
		WHERE id = $1
	`, id)

	return checkRowsAffetcted(result, err)
}

func (s *Service) CustomerUnblockByID(ctx context.Context, id uint64) (int64, error) {
	result, err := s.db.ExecContext(ctx, `
		UPDATE customers
		SET active = true
		WHERE id = $1
	`, id)

	return checkRowsAffetcted(result, err)
}
