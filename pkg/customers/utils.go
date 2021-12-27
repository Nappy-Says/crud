package customer

import (
	"database/sql"
	"log"
)





func checkRowsAffetcted(result sql.Result, err error) (int64, error) {
	if  err != nil {
		log.Println(err)       
		return 0, err
	}

	if status, err := result.RowsAffected(); status <= 0 || err != nil {
		log.Println(err)
		return 0, ErrNotFound
	}

	log.Println("----", result)
	
	return 1, nil
}
