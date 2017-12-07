package model

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kshvakov/clickhouse"
	"log"
)

const dbcountQuery = `select count(distinct point_id) from stat.statistics where created = toDate(today())`


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

func NewClick() *sql.DB {
	db, err := sql.Open("clickhouse", "tcp://192.168.0.145:9000?database=stat&read_timeout=10&write_timeout=20")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
	}
	return db
}

func СountQuery(db *sql.DB) (int, error) {
	var count int
	if err := db.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			return 0, fmt.Errorf("Error connect to GoodClick: %v", err)
		}
	}
	rows, err := db.Query(dbcountQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return count, nil
}