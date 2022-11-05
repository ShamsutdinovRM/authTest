package handler

import (
	"authTest/model"
	"context"
	"encoding/json"
	"io/ioutil"
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error read body: %s", err)
			SendErr(w, http.StatusBadRequest, "Invalid field")
			return
		}
		defer r.Body.Close()

		var user model.Name
		if err = json.Unmarshal(body, &user); err != nil {
			log.Printf("Error unmarshal body: %s", err)
			SendErr(w, http.StatusBadRequest, "Invalid field")
			return
		}

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

		if username != user.Username {
			SendErr(w, http.StatusUnauthorized, "User not authorized")
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
