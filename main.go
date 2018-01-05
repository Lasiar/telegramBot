package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"telega/lib"
	"telega/telegram"

	"io/ioutil"
	//"regexp"
	"gopkg.in/telegram-bot-api.v4"
)

func init() {
	var err error
	flag.StringVar(&lib.TelegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.Parse()
	if lib.TelegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}
	telegram.Bot, err = tgbotapi.NewBotAPI(lib.TelegramBotToken)
	if err != nil {
		log.Panic("ошибка подключения бота", err)
	}
	log.Println("Connect", telegram.Bot.Self.UserName)
}

func main() {
	logs := make(chan string)
	msg := make(chan tgbotapi.Update)
	go telegram.ReceivingMessageTelegram(msg)
	go telegram.Worker(msg, logs)
	handleHello := makeHello(logs)
	http.HandleFunc("/gateway/telegram/create/bad", handleHello)
	http.ListenAndServe(":8181", nil)

}

func makeHello(logger chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "all ok	")
		decoder := json.NewDecoder(r.Body)
		var t lib.BadJson

		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println(ioutil.ReadAll(r.Body))
			fmt.Fprint(w, "all ok	")
			log.Println(err)
			return
		}
		string := fmt.Sprintf("ip: *%v* json *%v*",t.Ip , t.Json)
		fmt.Println(string)
		select {
		case <-logger:
			fmt.Println("отправил")
			logger <- string
		default:
			fmt.Println("Не отправил")
			return
		}
	}
}
