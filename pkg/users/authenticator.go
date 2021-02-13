package users

import "net/http"

// UserAuthenticator provides methods for JWT management
type UserAuthenticator interface {
	// GenerateTokenForUser function generate JWT for given user
	// Return an error if such occurs during token signing or the token string
	GenerateTokenForUser(user *User) (string, error)

	// VerifyJWT function serves as a middleware handler for user token verification
	// Returns status Forbidden if token is invalid
	// or sets the user context otherwise
	VerifyJWT(next http.Handler) http.Handler
}
