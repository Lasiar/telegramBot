package telegram

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"regexp"
	"strconv"
	"telega/lib"
	"telega/model"
)

func MainTtelegram(logs chan string) {
	var id = regexp.MustCompile(`\d`)
	bot, err := tgbotapi.NewBotAPI(lib.TelegramBotToken)
	if err != nil {
		log.Panic("ошибка подключения бота", err)
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
		case m == "1":

		case m == "list":
			reply = ""
			keys, _ := model.List()
			for _, key := range keys {
				reply = reply + fmt.Sprint(key, "; ")
			}
		case m == "count":
			count, err := model.CountAllQuery()
			if err != nil {
				log.Println(err)
				reply = "ошибка"
				break
			}
			reply = strconv.Itoa(count)
		case m == "point today":
			reply = ""
			keys, _ := model.CountToDayQuery()
			for _, key := range keys {
				reply = reply + fmt.Sprint(key, "\n")
			}
		case id.MatchString(m):
			id := update.Message.Chat.ID
			go func(id int) {
			for {
				select {
				case u := <-updates:
				
					switch u.Message.Text {
					case "exit":
						return
					default:
						reply = "print exit"
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
						msg.ParseMode = "markdown"
						bot.Send(msg)
					}
				case reply := <-logs:
					fmt.Println("Жду получения")
					reply = <-logs
					fmt.Println("Получил")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					msg.ParseMode = "markdown"
					bot.Send(msg)
				}
			}
		}(id)
			/*infoPoint, _ := model.InfoPoint(m)
			if infoPoint.Success {
				reply = fmt.Sprintf("ip: *%v*; user info: *%v*", infoPoint.Ip, infoPoint.UserAgent)
			} else {
				reply = "такой машины нет"
			} */
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "markdown"
		bot.Send(msg)
	}
}
