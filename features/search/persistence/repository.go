package persistence

import (
	"gopkg.in/mgo.v2"
)

type Repository interface {
	Close()
}

type DB struct {
	*mgo.Session
}
