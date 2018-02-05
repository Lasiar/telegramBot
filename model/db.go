package model

import (
	"github.com/go-redis/redis"
	"log"
)

var dbRedis *redis.Client

func NewRedis(addr string, password string) *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := db.Ping().Result()
	if err != nil {
		log.Println(err)
	}
	return db
}
