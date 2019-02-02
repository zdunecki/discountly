package persistence

import "gopkg.in/mgo.v2/bson"

func (db *DB) DeleteDiscount(userId, discountId string) error {
	collection := db.getCollection()

	err := collection.Update(
		bson.M{
			"user_id": userId,
		},
		bson.M{
			"$pull": bson.M{
				"discounts": bson.M{"id": discountId},
			},
		},
	)

	return err
}
