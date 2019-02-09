package search

import (
	"github.com/dgrijalva/jwt-go"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/zdunecki/discountly/features/auth"
	"github.com/zdunecki/discountly/features/discounts/models"
	"testing"
)

func TestCreatingDiscountPromoCode(t *testing.T) {
	assert := assert2.New(t)

	discount := discounts.ProtectedDiscount{
		Id: "test-id",
	}

	promoCodeToken, err := createPromoCodeToken(discount)

	if err != nil {
		assert.Error(err)
	}

	jwtToken, err := auth.ParseJWT(promoCodeToken)
	if err != nil {
		assert.Error(err)
	}

	discountId := jwtToken.Claims.(jwt.MapClaims)["id"].(string)

	assert.Equal(discountId, "test-id")
}
