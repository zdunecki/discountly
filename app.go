package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zdunecki/discountly/features/auth"
	"github.com/zdunecki/discountly/features/discounts"
	"github.com/zdunecki/discountly/features/notifications"
	"github.com/zdunecki/discountly/features/search"
	"github.com/zdunecki/discountly/infra"
	"github.com/zdunecki/discountly/oauth"
	"net/http"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*") //only for tester

	searchRoute := r.Group("/search")
	searchRoute.POST("/", search.FindBestDiscounts)

	geoFencingRoute := r.Group("/geo-fencing")
	geoFencingRoute.POST("/nearby-hook", search.ReceiveNearbyHook)

	notificationsRoute := r.Group("/notifications")
	notificationsRoute.POST("/auth", notifications.Auth)
	notificationsRoute.POST("/receive-hook", notifications.ReceiveHook)

	oauthRoute := r.Group("/oauth")
	oauthRoute.GET("/login", oauth.Redirect)
	oauthRoute.GET("/callback", oauth.Callback)

	meRoute := r.Group("/me", auth.AuthorizedResources())
	meRoute.GET("/", auth.Me)

	authorizedDiscountsRoute := r.Group("/me/discounts", auth.AuthorizedOwnResources())
	authorizedDiscountsRoute.GET("/", discounts.GetUserDiscounts)
	authorizedDiscountsRoute.POST("/", discounts.CreateDiscounts)
	authorizedDiscountsRoute.PUT("/:id", discounts.UpdateDiscount)
	authorizedDiscountsRoute.DELETE("/:id", discounts.DeleteDiscount)

	discountsRoute := r.Group("/discounts")
	discountsRoute.GET("/", discounts.GetAllDiscounts)
	discountsRoute.POST("/promo-code", discounts.CreateDiscountPromoCode)

	r.GET("/demo", func(c *gin.Context) { //only for tester
		pusherKey := infra.GetEnv("PUSHER_APP_KEY")
		appURL := infra.GetEnv("EXPOSE_APP_URL")

		c.HTML(http.StatusOK, "demo.html", gin.H{"pusher_app_key": pusherKey, "app_url": appURL })
	})

	if err := r.Run(infra.GetExposedHost()); err != nil {
		panic(err)
	}
}
