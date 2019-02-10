package persistence

import (
	"github.com/zdunecki/discountly/features/auth/models"
	"gopkg.in/mgo.v2"
)

type Repository interface {
	Me(userId string) (auth.User, error)
	UserExists(user auth.User) (bool, error)
	Create(user auth.User) (interface{}, error)
	Close()
}

type DB struct {
	*mgo.Session
}
