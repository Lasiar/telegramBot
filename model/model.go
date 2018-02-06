package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"telega/lib"
)

func List() ([]int, error) {
	var point lib.Point
	res, err := http.Get("http://127.0.0.1:8181/gateway/telegram/list-point")
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
		return nil, fmt.Errorf("error unmarshal %v", err)
	}
	sort.Ints(point.Point)
	return point.Point, nil
}

func InitialGoodStatistic(js lib.RequestGoodStatistic) {
	url := fmt.Sprint("http://127.0.0.1:8181/gateway/telegram/initial-good-point")
	jsonStr, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func InfoPoint(point string) (lib.InfoPointJs, error) {
	var infoPointJs lib.InfoPointJs
	url := fmt.Sprint("http://127.0.0.1:8181/gateway/telegram/info-point?point=", point)
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

func CountAllQuery() (int, error) {
	var pointCount lib.PointCount
	res, err := http.Get("http://127.0.0.1:8181/gateway/telegram/count-point")
	if err != nil {
		fmt.Println("Error get to api", err)
	}
	count, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(count, &pointCount)
	if err != nil {
		return 0, fmt.Errorf("error unmarsha %v", err)
	}
	return pointCount.Count, nil
}

func CountToDayQuery() ([]int, error) {
	var pointToDayCount lib.Point
	res, err := http.Get("http://127.0.0.1:8181/gateway/telegram/list-point-today")
	if err != nil {
		fmt.Println("Error get to api", err)
	}
	count, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(count, &pointToDayCount)
	if err != nil {
		return pointToDayCount.Point, fmt.Errorf("error unmarsha %v", err)
	}
	fmt.Println(pointToDayCount)
	return pointToDayCount.Point, nil
}
