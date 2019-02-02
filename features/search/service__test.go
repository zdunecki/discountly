package search

import (
	"github.com/zdunecki/discountly/features/auth"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/dgrijalva/jwt-go"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatingDiscountPromoCode(t *testing.T) {
	assert := assert2.New(t)

	d := discounts.Discount{
		Id: "test-id",
	}

	response, err := createPromoCode(d)

	if err != nil {
		assert.Error(err)
	}

	jwtToken, err := auth.ParseJWT(response.Discount)
	if err != nil {
		assert.Error(err)
	}

	discountId := jwtToken.Claims.(jwt.MapClaims)["id"].(string)

	assert.Equal(discountId, "test-id")
}
