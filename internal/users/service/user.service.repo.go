package service

import (
	"database/sql"
	"errors"
	"github.com/krasimiraMilkova/cookit/internal/db"
	"github.com/krasimiraMilkova/cookit/pkg/users"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserRepository struct {
	*sql.DB
}

func GetUsersRepository() users.UserRepository {
	return &UserRepository{db.Get()}
}

func (userRepository *UserRepository) CreateUser(user *users.User) error {
	if user.Email == "" || user.Password == "" || user.Name == "" {
		return errors.New("user cannot have empty fields")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return errors.New("password encryption failed")
	}

	user.Password = string(pass)

	_, err = userRepository.Exec("insert into users(name,email,password)values(?,?,?)", user.Name, user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func (userRepository *UserRepository) ExistUser(email string) bool {
	user := &users.User{}

	if email == "" {
		return false
	}

	row := userRepository.QueryRow("select id from users where email = ?", email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func (userRepository *UserRepository) FindUser(email, password string) (*users.User, error) {
	user := &users.User{}

	if email == "" || password == "" {
		return nil, errors.New("email or password cannot be empty")
	}

	row := userRepository.QueryRow("select id,name,email,password from users where email = ?", email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password for email")
	}

	return user, nil
}
