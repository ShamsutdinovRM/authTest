package handler

import (
	"authTest/model"
	"context"
	"log"
	"net/http"
	"strings"
)

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (b *Repos) JWTMiddleware(next http.Handler) http.Handler {
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

		username, err := ParseToken(headerParse[1])
		if err != nil {
			SendErr(w, http.StatusUnauthorized, err.Error())
			return
		}

		//check redis cache
		err = b.RedisRepository.CheckUser(model.Name{Username: username})
		if err != nil {
			log.Printf("Error, check user: %s", err)
			SendErr(w, http.StatusUnauthorized, "Error, user Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), "username", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
