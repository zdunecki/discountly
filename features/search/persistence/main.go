package persistence

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) FindAllDefinitionsByKeywords(search search.Search) []discounts.DiscountDefinition {
	discountSession := db.Session.DB("discountly").C("discount_definitions")

	var definitions []discounts.DiscountDefinition

	_ = discountSession.Find(
		bson.M{
			"discounts.keywords": bson.M{"$in": search.Keywords},
		},
	).All(&definitions)

	return definitions
}
