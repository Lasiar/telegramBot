package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/go-redis/redis"
	"log"
	"os"
	"regexp"
	"telega/model"
	"strconv"
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
		case m == "help":
			reply = "/list - показ всех машин, о которых есть информация" +
				"/$(id) - вывод информации по конкретной машине"
		case m == "list":
			key, err := dbRedis.Keys("*_ip*").Result()
			if err != nil {
				reply = fmt.Sprint(err)
				break
			}
			if len(key) == 0 {
				reply = "Нет машин.."
				break
			}
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
			ip := m + "_ip"
			user := m + "_user"
			keysdb, err := dbRedis.MGet(ip, user).Result()
			if err != nil {
				reply = fmt.Sprint(err)
			}
			for i, mes := range keysdb {
				switch i {
				case 0:
					ip = fmt.Sprint(mes)
				case 1:
					user = fmt.Sprint(mes)
				}
			}
			if ip == "<nil>" {
				reply = "Нет такой машины"
			} else {
				reply = fmt.Sprintf("*ip:* _%v_ * user agent: *_%v_", ip, user)
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "markdown"
		bot.Send(msg)
	}
}
