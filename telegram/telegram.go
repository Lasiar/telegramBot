package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"fmt"
	"regexp"
	"strconv"
	"telega/model"
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

func Worker(msg chan tgbotapi.Update, json chan string) {
	var idListenGoodId []int64
	idArray := make(chan int64)
	chan1 := make(chan string)
	loop:
	for {
		select {
		case m := <-msg:
			for _, id := range idListenGoodId {
				if m.Message.Chat.ID == id {
					chan1 <- m.Message.Text
					continue loop
					}
			}
			v1 := regular(m, json, chan1, idArray)
			if v1 != 0 {
				idListenGoodId = append(idListenGoodId, v1)
			}
			case i := <-idArray:
				fmt.Println("Должен удалить: ", i)
				idListenGoodId = deleteByValue(i, idListenGoodId)
				fmt.Println(idListenGoodId)
		}
		fmt.Println(idListenGoodId)
	}
}

func regular(update tgbotapi.Update, json chan string, updateChan chan string, idArray chan int64) int64 {
	reply := "Не знаю что ответить"
	m := update.Message.Text
	switch {
	case m == "1":
		go test(update.Message.Chat.ID,updateChan, json, idArray)
		return  update.Message.Chat.ID
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



func test(chatID int64 ,update chan string, logs chan string, idArray chan int64) {
	sendMessage(chatID, "Трансляция началась")
	for {
		select {
		case u := <-update:
			switch u {
			case "exit":
				sendMessage(chatID, "Выход в обычный режим")
				idArray  <- chatID
				return
			default:
				sendMessage(chatID, "print exit")

			}
		case reply := <-logs:
			fmt.Println("Получил сообщение")
			sendMessage(chatID, reply)
		case logs <- "ok":
			continue
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