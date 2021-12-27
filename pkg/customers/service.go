package customer

import (
	"context"
	// "encoding/json"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrNotFound = errors.New("customer not found")
var ErrInternalServer = errors.New("internal server error")


type Service struct {
	pool *pgxpool.Pool
}

type Customer struct {
						// shit 
	ID		int64		`json:"id"`
	Name	string		`json:"name"`
	Phone	string		`json:"phone"`
	Active	bool		`json:"active"`
	Created time.Time	`json:"created"`
}
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}


func (s *Service) CustomerGetByID(ctx context.Context, customerID int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		SELECT id, name, phone, active, created
		FROM customers 
		WHERE id = $1;
	`, customerID).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	

	if errors.Is(err, pgx.ErrNoRows) {
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

	rows, err := s.pool.Query(ctx, `
		SELECT id, name, phone, active, created
		FROM customers;
	`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

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

	rows, err := s.pool.Query(ctx, `
		SELECT id, name, phone, active, created
		FROM customers
		WHERE active = true;
	`)

	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

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

func (s *Service) CustomerSave(ctx context.Context, customer *Customer) (*Customer, error) {
	var err error

	log.Println("<=====", customer)

	if customer.ID == 0 {
		err = s.pool.QueryRow(ctx, `
			INSERT INTO customers(name, phone)
			VALUES ($1, $2)
			RETURNING id, name, phone, active, created; 
		`, customer.Name, customer.Phone).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)
	} else {
		err = s.pool.QueryRow(ctx, `
			UPDATE customers 
			SET name = $1, phone = $2
			WHERE id = $3
			RETURNING id, name, phone, active, created; 
		`, customer.Name, customer.Phone, customer.ID).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)
	}


	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}

	log.Println("=====>", customer)

	return customer, nil
}


func (s *Service) CustomerRemoveByID(ctx context.Context, id uint64) (int64, error) {
	result, err := s.pool.Query(ctx, `
		DELETE FROM customers
		WHERE id = $1;
	`, id)

	return checkRowsAffetcted(result, err)
}


func (s *Service) CustomerBlockByID(ctx context.Context, id uint64) (customer *Customer, err error) {
	customer = &Customer{}

	log.Println(id)

	err = s.pool.QueryRow(ctx, `
		UPDATE customers
		SET active = false
		WHERE id = $1
		RETURNING id, name, phone, active, created;
	`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("no rows")
		return
	}

	return customer, nil
}

func (s *Service) CustomerUnblockByID(ctx context.Context, id uint64) (customer *Customer, err error) {
	customer = &Customer{}

	err = s.pool.QueryRow(ctx, `
		UPDATE customers
		SET active = true
		WHERE id = $1
		RETURNING id, name, phone, active, created;
	`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("no rows")
		return
	}

	return customer, nil
}
