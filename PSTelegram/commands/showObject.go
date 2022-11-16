package commands

import (
	datastructures "PSTelegram/dataStructures"
	"PSTelegram/tghelper"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ShowObjectCommand struct {
	StructureOfCommand
}

const ShowObjectCommandCallName = "show_object"

func NewShowObjectCommand() *ShowObjectCommand {
	return &ShowObjectCommand{
		StructureOfCommand: StructureOfCommand{
			CallName: ShowObjectCommandCallName,
		},
	}
}

func (c *ShowObjectCommand) Execute(message *tgbotapi.Message) error {
	arguments := tghelper.ParseArguments(message.CommandArguments())
	err := c.ValidateArguments(arguments)
	if err != nil {
		return fmt.Errorf("arguments validation error: %w", err)
	}

	req, err := http.NewRequest("GET", c.CommandURL, nil)
	if err != nil {
		return fmt.Errorf("http request creation error: %w", err)
	}

	res, err := http.Get(req.URL.String())
	if err != nil {
		return fmt.Errorf("http request error: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read request body error: %w", err)
	}

	var objects []datastructures.Object
	json.Unmarshal(body, &objects)

	formattedObjectsTable := datastructures.CreateObjectTable(objects)
	err = tghelper.SendTextMessage(formattedObjectsTable, c.ChatToWrite.ID, c.Bot)
	if err != nil {
		return fmt.Errorf("send message error: %w", err)
	}

	chatWithUser := message.Chat
	messageText := fmt.Sprintf("Objets table printed in chat @%s", c.ChatToWrite.UserName)
	tghelper.SendTextMessage(messageText, chatWithUser.ID, c.Bot)

	return nil
}
