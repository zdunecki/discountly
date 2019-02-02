package persistence

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/models"
	"gopkg.in/mgo.v2"
)

type Repository interface {
	FindAllDefinitionsByKeywords(search search.Search) []discounts.DiscountDefinition
}

type DB struct {
	*mgo.Session
}
