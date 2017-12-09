package model

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kshvakov/clickhouse"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

	const dbcountQuery = `select count(distinct point_id) from stat.statistics where created = toDate(today())`


type  PointCount struct {
	Count int
}
func NewRedis() *redis.Client {
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


func СountQuery() (int, error) {
	var pointCount PointCount
	res, err := http.Get("http://127.0.0.1:8080/gateway/statistics/all")
	if err != nil {
		fmt.Println("Error get to api",err)
	}
	count, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)
	err = json.Unmarshal([]byte(count), &pointCount)
	if err != nil {
		return 0, fmt.Errorf("error unmarshal", err)
	}
	return pointCount.Count, nil
}