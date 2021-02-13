// Package users provides handlers, db operations and models for creating and fetching users
// as well as methods for JWT authentication
package users

import "net/http"

// UserService interface provides handlers for user login and registration
type UserService interface {
	// CreateUser function handles payload for user registration
	// Returns Status BadRequest if cannot decode the payload
	// or the user with the same email already exists
	// Status InternalServerError if error occurs during user creation
	// Status Created if user is successfully inserted into the db
	CreateUser(w http.ResponseWriter, r *http.Request)

	// Login function handles payload for user login
	// Returns Status BadRequest if cannot decode the payload
	// Status NotFound if user does not exist
	// Status InternalServerError if error occurs during JWT generation
	// Status OK and a cookie containing the token
	Login(w http.ResponseWriter, r *http.Request)
}
