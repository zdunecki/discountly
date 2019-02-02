package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/zdunecki/discountly/infra"
	"time"
)

//TODO: generate private/public .pem instead of .env secret
var (
	invalidToken = errors.New("Invalid token")
)

var tokenExpiration = time.Hour * 4

var secret []byte

func init() {
	secret = []byte(infra.GetEnv("JWT_SECRET"))
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if token == nil {
		return nil, invalidToken
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, nil
	} else {
		return nil, err
	}
}

func CreateJwtToken(userId string) (string, error) {
	claims := struct {
		UserId string `json:"user_id"`
		jwt.StandardClaims
	}{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpiration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", nil
	}

	return tokenString, nil
}
