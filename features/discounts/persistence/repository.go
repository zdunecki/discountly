package persistence

import (
	authModels "github.com/zdunecki/discountly/features/auth/models"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Repository interface {
	CreateDefinition(user authModels.User) (bson.ObjectId, error)
	UserDiscounts(userId string) ([]discounts.Discount, error)
	CreateDiscounts(userId string, discounts []discounts.Discount) ([]discounts.Discount, error)
	UpdateUserDiscount(userId, discountId string, discount discounts.Discount) (interface{}, []discounts.Location, error)
	UpdateDiscount(discountId string, discount discounts.Discount) error
	DeleteDiscount(userId, discountId string) error
	Find(discountId string) (discounts.Discount, error)
	FindAll() ([]discounts.Discount, error)
	FindAllByKeywords(search search.Search) ([]discounts.Discount, error)
	Close()
}

type DB struct {
	*mgo.Session
}

func (db *DB) getDiscountDefinitionCollection() *mgo.Collection {
	return db.Session.DB("discountly").C("discount_definitions")
}
