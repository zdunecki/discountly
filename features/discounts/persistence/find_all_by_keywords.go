package persistence

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) FindAllByKeywords(search search.Search) ([]discounts.Discount, error) {
	type AllDiscounts struct {
		Discounts []discounts.Discount `bson:"discounts" json:"discounts"`
	}

	var all []AllDiscounts

	collection := db.getDiscountDefinitionCollection()

	if err := collection.Find(bson.M{
		"discounts.keywords": bson.M{"$in": search.Keywords},
	}).Select(bson.M{"discounts": 1}).All(&all); err != nil {
		return nil, err
	}

	var result []discounts.Discount

	for _, d := range all {
		result = append(result, d.Discounts...)
	}

	return result, nil
}