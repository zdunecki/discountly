package main

import (
	"github.com/zdunecki/discountly/features/auth"
	"github.com/zdunecki/discountly/features/discounts"
	"github.com/zdunecki/discountly/features/search"
	"github.com/zdunecki/discountly/infra"
	"github.com/zdunecki/discountly/oauth"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/search", search.FindBestDiscounts)

	oauthRoute := r.Group("/oauth")
	oauthRoute.GET("/login", oauth.Redirect)
	oauthRoute.GET("/callback", oauth.Callback)

	authorizedDiscounts := r.Group("/discounts", auth.AuthorizedOwnResources())
	authorizedDiscounts.GET("/", discounts.GetDiscounts)
	authorizedDiscounts.POST("/", discounts.CreateDiscounts)
	authorizedDiscounts.PUT("/:id", discounts.UpdateDiscount)
	authorizedDiscounts.DELETE("/:id", discounts.DeleteDiscount)

	r.POST("/discounts/promo-code", discounts.CreateDiscountPromoCode)

	if err := r.Run(infra.GetHost()); err != nil {
		panic(err)
	}
}
