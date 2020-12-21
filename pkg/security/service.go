package security

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx/v4"
	"encoding/hex"
	"crypto/rand"
	"context"
	"errors"
	"time"
	"log"
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


func (s *Service) TokenForCustomer(ctx context.Context, phone, password string) (string, error) {

	var hash string
	var id int64

	err := s.db.QueryRow(
		ctx,
		"SELECT id, password FROM customers WHERE phone = $1",
		phone).Scan(&id, &hash)
	if err == pgx.ErrNoRows {

		return "", ErrNoSuchUser
	}

	if err != nil {

		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {

		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {

		return "", ErrInternal
	}

	token := hex.EncodeToString(buffer)
	_, err = s.db.Exec(
		ctx,
		"INSERT INTO customers_tokens (token, customer_id) VALUES ($1, $2)",
		token, id)
	if err != nil {

		return "", ErrInternal
	}

	return token, nil
}

func (s *Service) AuthenticateCustomer(ctx context.Context, token string) (int64, error) {
	var id int64
	var expire time.Time
	err := s.db.QueryRow(
		ctx,
		"SELECT customer_id, expire FROM customers_tokens WHERE token=$1",
		token).Scan(&id, &expire)
	if err == pgx.ErrNoRows {

		return 0, ErrNoSuchUser
	}

	if err != nil {

		return 0, ErrInternal
	}

	timeNow := time.Now().Format("2020-07-07 10:36:21")
	timeEnd := expire.Format("2020-07-07 10:36:21")

	if timeNow > timeEnd {

		return 0, ErrExpireToken
	}

	return id, nil
}
