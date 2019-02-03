package search

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/copier"
	"github.com/zdunecki/discountly/db"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/notifications"
	"github.com/zdunecki/discountly/features/search/finder"
	"github.com/zdunecki/discountly/features/search/models"
	"github.com/zdunecki/discountly/infra"
	"sync"
	"time"
)

var secret []byte

func init() {
	secret = []byte(infra.GetEnv("JWT_SECRET"))
}

//TODO: find faster and solution that jwt token
func createPromoCodeToken(discount discounts.ProtectedDiscount) (string, error) {
	claims := struct {
		discounts.ProtectedDiscount
		jwt.StandardClaims
	}{
		discount,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func findBestDiscounts(wg *sync.WaitGroup, criteria search.Search) ([]search.SearchDiscount, error) {
	discountChannel := make(chan search.SearchDiscount)
	errorChannel := make(chan error, 1)

	repo, err := database.NewRepo(database.DbConnection())
	if err != nil {
		return nil, err
	}

	defer repo.Discounts.Close()

	allByKeywords, err := repo.Discounts.FindAllByKeywords(criteria)
	if err != nil {
		return nil, err
	}
	bestDiscounts := finder.FindBestDiscounts(allByKeywords, criteria)

	if len(bestDiscounts) == 0 {
		return nil, nil
	}

	for ii, b := range bestDiscounts {
		wg.Add(1)
		go func(bestDiscount discounts.Discount, i int) {
			defer wg.Done()

			protectedDiscount := &discounts.ProtectedDiscount{}
			err := copier.Copy(&protectedDiscount, &bestDiscount)

			if err != nil {
				errorChannel <- err
				close(errorChannel)
				return
			}

			promoCodeToken, err := createPromoCodeToken(*protectedDiscount)

			if err != nil {
				errorChannel <- err
				close(errorChannel)
				return
			}

			discountChannel <- search.SearchDiscount{
				ProtectedDiscount: *protectedDiscount,
				Token: promoCodeToken,
			}

			if i == (len(bestDiscounts) - 1) {
				close(discountChannel)
			}
		}(b, ii)
	}

	select {
	case <-errorChannel:
		return nil, err
	default:
		var response []search.SearchDiscount

		for discount := range discountChannel {
			response = append(response, discount)
		}
		return response, nil
	}
}

func receiveHook(hook search.Hook) error {
	split := hook.Split()
	key := finder.NewGeoFenceUserId(hook.Id)

	conn := database.RedisPool.Get()

	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	values, err := redis.Values(
		conn.Do("SMEMBERS", key),
	)

	if err != nil {
		return err
	}

	for _, val := range values {
		userId := fmt.Sprintf("%s",  val)

		data := map[string]string{"detect": hook.Detect, "discount_id": split.DiscountId}

		if _, err := notifications.Client.Trigger("private-"+userId, finder.HookPointKey, data); err != nil {
			return err
		}
	}

	return nil
}
