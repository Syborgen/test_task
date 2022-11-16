package tghelper

import (
	"fmt"
	"regexp"
	"strings"

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
		message := tgbotapi.NewMessage(to, messageText[i:i+min(maxMessageSize, len(messageText)-i)])

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

var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func ParseArguments(argumentsAsString string) []string {
	if argumentsAsString == "" {
		return []string{}
	}

	arguments := strings.Split(argumentsAsString, " ")
	for i := 0; i < len(arguments); i++ {
		if isDate(arguments[i]) {
			arguments[i] += " " + arguments[i+1]
			arguments = append(arguments[:i+1], arguments[i+2:]...)
		}
	}

	return arguments
}

func isDate(str string) bool {
	return datePattern.MatchString(str)
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
