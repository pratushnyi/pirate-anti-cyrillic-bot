package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	updateConfigTimeoutString := os.Getenv("UPDATE_CONFIG_TIMEOUT")

	botState := "enabled"

	if botToken == "" {
		panic("Environment variable TELEGRAM_BOT_TOKEN should be filled")
	}

	if updateConfigTimeoutString == "" {
		updateConfigTimeoutString = "1"
	}

	updateConfigTimeout, err := strconv.Atoi(updateConfigTimeoutString)
	if err != nil {
		panic("Environment variable UPDATE_CONFIG_TIMEOUT should be an integer")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = updateConfigTimeout

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
			case "enable":
				botState = "enabled"
				botMessage.Text = "Bot is enabled now"
			case "disable":
				botState = "disabled"
				botMessage.Text = "Bot is disabled now"
			case "status":
				botMessage.Text = fmt.Sprintf("Bot is %s now", botState)
			}

			if botMessage.Text != "" {
				if _, err := bot.Send(botMessage); err != nil {
					log.Panic(err)
				}
			}

			continue
		}

		if botState == "disabled" {
			log.Println("Bot is disabled now")
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
