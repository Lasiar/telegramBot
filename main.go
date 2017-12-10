package main

import (
	"flag"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/go-redis/redis"
	"log"
	"os"
	"regexp"
	"strconv"
	"telega/model"
	//	"github.com/prometheus/common/version"
)

var (
	telegramBotToken string
	dbRedis          *redis.Client
)

func init() {
	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.Parse()

	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}
	dbRedis = model.NewRedis()
}

func main() {
	var id = regexp.MustCompile(`\d`)
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		m := update.Message.Text
		switch {
		case m == "list":
			key, _ := model.List()
			reply = fmt.Sprint(key)
		case m == "count":
			count, err := model.СountQuery()
			if err != nil {
				log.Println(err)
				reply = "ошибка"
				break
			}
			reply = strconv.Itoa(count)
		case id.MatchString(m):
			infoPoint, _ := model.InfoPoint(m)
			fmt.Println(infoPoint.Success)
			if infoPoint.Success {
				reply = fmt.Sprintf("ip: *%v*; user info: *%v*", infoPoint.Ip, infoPoint.UserAgent)
			} else {
				reply = "такой машины нет"
			}
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "markdown"
		bot.Send(msg)
	}
}
