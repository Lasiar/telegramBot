package telegram

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"regexp"
	"strconv"
	"telega/model"
	//"net/http"
)

var Bot *tgbotapi.BotAPI

var id = regexp.MustCompile(`\d`)

func ReceivingMessageTelegram(msg chan tgbotapi.Update) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println(err)
	}
	for update := range updates {
		msg <- update
	}
}

func Worker(update chan tgbotapi.Update, msgFromMachine chan string) {
	var idListenGoodId []int64
	idReturn := make(chan int64)
	broadcast := make(chan string)
	msgForListen := make(chan string)
loop:
	for {
		select {
		case u := <-update:
			for _, id := range idListenGoodId {
				if u.Message.Chat.ID == id {
					msgForListen <- u.Message.Text
					continue loop
				}
			}
			v1 := regular(u, broadcast, msgForListen, idReturn)
			if v1 != 0 {
				idListenGoodId = append(idListenGoodId, v1)
			}

		case b := <-msgFromMachine:
			for range idListenGoodId {
				broadcast <- b
			}

		case i := <-idReturn:
			fmt.Println("Должен удалить: ", i)
			idListenGoodId = deleteByValue(i, idListenGoodId)
			fmt.Println(idListenGoodId)
		}
		fmt.Println(idListenGoodId)
	}
}

func regular(update tgbotapi.Update, msgFromMachine chan string, msgForListen chan string, idReturn chan int64) int64 {
	reply := "Не знаю что ответить"
	m := update.Message.Text
	switch {
	case m == "bad":
		go sendStatistics(update.Message.Chat.ID, msgForListen, msgFromMachine, idReturn)
		return update.Message.Chat.ID
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
		infoPoint, _ := model.InfoPoint(m)
		if infoPoint.Success {
			reply = fmt.Sprintf("ip: *%v*; user info: *%v*", infoPoint.Ip, infoPoint.UserAgent)
		} else {
			reply = "такой машины нет"
		}
	}
	sendMessage(update.Message.Chat.ID, reply)
	return 0
}

func sendMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "markdown"
	Bot.Send(msg)
}

func sendStatistics(chatID int64, update chan string, msgFromMachine chan string, idReturn chan int64) {
	sendMessage(chatID, "Трансляция началась")
	//	http.Get("")
	for {
		select {
		case u := <-update:
			switch u {
			case "exit":
				idReturn <- chatID
				sendMessage(chatID, "Выход в обычный режим")
				return
			default:
				sendMessage(chatID, "print exit")

			}
		case reply := <-msgFromMachine:
			fmt.Println("Получил сообщение")
			sendMessage(chatID, reply)
			//case msgFromMachine <- "ok":
			//	continue
		}
	}
}

func deleteByValue(value int64, array []int64) []int64 {
	var arrayReturn []int64
	for _, a := range array {
		if value != a {
			arrayReturn = append(arrayReturn, a)
		}
	}
	return arrayReturn
}
