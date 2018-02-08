package telegram

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"regexp"
	"strconv"
	"telega/lib"
	"telega/model"
	"telega/system"
)

var idInfo = regexp.MustCompile(`\d`)

func ReceivingMessageTelegram(msg chan tgbotapi.Update) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := lib.Bot.GetUpdatesChan(u)
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
		case u := <-sendWarning:
			system.SendMessageWithoutParse(u.ChatId, u.Message)
		case u := <-update:
			if u.Message.From.UserName != "" {
				log.Printf("[%s] %s", u.Message.From.UserName, u.Message.Text)
			} else {
				log.Printf("[%s %s] %s", u.Message.From.FirstName, u.Message.From.LastName, u.Message.Text)

			}
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
			v1, v2 := regular(u, broadcastBad, msgForBadListen, chatIdBadReturn, chatIdGoodReturn, broadcastGood, msgForGoodListen)
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
			idListenBadId = system.DeleteByValue(i, idListenBadId)
		case idGoodReturn := <-chatIdGoodReturn:
			fmt.Println(idGoodReturn)
			idListenGooodId = system.DeleteByValue(idGoodReturn, idListenGooodId)
		}
	}
}

func regular(update tgbotapi.Update, msgFromMachine chan string, msgForListen chan string, idGoodReturn chan int64, idBadReturn chan int64, goodJson chan lib.GoodJson, msgForGoodListen chan string) (int64, int64) {
	reply := "Не знаю что ответить"
	m := update.Message
	switch {
	case m.Command() == "listen":
		var js lib.RequestGoodStatistic
		if m.CommandArguments() == "" {
			reply = "Введите после команды id-машины"
			break
		}
		js.ChatId = m.Chat.ID
		id := m.CommandArguments()
		v1, err := strconv.Atoi(id)
		if err != nil {
			reply = fmt.Sprint(err)
			break
		}
		js.Point = append(js.Point, v1)
		model.InitialGoodStatistic(js)
		go consumerGoodStatistics(update.Message.Chat.ID, msgForGoodListen, goodJson, idBadReturn)
		return 0, update.Message.Chat.ID
	case m.Command() == "bad":
		go consumerBadStatistics(update.Message.Chat.ID, msgForListen, msgFromMachine, idGoodReturn)
		return update.Message.Chat.ID, 0
	case m.Command() == "list":
		reply = ""
		keys, _ := model.List()
		for _, key := range keys {
			reply = reply + fmt.Sprint(key, "; ")
		}
	case m.Command() == "count":
		count, err := model.CountAllQuery()
		if err != nil {
			log.Println(err)
			reply = "ошибка"
			break
		}
		reply = strconv.Itoa(count)
	case m.Command() == "point_today":
		reply = ""
		keys, _ := model.CountToDayQuery()
		for _, key := range keys {
			reply = reply + fmt.Sprint(key, "\n")
		}
	case idInfo.MatchString(m.Command()):
		infoPoint, _ := model.InfoPoint(m.Command())
		if infoPoint.Success {
			reply = fmt.Sprintf("ip: *%v*; user info: *%v*", infoPoint.Ip, infoPoint.UserAgent)
		} else {
			reply = "такой машины нет"
		}
	}

	system.SendMessage(update.Message.Chat.ID, reply)
	return 0, 0
}

func consumerBadStatistics(chatID int64, update chan string, msgFromMachine chan string, idReturn chan int64) {
	system.SendMessage(chatID, "Трансляция началась \n для выхода напишите: *exit*")
	for {
		select {
		case u := <-update:
			switch u {
			case "exit":
				idReturn <- chatID
				system.SendMessage(chatID, "Выход в обычный режим")
				return
			default:
				system.SendMessage(chatID, "Для выхода напишите exit")
			}
		case reply := <-msgFromMachine:
			system.SendMessage(chatID, reply)
		}
	}
}

func consumerGoodStatistics(chatID int64, update chan string, goodJson chan lib.GoodJson, idReturn chan int64) {
	system.SendMessage(chatID, "Трансляция началась \n для выхода напишите: *exit*")
	for {
		select {
		case u := <-update:
			switch u {
			case "exit":
				idReturn <- chatID
				system.SendMessage(chatID, "Выход в обычный режим")
				return
			default:
				system.SendMessage(chatID, "Для выхода напишите exit")

			}
		case reply := <-goodJson:
			system.SendMessage(chatID, fmt.Sprint(reply))
		}
	}
}
