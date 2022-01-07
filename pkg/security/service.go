package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNoSuchUser      = errors.New("no such user")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInternal        = errors.New("internal error")
	ErrExpireToken     = errors.New("token expired")
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



func (s *Service) TokenForCustomer(ctx context.Context, phone, password string) (string, error) {
	var id int64
	var hash string

	err := s.db.QueryRow(ctx, `
		SELECT id, password FROM customers WHERE phone = $1
		`, phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return "", ErrNoSuchUser
	}

	if err != nil {
		log.Println(err)
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		log.Println(err)
		return "", ErrInvalidPassword
	}

	buf := make([]byte, 1024)
	n, err := rand.Read(buf)

	if n != len(buf) || err != nil {
		log.Println(err)
		return "", ErrInternal
	}

	token := hex.EncodeToString(buf)
	_, err = s.db.Exec(ctx,`
		INSERT INTO customers_tokens (token, customer_id) VALUES ($1, $2)
		`,token, id)

	if err != nil {
		log.Println(err)
		return "", ErrInternal
	}

	return token, nil
}

func (s *Service) AuthenticateCustomer(ctx context.Context, token string) (int64, error) {
	var id int64
	var expire time.Time

	err := s.db.QueryRow(ctx,`
			SELECT customer_id, expire 
			FROM customers_tokens 
			WHERE token = $1
		`,token).Scan(&id, &expire)

	if err == pgx.ErrNoRows {
		return 0, ErrNoSuchUser
	}

	if err != nil {
		return 0, ErrInternal
	}

	timeNow := time.Now().Format("2020-01-02 01:02:03")
	timeEnd := expire.Format("2020-01-02 05:06:07")

	if timeNow > timeEnd {
		return 0, ErrExpireToken
	}

	return id, nil
}
