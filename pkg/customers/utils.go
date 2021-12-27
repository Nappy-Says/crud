package customer

import (
	"log"

	"github.com/jackc/pgx/v4"
)





func checkRowsAffetcted(result pgx.Rows, err error) (int64, error) {
	if  err != nil {
		log.Println(err)       
		return 0, err
	}

	// if status := result.RawValues(); status <= 0 || err != nil {
		// log.Println(err)
		// return 0, ErrNotFound
	// }

	status := result.RawValues();

	log.Println("----", result, status)
	
	return 1, nil
}
