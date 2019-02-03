package oauth

import (
	"encoding/json"
	"github.com/zdunecki/discountly/db"
	"github.com/zdunecki/discountly/features/auth"
	models "github.com/zdunecki/discountly/features/auth/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Redirect(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, getAuthUrl())
}

func signUpUserOrDefault(user models.User) error {
	repo, err := database.NewRepo(database.DbConnection())
	if err != nil {
		return err
	}
	defer repo.Auth.Close()
	defer repo.Discounts.Close()

	exists, err := repo.Auth.UserExists(user)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	if _, err := repo.Auth.Create(user); err != nil {
		return err
	}

	_, err = repo.Discounts.CreateDefinition(user)

	return err
}

func Callback(c *gin.Context) {
	var content []byte
	var err error = nil

	code := c.Query("code")
	token := c.Query("token")

	if token != "" {
		content, _, err = userInfoFromToken(token)
	} else {
		content, _, err = getUserInfo(c.Query("state"), code, false)
	}

	if err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	jwtToken, err := callback(content)
	if err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, jwtToken)
}

func callback(content []byte) (string, error) {
	var user models.User

	if err := json.Unmarshal(content, &user); err != nil {
		return "", err
	}
	if err := signUpUserOrDefault(user); err != nil {
		return "", err
	}
	if jwtToken, err := auth.CreateJwtToken(user.Id); err != nil {
		return "", err
	} else {
		return jwtToken, nil
	}
}
