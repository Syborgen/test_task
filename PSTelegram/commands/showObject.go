package commands

import (
	datastructures "PSTelegram/dataStructures"
	"PSTelegram/helper"
	"PSTelegram/tghelper"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const ShowObjectCommandCallName = "show_object"

type ShowObjectCommand struct {
	StructureOfCommand
}

func NewShowObjectCommand() *ShowObjectCommand {
	return &ShowObjectCommand{
		StructureOfCommand: StructureOfCommand{
			CallName: ShowObjectCommandCallName,
		},
	}
}

func (c *ShowObjectCommand) Execute(message *tgbotapi.Message) error {
	arguments := helper.ParseArguments(message.CommandArguments())
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

	err = c.checkError(body)
	if err != nil {
		return err
	}

	err = c.showResult(body, message.Chat)
	if err != nil {
		return err
	}

	return nil
}

func (c *ShowObjectCommand) showResult(body []byte, chatWithUser *tgbotapi.Chat) error {
	var objects []datastructures.Object
	err := json.Unmarshal(body, &objects)
	if err != nil {
		return fmt.Errorf("unmarshal body error: %w", err)
	}

	if len(objects) == 0 {
		tghelper.SendTextMessage("There is no objects in database.", c.ChatToWrite.ID, c.Bot)
		return nil
	}

	formattedObjectsTable := datastructures.CreateObjectTable(objects)
	err = tghelper.SendTextMessage(formattedObjectsTable, c.ChatToWrite.ID, c.Bot)
	if err != nil {
		return fmt.Errorf("send message error: %w", err)
	}

	messageText := fmt.Sprintf("Objets table printed in chat @%s", c.ChatToWrite.UserName)
	tghelper.SendTextMessage(messageText, chatWithUser.ID, c.Bot)

	return nil
}
