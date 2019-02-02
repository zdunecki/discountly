package persistence

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/lib"
	"gopkg.in/mgo.v2/bson"
)

func updateRules(rules []discounts.Rule) []discounts.Rule {
	var makeEditRules = make([]discounts.Rule, len(rules))

	for i, rule := range rules {
		editRule, err := rule.Edit()
		if err != nil {
			panic(err)
		}
		makeEditRules[i] = editRule
	}
	return makeEditRules
}

func updateLocations(locations []discounts.Location) []discounts.Location {
	var makeEditLocations = make([]discounts.Location, len(locations))

	for i, location := range locations {
		if location.Id == "" {
			makeEditLocations[i] = location.New()
		} else {
			makeEditLocations[i] = location
		}
	}
	return makeEditLocations
}

func (db *DB) UpdateUserDiscount(userId string, discountId string, discount discounts.Discount) (interface{}, []discounts.Location, error) {
	collection := db.getCollection()

	updateDiscount := discounts.Discount{
		Name:       discount.Name,
		ImageUrl:   discount.ImageUrl,
		Keywords:   discount.Keywords,
		PromoCodes: discount.PromoCodes,
		Rules:      updateRules(discount.Rules),
		Locations:  updateLocations(discount.Locations),
	}

	if response, err := collection.Upsert(
		bson.M{
			"user_id":      userId,
			"discounts.id": discountId,
		},
		bson.M{
			"$set": lib.PrettyBsonSet("discounts.$", &updateDiscount),
		}); err != nil {
		return nil, nil, err
	} else {
		return response.UpsertedId, updateDiscount.Locations, nil
	}
}

func (db *DB) UpdateDiscount(discountId string, discount discounts.Discount) error {
	collection := db.getCollection()

	updateDiscount := discounts.Discount{
		Name:       discount.Name,
		ImageUrl:   discount.ImageUrl,
		Keywords:   discount.Keywords,
		PromoCodes: discount.PromoCodes,
		Rules:      updateRules(discount.Rules),
		Locations:  updateLocations(discount.Locations),
	}

	if _, err := collection.Upsert(
		bson.M{
			"discounts.id": discountId,
		},
		bson.M{
			"$set": lib.PrettyBsonSet("discounts.$", &updateDiscount),
		}); err != nil {
		return err
	}

	return nil
}
