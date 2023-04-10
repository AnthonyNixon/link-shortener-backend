package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"strings"
	"time"
)

var (
	JWT_SIGNING_KEY []byte
	AUDIENCE        = os.Getenv("DATASTORE_NAMESPACE")
)

var (
	errorNoAuthHeader     = errors.New("no authorization header content present")
	errorAuthHeaderFormat = errors.New("authorization header format incorrect, should be 'bearer <token>`")
)

const (
	TOKEN_VALID_TIME = 24 * time.Hour
	ISSUER           = "link-shortener-backend-api"
)

func init() {
	log.Print("Initializing Authentication")
	signingKey := os.Getenv("JWT_SIGNING_KEY")
	if signingKey == "" {
		log.Fatal("No Signing Key Present.")
	}

	JWT_SIGNING_KEY = []byte(signingKey)
	log.Print("done")
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func New(username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(TOKEN_VALID_TIME)
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ISSUER,
			Audience:  []string{AUDIENCE},
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err = token.SignedString(JWT_SIGNING_KEY)

	return

}

func parseHeader(header string) (token string, err error) {
	if header == "" {
		return "", errorNoAuthHeader
	}

	// Tokens will be of format "bearer <token>", split on ' ' space
	content := strings.Split(header, " ")
	if len(content) != 2 {
		return "", errorAuthHeaderFormat
	}

	token = content[1]
	return
}

func ParseToken(header string) (username string, err error) {
	tokenString, err := parseHeader(header)
	if err != nil {
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return JWT_SIGNING_KEY, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Printf("%v %v %v\n", claims.Username, claims.Issuer, claims.Audience)
		username = claims.Username
	} else {
		fmt.Println(err)
	}

	return
}
