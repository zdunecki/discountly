package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zdunecki/discountly/db"
	"net/http"
)

func Me(c *gin.Context) {
	user, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	repo, err := database.NewRepo(database.DbConnection())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		c.Abort()
		return
	}

	defer repo.Auth.Close()

	userId := fmt.Sprintf("%s", user)

	me, err := repo.Auth.Me(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, me)
}