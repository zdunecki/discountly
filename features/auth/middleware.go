package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func notAuthorized(c *gin.Context) {
	c.Abort()
	c.JSON(http.StatusUnauthorized, http.StatusUnauthorized)
}

func getJwtToken(c *gin.Context) (*jwt.Token, error) {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		notAuthorized(c)
		return nil, nil
	}
	token := strings.Split(bearerToken, " ")[1]

	return ParseJWT(token)
}

func authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := getJwtToken(c)

		if err != nil {
			notAuthorized(c)
			return
		}
		c.Next()
	}
}

func AuthorizedOwnResources() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := getJwtToken(c)

		if err != nil {
			notAuthorized(c)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)

		if claims["user_id"] == c.GetHeader("x-user-id") {
			c.Next()
			return
		}
		notAuthorized(c)
		return
	}
}
