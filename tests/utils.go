package tests

import (
	"github.com/zdunecki/discountly/infra"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/mgo.v2"
)

var dbAddress string
var geoAddress string

func init() {
	dbAddress = infra.GetEnv("DB_CONNECTION")
	geoAddress = infra.GetEnv("GEO_CONNECTION")
}

func DeleteAll() {
	session, err := mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	conn, err := redis.Dial("tcp", geoAddress)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	_ = session.DB("discountly").DropDatabase()
	_ = conn.Flush()
}
