package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/krasimiraMilkova/cookit/internal/appconfig"
	"github.com/krasimiraMilkova/cookit/pkg/users"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type JwtAuthenticator struct{}

const (
	privateKeyPath = "/keys/app.rsa"
	publicKeyPath  = "/keys/app.rsa.pub"
	TokenName      = "cookit-access-token"
	Expiration     = 30 * time.Minute
)

var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
)

var authenticator *JwtAuthenticator

func GetAuthenticator() *JwtAuthenticator {
	if authenticator == nil {
		authenticator = &JwtAuthenticator{}
		initKeys()
	}

	return authenticator
}

func initKeys() {
	projectDir := appconfig.Get().GetProjectDir()

	privateKeyFileContent, err := ioutil.ReadFile(projectDir + privateKeyPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyFileContent)
	if err != nil {
		log.Fatal(err)
		return
	}

	publicKeyFileContent, err := ioutil.ReadFile(projectDir + publicKeyPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyFileContent)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (jwtAuth JwtAuthenticator) GenerateTokenForUser(user *users.User) (string, error) {
	token := &Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(Expiration).Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, token)
	tokenString, err := jwtToken.SignedString(privateKey)
	if err != nil {
		return "", errors.New("error while signing generated jwt token")
	}

	return tokenString, nil
}

func (jwtAuth JwtAuthenticator) VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := getToken(r)
		if token == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		valid := isTokenValid(token)
		if !valid {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		user, err := userFromToken(token)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getToken(r *http.Request) string {
	var token = ""
	cookie, err := r.Cookie(TokenName)
	if err == nil {
		token = cookie.Value
	}
	return token
}

func isTokenValid(token string) bool {
	tk := &Token{}

	_, err := jwt.ParseWithClaims(token, tk, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return false
	}

	return true
}

func userFromToken(tokenString string) (*users.User, error) {
	token := &Token{}

	_, err := jwt.ParseWithClaims(tokenString, token, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	var usr = users.User{
		ID:    token.UserID,
		Email: token.Email,
		Name:  token.Name,
	}
	return &usr, err
}
