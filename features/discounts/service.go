package discounts

import (
	"github.com/zdunecki/discountly/db"
	"github.com/zdunecki/discountly/features/auth"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/finder"
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"os"
)

func getDiscounts(userId string) ([]discounts.Discount, error) {
	repo, err := database.NewRepo(database.DbConnection())
	defer repo.Session.Close()

	if err != nil {
		return nil, err
	}

	userDiscounts, err := repo.Discounts.UserDiscounts(userId)

	if err != nil {
		return nil, err
	}

	return userDiscounts, nil
}

func createDiscounts(d []discounts.Discount, userId string) ([]discounts.Discount, error) {
	repo, err := database.NewRepo(database.DbConnection())
	defer repo.Session.Close()

	if err != nil {
		return nil, err
	}

	newDiscounts, err := repo.Discounts.CreateDiscounts(userId, d)

	if err != nil {
		return nil, err
	}

	for _, discount := range newDiscounts {
		go func(locations []discounts.Location) {
			finder.SetPoint(userId, locations)
		}(discount.Locations)
	}

	return newDiscounts, nil
}

func updateDiscount(discount discounts.Discount, discountId, userId string) error {
	repo, err := database.NewRepo(database.DbConnection())
	defer repo.Session.Close()

	if err != nil {
		return err
	}

	_, updatedLocations, err := repo.Discounts.UpdateUserDiscount(userId, discountId, discount)

	if err != nil {
		return err
	}

	finder.SetPoint(userId, updatedLocations)

	return nil
}

func deleteDiscount(discountId, userId string) error {
	repo, err := database.NewRepo(database.DbConnection())
	defer repo.Session.Close()

	if err != nil {
		return err
	}

	if err := repo.Discounts.DeleteDiscount(userId, discountId); err != nil {
		return err
	}

	return nil
}

func createDiscountPromoCode(discountToken string) error {
	repo, err := database.NewRepo(database.DbConnection())
	if err != nil {
		return err
	}
	defer repo.Session.Close()

	jwtToken, err := auth.ParseJWT(discountToken)
	if err != nil {
		return err.(*jwt.ValidationError).Inner
	}

	discountId := jwtToken.Claims.(jwt.MapClaims)["id"].(string)

	updateDiscount, err := repo.Discounts.Find(discountId)
	if err != nil {
		return err
	}

	promoCode := &discounts.PromoCode{}

	updateDiscount.PromoCodes = append(updateDiscount.PromoCodes, promoCode.New())

	if err := repo.Discounts.UpdateDiscount(discountId, updateDiscount); err != nil {
		return err
	}

	return nil
}

func revokeToken(discountToken string) error {
	conn, err := redis.Dial("tcp", os.Getenv("REDIS_CONNECTION"))
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Do("SET", "discountToken:"+discountToken, ""); err != nil {
		return err
	}

	if _, err := conn.Do("EXPIRE", "discountToken:"+discountToken, 60 * 5); err != nil {
		return err
	}

	return nil
}

func findRevokeToken(discountToken string) (interface{}, error) {
	conn, err := redis.Dial("tcp", os.Getenv("REDIS_CONNECTION"))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	find, err := conn.Do("GET", "discountToken:"+discountToken)
	if err != nil {
		return nil, err
	}
	return find, nil
}

