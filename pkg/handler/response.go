package handler

import (
	"authTest/model"
	"authTest/pkg/repository"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
)

type Repos struct {
	Repository repository.Operations
}

func (b *Repos) SignUp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error read body: %s", err)
		SendErr(w, http.StatusBadRequest, "Invalid field")
		return
	}
	defer r.Body.Close()

	var user model.User
	if err = json.Unmarshal(body, &user); err != nil {
		log.Printf("Error unmarshal body: %s", err)
		SendErr(w, http.StatusBadRequest, "Invalid field")
		return
	}

	user.Password = getHash([]byte(user.Password))

	sign, err := b.Repository.CreateUser(user)
	if err != nil {
		SendErr(w, http.StatusBadRequest, err.Error())
		return
	}

	SendOK(w, http.StatusOK, sign.Username)
}

func (b *Repos) LogIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error read body: %s", err)
		SendErr(w, http.StatusBadRequest, "Invalid field")
		return
	}
	defer r.Body.Close()

	var user model.User
	if err = json.Unmarshal(body, &user); err != nil {
		log.Printf("Error unmarshal body: %s", err)
		SendErr(w, http.StatusBadRequest, "Invalid field")
		return
	}

	userDB, err := b.Repository.LoginUser(user)
	if err != nil {
		log.Printf("Error check user: %s", err)
		SendErr(w, http.StatusBadRequest, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if err != nil {
		log.Printf("Wrong password: %s", err)
		SendErr(w, http.StatusBadRequest, "Wrong password or Username")
		return
	}

	JWTToken, err := GenerateJWT(userDB.Username)
	if err != nil {
		log.Printf("Error create token: %s", err)
		SendErr(w, http.StatusInternalServerError, "Error login to service")
		return
	}

	LogUser := model.LogUser{
		Username: userDB.Username,
		Token:    JWTToken,
	}

	SendOK(w, http.StatusOK, LogUser)
}

func (b *Repos) LogOut(w http.ResponseWriter, r *http.Request) {

}

func (b *Repos) Hello(w http.ResponseWriter, r *http.Request) {
	SendOK(w, http.StatusOK, "Hello From Wrong Side Of Heaven")
}
