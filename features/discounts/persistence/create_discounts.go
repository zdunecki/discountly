package persistence

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

func newRules(rules []discounts.Rule) []discounts.Rule {
	var makeRules = make([]discounts.Rule, len(rules))

	for i, rule := range rules {
		newRule, err := rule.New()
		if err != nil {
			panic(err)
		}
		makeRules[i] = newRule
	}
	return makeRules
}

func newLocations(locations []discounts.Location) []discounts.Location {
	var newLocations = make([]discounts.Location, len(locations))

	for i, r := range locations {
		newLocations[i] = r.New()
	}
	return newLocations
}

func (db *DB) CreateDiscounts(userId string, d []discounts.Discount) ([]discounts.Discount, error) {
	collection := db.getCollection()

	var newDiscounts = make([]discounts.Discount, len(d))

	for i, d := range d {
		id := uuid.NewV4()
		newDiscounts[i] = discounts.Discount{
			Id:         id.String(),
			Name:       d.Name,
			ImageUrl:   d.ImageUrl,
			Keywords:   d.Keywords,
			PromoCodes: d.PromoCodes,
			Locations:  newLocations(d.Locations),
			Rules:      newRules(d.Rules),
		}
	}

	if _, err := collection.Upsert(
		bson.M{
			"user_id": userId,
		},
		bson.M{
			"$push": bson.M{"discounts": bson.M{"$each": newDiscounts}},
		}); err != nil {
		return nil, err
	} else {
		return newDiscounts, err
	}
}
