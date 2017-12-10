package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"net/http"
)

const dbcountQuery = `select count(distinct point_id) from stat.statistics where created = toDate(today())`

type InfoPointJs struct {
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Success   bool   `json:"success"`
}

type PointCount struct {
	Count int `json:"count"`
}

type Point struct {
	Point []int `json:"point"`
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

func List() ([]int, error) {
	var point Point
	res, err := http.Get("http://127.0.0.1:8080/gateway/statistics/list")
	if err != nil {
		fmt.Println("Error get to api", err)
	}
	pointArr, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pointArr)
	err = json.Unmarshal([]byte(pointArr), &point)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal", err)
	}
	return point.Point, nil
}



func InfoPoint(point string) (InfoPointJs, error) {
	var infoPointJs InfoPointJs
	url := fmt.Sprint("http://127.0.0.1:8080/gateway/statistics/info-point?point=", point)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error get to api", err)
	}
	infoPoint, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal([]byte(infoPoint), &infoPointJs)
	if err != nil {
		return infoPointJs, fmt.Errorf("error unmarshal %v", err)
	}
	return infoPointJs, nil
}

func СountQuery() (int, error) {
	var pointCount PointCount
	res, err := http.Get("http://127.0.0.1:8080/gateway/statistics/all")
	if err != nil {
		fmt.Println("Error get to api", err)
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
