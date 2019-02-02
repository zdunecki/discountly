package persistence

import (
	"github.com/zdunecki/discountly/features/auth/models"
	"gopkg.in/mgo.v2"
)

type Repository interface {
	UserExists(user auth.User) (bool, error)
	Create(user auth.User) (interface{}, error)
}

type DB struct {
	*mgo.Session
}
