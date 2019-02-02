package persistence

import (
	"github.com/zdunecki/discountly/features/auth/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) UserExists(user auth.User) (bool, error) {
	userSession := db.Session.DB("discountly").C("users")

	if count, err := userSession.Find(
		bson.M{
			"id": user.Id,
		}).Count(); err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}

func (db *DB) Create(user auth.User) (interface{}, error) {
	userSession := db.Session.DB("discountly").C("users")

	if response, err := userSession.Upsert(user, &user); err != nil {
		return nil, err
	} else {
		return response.UpsertedId, nil
	}
}
