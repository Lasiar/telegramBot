package telegram

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"regexp"
	"strconv"
	"telega/model"
	"telega/lib"
)

var Bot *tgbotapi.BotAPI
var idInfo = regexp.MustCompile(`\d`)
var idListen = regexp.MustCompile(`listen \d`)

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

func Worker(update chan tgbotapi.Update, msgFromMachine chan string, goodJson chan lib.GoodJson, sendWarning chan lib.MessageChat) {
	var idListenBadId []int64
	var idListenGooodId []int64
	chatIdBadReturn := make(chan int64)
	chatIdGoodReturn := make(chan int64)
	broadcastBad := make(chan string)
	broadcastGood := make(chan lib.GoodJson)
	msgForBadListen := make(chan string)
	msgForGoodListen := make(chan string)

loop:
	for {
		select {
		case u := <- sendWarning:
			sendMessage(379572314, u.Message)
		case u := <-update:
			for _, id := range idListenBadId {
				if u.Message.Chat.ID == id {
					msgForBadListen <- u.Message.Text
					continue loop
				}

			}
			for _, id := range idListenGooodId {
				fmt.Println("in if")
				if u.Message.Chat.ID == id {
					msgForGoodListen <- u.Message.Text
					fmt.Println("send msg")
					continue loop
				}
			}
			v1, v2 := regular(u, broadcastBad, msgForBadListen, chatIdBadReturn, chatIdGoodReturn, broadcastGood,msgForGoodListen)
			if v1 != 0 {
				idListenBadId = append(idListenBadId, v1)
			}
			if v2 != 0 {
				idListenGooodId = append(idListenGooodId, v2)
				fmt.Println("add GoodID: ", v2)
			}
		case g := <-goodJson:
			for range idListenGooodId {
				broadcastGood <- g
				fmt.Println("Отправил")
			}
		case b := <-msgFromMachine:
			for range idListenBadId {
				broadcastBad <- b
			}
		case i := <-chatIdBadReturn:
			idListenBadId = deleteByValue(i, idListenBadId)
			fmt.Println(idListenBadId)
		case i := <-chatIdGoodReturn:
			idListenBadId = deleteByValue(i, idListenBadId)
		}
	}
}

func regular(update tgbotapi.Update, msgFromMachine chan string, msgForListen chan string, idGoodReturn chan int64, idBadReturn chan int64, goodJson chan lib.GoodJson, msgForGoodListen chan string) (int64, int64) {
	reply := "Не знаю что ответить"
	m := update.Message.Text
	switch {
	case idListen.MatchString(m):
		var js lib.RequestGoodStatistic
		js.ChatId = update.Message.Chat.ID
		//FIXME нужно сделать чтобы не только для одного айди было
		id := m[7:]
		v1, err := strconv.Atoi(id)
		if err != nil {
			reply = fmt.Sprint(err)
			break
		}
		js.Point = append(js.Point, v1)
		model.InitialGoodStatistic(js)
		go consumerGoodStatistics(update.Message.Chat.ID, msgForGoodListen, goodJson, idBadReturn)
		return  0, update.Message.Chat.ID
	case m == "bad":
		go consumerBadStatistics(update.Message.Chat.ID, msgForListen, msgFromMachine, idGoodReturn)
		return update.Message.Chat.ID, 0
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
	case idInfo.MatchString(m):
		infoPoint, _ := model.InfoPoint(m)
		if infoPoint.Success {
			reply = fmt.Sprintf("ip: *%v*; user info: *%v*", infoPoint.Ip, infoPoint.UserAgent)
		} else {
			reply = "такой машины нет"
		}
	}
	sendMessage(update.Message.Chat.ID, reply)
	return 0, 0
}

func sendMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "markdown"
	Bot.Send(msg)
}

func consumerBadStatistics(chatID int64, update chan string, msgFromMachine chan string, idReturn chan int64) {
	sendMessage(chatID, "Трансляция началась \n для выхода напишите: *exit*")
	for {
		select {
		case u := <-update:
			switch u {
			case "exit":
				idReturn <- chatID
				sendMessage(chatID, "Выход в обычный режим")
				return
			default:
				sendMessage(chatID, "Для выхода напишите exit")

			}
		case reply := <-msgFromMachine:
			sendMessage(chatID, reply)
		}
	}
}

func consumerGoodStatistics(chatID int64, update chan string, goodJson chan lib.GoodJson, idReturn chan int64) {
	sendMessage(chatID, "Трансляция началась \n для выхода напишите: *exit*")
	for {
		select {
		case u := <-update:
			switch u {
			case "exit":
				idReturn <- chatID
				sendMessage(chatID, "Выход в обычный режим")
				return
			default:
				sendMessage(chatID, "Для выхода напишите exit")

			}
		case reply := <-goodJson:
			sendMessage(chatID, fmt.Sprint(reply))
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
