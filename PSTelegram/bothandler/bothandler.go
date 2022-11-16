package bothandler

import (
	"PSTelegram/commands"
	"PSTelegram/tghelper"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	Bot            *tgbotapi.BotAPI
	ChatForResults tgbotapi.Chat
	Commands       map[string]commands.Command
	Updates        tgbotapi.UpdatesChannel
}

func New(configFileName string, commands map[string]commands.Command) (*BotHandler, error) {
	fmt.Printf("Bot use \"%s\" file as config\n", configFileName)

	cfg, err := createConfigFromFile(configFileName)
	if err != nil {
		fmt.Println(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, errors.New("bot connection error")
	}

	fmt.Println("Connected to account", bot.Self.FirstName)
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	chatToWrite, err := tghelper.FindChatByName(cfg.ChatToWriteName, bot)
	if err != nil {
		return nil, fmt.Errorf("find chat by name error: %w", err)
	}

	for _, command := range commands {
		command.SetExecutor(bot)
		command.SetChatToWrite(&chatToWrite)
	}

	return &BotHandler{
		Bot:            bot,
		Updates:        updates,
		ChatForResults: chatToWrite,
		Commands:       commands,
	}, nil
}

func (botHandler *BotHandler) Start() {
	for update := range botHandler.Updates {
		go botHandler.handleOneUpdate(update)
	}
}

func (botHandler *BotHandler) handleOneUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatWithUserID := update.Message.Chat.ID
	if !update.Message.IsCommand() {
		botHandler.SendTextMessage("It's not a command.", chatWithUserID)
		return
	}

	command, ok := botHandler.Commands[update.Message.Command()]
	if !ok {
		botHandler.SendTextMessage("Unknown command.", chatWithUserID)
		return
	}

	err := command.Execute(update.Message)
	if err != nil {
		errorMessageText := fmt.Sprintf("Command execution error: %v", err)
		botHandler.SendTextMessage(errorMessageText, chatWithUserID)
		return
	}
}

func (botHandler *BotHandler) SendTextMessage(messageText string, to int64) error {
	return tghelper.SendTextMessage(messageText, to, botHandler.Bot)
}
