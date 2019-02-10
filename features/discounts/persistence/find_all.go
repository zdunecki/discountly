package persistence

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"gopkg.in/mgo.v2/bson"
)

type projectDiscounts struct {
	Discounts []discounts.Discount `bson:"discounts" json:"discounts"`
}

// TODO: pagination
func (db *DB) FindAll() ([]discounts.Discount, error) {
	var all []projectDiscounts

	collection := db.getDiscountDefinitionCollection()

	if err := collection.Find(nil).Select(bson.M{"discounts": 1}).All(&all); err != nil {
		return nil, err
	}

	var result []discounts.Discount

	for _, discount := range all {
		result = append(result, discount.Discounts...)
	}

	return result, nil
}