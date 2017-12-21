package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"net/http"
	"encoding/json"
	"telega/lib"
	"telega/telegram"

	"io/ioutil"
	//"regexp"
)



func init() {
	flag.StringVar(&lib.TelegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.Parse()
	if lib.TelegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}
}

func main() {
	logs := make(chan string)
	go telegram.MainTtelegram(logs)
	handleHello := makeHello(logs)
	http.HandleFunc("/listen", handleHello)
	http.ListenAndServe(":8282", nil)

}


func makeHello(logger chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "all ok	")
		decoder := json.NewDecoder(r.Body)
		var t lib.Json

		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println(ioutil.ReadAll(r.Body))
			fmt.Fprint(w, "all ok	")
			log.Println(err)
			return
		}
		fmt.Println(t.Point)
		string := fmt.Sprint("id: ", t.Point, "info: ", t.Statistics)
		fmt.Println("Отправил в логгер")
		logger <- string
	}
}

