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
