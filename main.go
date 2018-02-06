package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"telega/lib"
	"telega/telegram"
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
	msgFromMachine := make(chan string)
	msg := make(chan tgbotapi.Update)

	go telegram.ReceivingMessageTelegram(msg)
	go telegram.Worker(msg, msgFromMachine)

	handleHello := recivingBadStatistic(msgFromMachine)
	http.HandleFunc("/bad", handleHello)
	http.ListenAndServe(":8282", nil)
}

func recivingBadStatistic(messageFromMachine chan string) func(http.ResponseWriter, *http.Request) {
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
		msg := fmt.Sprintf("ip: *%v* json *%v*", t.Ip, t.Json)
		select {
		case messageFromMachine <- msg:
			return
		default:
			return
		}
	}
}
