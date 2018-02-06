package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"telega/lib"
)

func AdmissionStatistic(goodStat chan lib.GoodJson) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "all ok	")
		decoder := json.NewDecoder(r.Body)
		var t lib.GoodJson

		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println(ioutil.ReadAll(r.Body))
			fmt.Fprint(w, "all ok	")
			log.Println(err)
			return
		}
		select {
		case goodStat <- t:
			return
		default:
			return
		}
	}
}

func SendWarning(message chan lib.MessageChat) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "all ok	")
		var t lib.MessageChat
		msg	:= r.FormValue("message")
		t.Message=msg
		message <-t
	}
	}