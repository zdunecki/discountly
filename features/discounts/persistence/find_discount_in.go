package persistence

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) FindInIds(discountIds []string) ([]discounts.Discount, error) {
	var project []projectDiscounts

	collection := db.getDiscountDefinitionCollection()

	if err := collection.
		Find(bson.M{"discounts.id": bson.M{"$in": discountIds}}).
		Select(bson.M{"discounts.$": 1}).
		All(&project); err != nil {
		return nil, err
	}

	var result []discounts.Discount
	for _, discount := range project {
		result = append(result, discount.Discounts...)
	}

	return result, nil
}
