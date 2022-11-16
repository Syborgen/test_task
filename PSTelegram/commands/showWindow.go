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

type ShowWindowCommand struct {
	StructureOfCommand
}

const ShowWindowCommandCallName = "show_window"

func NewShowWindowCommand() *ShowWindowCommand {
	return &ShowWindowCommand{
		StructureOfCommand: StructureOfCommand{
			CallName: ShowWindowCommandCallName,
		},
	}
}

func (c *ShowWindowCommand) Execute(message *tgbotapi.Message) error {
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

	var groupedTechWindows []datastructures.GroupedTechWindow
	json.Unmarshal(body, &groupedTechWindows)

	chatWithUser := message.Chat

	formattedTechWindowsTable := datastructures.CreateGroupedTechWindowTable(groupedTechWindows)
	err = tghelper.SendTextMessage(formattedTechWindowsTable, c.ChatToWrite.ID, c.Bot)
	if err != nil {
		tghelper.SendTextMessage("Send message error: "+err.Error(), chatWithUser.ID, c.Bot)
	}

	messageText := fmt.Sprintf("Technical windows table printed in chat @%s", c.ChatToWrite.UserName)
	tghelper.SendTextMessage(messageText, chatWithUser.ID, c.Bot)

	return nil
}
