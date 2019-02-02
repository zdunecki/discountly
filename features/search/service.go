package search

import (
	"github.com/zdunecki/discountly/db"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/finder"
	"github.com/zdunecki/discountly/features/search/models"
	"github.com/zdunecki/discountly/infra"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
	"sync"
	"time"
)

var secret []byte

func init() {
	secret = []byte(infra.GetEnv("JWT_SECRET"))
}


//TODO: find faster and solution that jwt token
func createPromoCode(discount discounts.Discount) (*discounts.PromoCodeToken, error) {
	claims := struct {
		discounts.Discount
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
		return nil, err
	}

	return &discounts.PromoCodeToken{Discount: tokenString}, nil
}

func findBestDiscounts(wg *sync.WaitGroup, s search.Search) ([]discounts.Discount, error) {
	discountChannel := make(chan *discounts.Discount)
	errorChannel := make(chan error, 1)

	repo, err := database.NewRepo(database.DbConnection())
	if err != nil {
		return nil, err
	}

	defer repo.Session.Close()

	result := repo.Search.FindAllDefinitionsByKeywords(s)

	cc := finder.FindBestDiscounts(result, s)

	if len(cc) == 0 {
		return nil, nil
	}

	for ii, best := range cc {
		wg.Add(1)
		go func(b discounts.Discount, i int) {
			defer wg.Done()

			promoCode, err := createPromoCode(b)

			if err != nil {
				errorChannel <- err
				close(errorChannel)
				return
			}

			newDiscount := &discounts.Discount{}
			err = copier.Copy(&newDiscount, &b)

			if err != nil {
				errorChannel <- err
				close(errorChannel)
				return
			}

			newDiscount.Token = promoCode.Discount
			discountChannel <- newDiscount

			if i == (len(cc) - 1) {
				close(discountChannel)
			}
		}(best, ii)
	}

	select {
	case <-errorChannel:
		return nil, err
	default:
		var response []discounts.Discount

		for d := range discountChannel {
			response = append(response, *d)
		}
		return response, nil
	}
}
