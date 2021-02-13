package users

// UserRepository interface provides functions for CRUD operations for user entity
type UserRepository interface {
	// CreateUser function provides create db operation for user entity that do not exist yet
	// Returns an error if such occurs during the db query execution or if the password encryption fails
	CreateUser(user *User) error

	// ExistUser function checks if a user with given email exists in the db
	// Returns a boolean value
	ExistUser(email string) bool

	// FindUser function fetches a user with given email and password if such exists in the db
	// Returns the user if such exists or an error otherwise
	FindUser(email, password string) (*User, error)
}
