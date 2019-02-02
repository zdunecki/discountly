package discounts

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"net/http"
)

func GetDiscounts(c *gin.Context) {
	userId := c.GetHeader("x-user-id")

	response, err := getDiscounts(userId)

	if err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, response)
}

func CreateDiscounts(c *gin.Context) {
	var bodyDiscounts []discounts.Discount
	if err := c.Bind(&bodyDiscounts); err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	userId := c.GetHeader("x-user-id")

	response, err := createDiscounts(bodyDiscounts, userId)

	if err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, response)
}

func UpdateDiscount(c *gin.Context) {
	var discount discounts.Discount
	if err := c.Bind(&discount); err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	userId := c.GetHeader("x-user-id")

	err := updateDiscount(discount, c.Param("id"), userId)

	if err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, http.StatusOK)
}

func DeleteDiscount(c *gin.Context) {
	userId := c.GetHeader("x-user-id")

	err := deleteDiscount(c.Param("id"), userId)

	if err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, http.StatusOK)
}

func CreateDiscountPromoCode(c *gin.Context) {
	token, err := findRevokeToken(c.Query("token"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"errors": err.Error()})
		return
	}

	if token != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"errors": "Token is already used"})
		return
	}

	if err := createDiscountPromoCode(c.Query("token")); err != nil {
		c.JSON(http.StatusInternalServerError,  map[string]interface{}{"errors": err.Error()})
		return
	}

	if err := revokeToken(c.Query("token")); err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, http.StatusOK)
}
