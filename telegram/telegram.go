package telegram

import (
	//	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	//	"regexp"
	//	"strconv"
	//	"telega/lib"
	//	"telega/model"
	"fmt"
	"regexp"
	"strconv"
	"telega/model"
//	"reflect"
)

var Bot *tgbotapi.BotAPI

/*func MainTtelegram(logs chan string) {
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
		go func(id int64) {
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
		} */ /*
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "markdown"
		bot.Send(msg)
	}
} */

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
	for {
		select {
		case u := <-update:
			switch u {
			case "exit":
				idArray  <- chatID
				return
			default:
				reply := "print exit"
				sendMessage(chatID, reply)

			}
		case reply := <-logs:
			fmt.Println("Жду получения")
			reply = <-logs
			fmt.Println("Получил")
			sendMessage(chatID, reply)
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