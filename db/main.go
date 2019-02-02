package database

import (
	authPersistence "github.com/zdunecki/discountly/features/auth/persistence"
	discountsPersistence "github.com/zdunecki/discountly/features/discounts/persistence"
	searchPersistence "github.com/zdunecki/discountly/features/search/persistence"
	"github.com/zdunecki/discountly/infra"
	"gopkg.in/mgo.v2"
)

type Session interface {
	Close()
}

type Repository struct {
	Discounts discountsPersistence.Repository
	Auth      authPersistence.Repository
	Search    searchPersistence.Repository
	Session   Session
}

func DbConnection() string {
	return infra.GetEnv("DB_CONNECTION")
}

func NewRepo(dbURI string) (*Repository, error) {
	session, err := mgo.Dial(dbURI)

	if err != nil {
		panic(err)
	}
	discounts := &discountsPersistence.DB{
		Session: session,
	}
	auth := &authPersistence.DB{
		Session: session,
	}
	search := &searchPersistence.DB{
		Session: session,
	}
	return &Repository{
		discounts,
		auth,
		search,
		session,
	}, err
}
