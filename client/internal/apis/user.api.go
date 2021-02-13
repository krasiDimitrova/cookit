// Package apis provides clients for communication with the server
package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UserApi struct {
	*http.Client
}

func GetUserApi() *UserApi {
	return &UserApi{&http.Client{}}
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser function sends registration request to the server for the given user information
// Returns an error if such occurs
func (ua *UserApi) RegisterUser(user User) error {
	payloadBuffer := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuffer).Encode(user)

	if err != nil {
		return err
	}

	request, _ := http.NewRequest("POST", serverUrl+"/register", payloadBuffer)
	response, err := ua.Client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if !(response.StatusCode >= 200 && response.StatusCode < 300) {
		return errors.New("failed to register")
	}

	return nil
}

// RegisterUser function sends login request to the server for the given user information
// Returns the obtained cookies
func (ua *UserApi) LogIn(email string, password string) []*http.Cookie {
	user := User{}
	user.Email = email
	user.Password = password

	payloadBuffer := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuffer).Encode(user)

	if err != nil {
		fmt.Println("Failed to login!")
		return nil
	}

	request, _ := http.NewRequest("POST", serverUrl+"/login", payloadBuffer)
	response, err := ua.Client.Do(request)

	if err != nil {
		fmt.Println("Failed to login!")
		return nil
	}

	defer response.Body.Close()

	if !(response.StatusCode >= 200 && response.StatusCode < 300) {
		fmt.Println("Failed to login!")
	}

	return response.Cookies()
}
