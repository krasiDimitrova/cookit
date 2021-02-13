package service

import (
	"encoding/json"
	"github.com/krasimiraMilkova/cookit/internal/users/auth"
	"github.com/krasimiraMilkova/cookit/pkg/users"
	"log"
	"net/http"
	"time"
)

type UserService struct {
	UserRepository    users.UserRepository
	UserAuthenticator users.UserAuthenticator
}

var usersService *UserService

func Get() *UserService {
	if usersService == nil {
		usersService = &UserService{UserRepository: GetUsersRepository(), UserAuthenticator: auth.GetAuthenticator()}
	}

	return usersService
}

func (us *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &users.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		log.Print("Error occurred when decoding user payload", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existUser := us.UserRepository.ExistUser(user.Email)

	if existUser {
		log.Print("User already exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := us.UserRepository.CreateUser(user); err != nil {
		log.Print("Error occurred when creating a user", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (us *UserService) Login(w http.ResponseWriter, r *http.Request) {
	user := &users.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		log.Print("Error occurred when decoding user payload", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundUser, err := us.UserRepository.FindUser(user.Email, user.Password)

	if err != nil {
		log.Print("Error occurred when fetching user ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	token, err := us.UserAuthenticator.GenerateTokenForUser(foundUser)

	if err != nil {
		log.Print("Error occurred when generating user token", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    auth.TokenName,
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(auth.Expiration),
	})
}
