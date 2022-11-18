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

const ShowWindowAllCommandCallName = "show_window_all"

type ShowWindowAllCommand struct {
	StructureOfCommand
}

func NewShowWindowAllCommand() *ShowWindowAllCommand {
	return &ShowWindowAllCommand{
		StructureOfCommand: StructureOfCommand{
			CallName: ShowWindowAllCommandCallName,
		},
	}
}

func (c *ShowWindowAllCommand) Execute(message *tgbotapi.Message) error {
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

func (c *ShowWindowAllCommand) showResult(body []byte, chatWithUser *tgbotapi.Chat) error {
	var techWindows []datastructures.TechWindow
	err := json.Unmarshal(body, &techWindows)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	if len(techWindows) == 0 {
		tghelper.SendTextMessage("There is no tech windows in database.", c.ChatToWrite.ID, c.Bot)
		return nil
	}

	for i := range techWindows {
		err := techWindows[i].Duration.ConvertToLocalTime()
		if err != nil {
			return fmt.Errorf("time convertation error: %w", err)
		}
	}

	formattedTechWindowsTable := datastructures.CreateTechWindowTable(techWindows)
	err = tghelper.SendTextMessage(formattedTechWindowsTable, c.ChatToWrite.ID, c.Bot)
	if err != nil {
		return fmt.Errorf("send message error: %w", err)
	}

	messageText := fmt.Sprintf("Technical windows table printed in chat @%s", c.ChatToWrite.UserName)
	tghelper.SendTextMessage(messageText, chatWithUser.ID, c.Bot)

	return nil
}
