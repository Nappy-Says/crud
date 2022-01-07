package middleware

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"

	// "github.com/nappy-says/crud/pkg/security"
)

func Basic(auth func(string, string) bool) func(handler http.Handler) http.Handler {

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {
			login, password, err := dataExtraction(request)

			if err != nil {
				log.Println(err)
				http.Error(write, http.StatusText(401), 401)
				return
			}

			if !auth(login, password) {
				http.Error(write, http.StatusText(401), 401)
				return
			}

			handler.ServeHTTP(write, request)
		})
	}
}


func dataExtraction(request *http.Request) (string, string, error) {
	authHeader := strings.SplitN(request.Header.Get("Authorization"), " ", 2)

	if len(authHeader) != 2 || authHeader[0] != "Basic" {
		return "", "", errors.New("wrong header data")
	}

	base64, _ := base64.StdEncoding.DecodeString(authHeader[1])
	res := strings.SplitN(string(base64), ":", 2)

	if len(res) != 2 {
		return "", "", errors.New("wrong header data")
	}

	return res[0], res[1], nil
}
