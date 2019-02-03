package database

import (
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

var RedisPool *redis.Pool
var GeoRedisPool *redis.Pool

func init() {
	RedisPool = &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", os.Getenv("REDIS_CONNECTION"))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	GeoRedisPool = &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", os.Getenv("GEO_CONNECTION"))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}