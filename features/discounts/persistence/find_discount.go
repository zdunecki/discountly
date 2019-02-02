package persistence

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) Find(discountId string) (discounts.Discount, error) {
	result := discounts.DiscountDefinition{}

	collection := db.getCollection()

	if err := collection.Find(bson.M{"discounts.id": discountId}).One(&result); err != nil {
		return discounts.Discount{}, err
	}

	if result.Discounts == nil {
		return discounts.Discount{}, nil
	}

	for _, d := range result.Discounts {
		if d.Id == discountId {
			return d, nil
		}
	}

	return discounts.Discount{}, nil
}
