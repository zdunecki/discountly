package database

import (
	authPersistence "github.com/zdunecki/discountly/features/auth/persistence"
	discountsPersistence "github.com/zdunecki/discountly/features/discounts/persistence"
	searchPersistence "github.com/zdunecki/discountly/features/search/persistence"
	"github.com/zdunecki/discountly/infra"
	"gopkg.in/mgo.v2"
)

type Repository struct {
	Discounts discountsPersistence.Repository
	Auth      authPersistence.Repository
	Search    searchPersistence.Repository
}

func DbConnection() string {
	return infra.GetEnv("DB_CONNECTION")
}
var mgoSession *mgo.Session

func NewRepo(dbURI string) (*Repository, error) {
	if mgoSession == nil {
		session, err := mgo.Dial(dbURI)
		if err != nil {
			return nil, err
		}

		mgoSession = session
	}

	discounts := &discountsPersistence.DB{
		Session: mgoSession.Copy(),
	}
	auth := &authPersistence.DB{
		Session: mgoSession.Copy(),
	}
	search := &searchPersistence.DB{
		Session: mgoSession.Copy(),
	}
	return &Repository{
		discounts,
		auth,
		search,
	}, nil
}
