package system

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"telega/lib"
)

func SendMessage(chatID int64, message string) {
	log.Printf("[chat Bot] %s", message)
	err := SendMessageParse(chatID, message)
	if err == nil {
		return
	} else {
		SendMessageWithoutParse(chatID, message)
	}

}
func SendMessageParse(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "markdown"
	_, err := lib.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func SendMessageWithoutParse(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := lib.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
