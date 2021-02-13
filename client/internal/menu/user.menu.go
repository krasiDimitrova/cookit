package menu

import (
	"client/internal/apis"
	"fmt"
	"golang.org/x/term"
	"net/http"
	"syscall"
)

type UserMenu struct {
	UserApi *apis.UserApi
}

var userMenu *UserMenu

func GetUserMenu() *UserMenu {
	if userMenu == nil {
		userMenu = &UserMenu{UserApi: apis.GetUserApi()}
	}

	return userMenu
}

// GetAccessToken function handles stdin/stdout operations for registration and login
// Uses the user api for the requests and to obtain the user token
func (um *UserMenu) GetAccessToken() []*http.Cookie {
	fmt.Print("Registration (1) or Log in (2): ")

	var command int
	fmt.Scan(&command)

	if command == 1 {
		err := um.register()
		if err != nil {
			fmt.Println("Failed to register")
			return um.GetAccessToken()
		}
	}

	return um.logIn()

}

func (um *UserMenu) register() error {
	fmt.Println("Register")

	var email, password, name string
	var err error
	for {
		email, password, name, err = getCredentials(true)
		if err == nil {
			break
		}
	}

	var user = apis.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	err = um.UserApi.RegisterUser(user)
	return err
}

func (um *UserMenu) logIn() []*http.Cookie {
	fmt.Println("Log in")

	var email, password string
	var err error
	for {
		email, password, _, err = getCredentials(false)
		if err == nil {
			break
		}
	}

	cookies := um.UserApi.LogIn(email, password)
	return cookies
}

func getCredentials(full bool) (string, string, string, error) {
	fmt.Print("Enter user email: ")
	var email string
	_, err := fmt.Scan(&email)
	if err != nil {
		fmt.Println("Error parsing email. Try again! ")
		return "", "", "", err
	}

	fmt.Print("Enter Password: ")
	//Used when testing through the intellij terminal since the term.ReadPassword is not working there
	//var password string
	//_, err = fmt.Scan(&password)
	//if err != nil {
	//	fmt.Println("Error parsing password. Try again! ")
	//	return "", "", "", err
	//}

	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Error parsing password. Try again! ")
		return "", "", "", err
	}
	password := string(bytePassword)

	var name = ""
	if full {
		fmt.Print("Enter name: ")
		_, err = fmt.Scan(&name)
		if err != nil {
			fmt.Println("Error parsing name. Try again! ")
			return "", "", "", err
		}
	}

	return email, password, name, nil
}
