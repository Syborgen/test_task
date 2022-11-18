package tghelper

import (
	"PSTelegram/helper"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func FindChatByID(chatID int64, bot *tgbotapi.BotAPI) (tgbotapi.Chat, error) {
	chatConfig := tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	}

	chat, err := bot.GetChat(chatConfig)
	if err != nil {
		return tgbotapi.Chat{}, fmt.Errorf("can't find chat by ID='%d': %w", chatID, err)
	}

	return chat, nil
}

func FindChatByName(chatName string, bot *tgbotapi.BotAPI) (tgbotapi.Chat, error) {
	chatConfig := tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			SuperGroupUsername: "@" + chatName,
		},
	}

	chat, err := bot.GetChat(chatConfig)
	if err != nil {
		return tgbotapi.Chat{}, fmt.Errorf("can't find chat by name='%s': %w", chatName, err)
	}

	return chat, nil
}

func SendTextMessage(messageText string, to int64, bot *tgbotapi.BotAPI) error {
	maxMessageSize := 4096

	for i := 0; i < len(messageText); i += maxMessageSize {
		message := tgbotapi.NewMessage(
			to, messageText[i:i+helper.Min(maxMessageSize, len(messageText)-i)])

		if _, err := bot.Send(message); err != nil {
			chat, _ := FindChatByID(to, bot)
			return fmt.Errorf("can't send message to chat \"%s\" (@%s): %w",
				chat.Title,
				chat.UserName,
				err,
			)
		}
	}

	return nil
}
