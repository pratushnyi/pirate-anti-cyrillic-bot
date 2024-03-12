package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"regexp"
)

func main() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 1 // todo to config

	re := regexp.MustCompile(`(?i)[а-яё]`)

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		// ignore any updates that does not contain messages
		if update.Message == nil {
			continue
		}

		botMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			default:
				botMessage.Text = "I don't know that command yet"
			}

			if _, err := bot.Send(botMessage); err != nil {
				log.Panic(err)
			}

			continue
		}

		msgText := update.Message.Text

		if msgText == "" {
			msgText = update.Message.Caption
		}

		if re.MatchString(msgText) {
			deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)

			if _, err := bot.Request(deleteMsg); err != nil {
				log.Panic(err)
			}

			botMsgText := fmt.Sprintf("@%s's message has been deleted because it contains cyrillic letters", update.Message.From.UserName)

			botMessage := tgbotapi.NewMessage(update.Message.Chat.ID, botMsgText)

			if _, err := bot.Send(botMessage); err != nil {
				log.Panic(err)
			}
		}
	}
}
