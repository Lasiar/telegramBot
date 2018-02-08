package lib

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

type Json struct {
	Point      interface{}     `json:"point"`
	Statistics [][]interface{} `json:"statistics"`
}

type BadJson struct {
	Ip   string
	Json interface{}
}

type RequestGoodStatistic struct {
	ChatId int64 `json:"chat_id"`
	Point  []int `json:"point"`
}

type GoodJson struct {
	Point    int
	Datetime int64
	Md5      string
	Len      int
}

type MessageChat struct {
	ChatId  int64
	Message string
}
