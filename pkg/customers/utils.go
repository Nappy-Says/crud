package customer

import (
	"database/sql"
	"log"
)





func checkRowsAffetcted(result sql.Result, err error) (string, error) {
	if  err != nil {
		log.Println(err)       
		return "", nil
	}

	if status, err := result.RowsAffected(); status <= 0 || err != nil {
		log.Println(err)
		return "", ErrNotFound
	}

	return "ok", nil
}
