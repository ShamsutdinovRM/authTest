package handler

import (
	"net/http"
	"strings"
)

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			SendErr(w, http.StatusUnauthorized, "Empty auth header")
			return
		}

		headerParse := strings.Split(header, " ")
		if len(headerParse) != 2 {
			SendErr(w, http.StatusUnauthorized, "Invalid auth header")
			return
		}

		_, err := ParseToken(headerParse[1])
		if err != nil {
			SendErr(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
