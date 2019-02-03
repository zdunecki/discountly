package persistence

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) UserDiscounts(userId string) ([]discounts.Discount, error) {
	result := discounts.DiscountDefinition{}

	collection := db.getDiscountDefinitionCollection()

	if err := collection.Find(bson.M{"user_id": userId}).One(&result); err != nil {
		return nil, err
	}

	if result.Discounts == nil {
		return []discounts.Discount{}, nil
	}
	return result.Discounts, nil
}
