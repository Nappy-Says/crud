package app

import (
	"encoding/json"
	"log"
	"net/http"
)


func respondJSON(w http.ResponseWriter, data interface{}) {
	item, err := json.Marshal(data)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(item)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}


func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	log.Print(err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}
