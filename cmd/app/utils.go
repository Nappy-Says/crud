package app

import (
	"log"
	"net/http"
)


func respondJSON(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(data)
	if err != nil {
		log.Print(err)
	}
}
