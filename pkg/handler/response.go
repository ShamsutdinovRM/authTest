package handler

import (
	"authTest/model"
	"authTest/pkg/repository"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
)

type Repos struct {
	Repository      repository.Operations
	RedisRepository repository.OperationsRedis
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

	logUser := model.LogUser{
		Username: userDB.Username,
		Token:    JWTToken,
	}

	loginUser, err := b.RedisRepository.LoginUser(logUser)
	if err != nil {
		log.Printf("Error save data User: %s", err)
		SendErr(w, http.StatusInternalServerError, "Error login to service")
		return
	}

	SendOK(w, http.StatusOK, loginUser)
}

func (b *Repos) LogOut(w http.ResponseWriter, r *http.Request) {
	user := fmt.Sprint(r.Context().Value("username"))

	err := b.RedisRepository.LogoutUser(model.Name{Username: user})
	if err != nil {
		log.Printf("Error logout User: %s", err)
		SendErr(w, http.StatusInternalServerError, "Error logout to service")
		return
	}

	SendOK(w, http.StatusOK, "Goodbye")
}

func (b *Repos) Hello(w http.ResponseWriter, r *http.Request) {
	SendOK(w, http.StatusOK, "Hello From Wrong Side Of Heaven "+fmt.Sprint(r.Context().Value("username")))
}
