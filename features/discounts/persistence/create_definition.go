package persistence

import (
	"github.com/zdunecki/discountly/features/auth/models"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) CreateDefinition(user auth.User) (bson.ObjectId, error) {
	collection := db.getDiscountDefinitionCollection()

	id := uuid.NewV4()

	discountDefinition := struct {
		Id                           string
		discounts.DiscountDefinition `bson:",inline"`
	}{
		id.String(),
		discounts.DiscountDefinition{
			UserId:    user.Id,
			Company:   user.Name,
			Discounts: []discounts.Discount{},
		},
	}

	if response, err := collection.Upsert(discountDefinition, &discountDefinition); err != nil {
		return "", err
	} else {
		return response.UpsertedId.(bson.ObjectId), nil
	}
}
