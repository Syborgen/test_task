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

const ShowWindowSortCommandCallName = "show_window_sort"

type ShowWindowSortCommand struct {
	StructureOfCommand
}

func NewShowWindowSortCommand() *ShowWindowSortCommand {
	return &ShowWindowSortCommand{
		StructureOfCommand: StructureOfCommand{
			CallName:          ShowWindowSortCommandCallName,
			ExpectedArguments: []string{"s(asc|desc)", "d", "d", "s(query|proc)"},
		},
	}
}

func (c *ShowWindowSortCommand) Execute(message *tgbotapi.Message) error {
	arguments := helper.ParseArguments(message.CommandArguments())
	err := c.ValidateArguments(arguments)
	if err != nil {
		return fmt.Errorf("arguments validation error: %w", err)
	}

	req, err := c.createGetRequestWithArguments(arguments)
	if err != nil {
		return fmt.Errorf("http reqest creation error: %w", err)
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

func (c *ShowWindowSortCommand) showResult(body []byte, chatWithUser *tgbotapi.Chat) error {
	var groupedTechWindows []datastructures.GroupedTechWindow
	err := json.Unmarshal(body, &groupedTechWindows)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	if len(groupedTechWindows) == 0 {
		tghelper.SendTextMessage("There is no tech windows in database.", c.ChatToWrite.ID, c.Bot)
		return nil
	}

	formattedTechWindowsTable := datastructures.CreateGroupedTechWindowTable(groupedTechWindows)
	err = tghelper.SendTextMessage(formattedTechWindowsTable, c.ChatToWrite.ID, c.Bot)
	if err != nil {
		return fmt.Errorf("send message error: %w", err)
	}

	messageText := fmt.Sprintf("Sorted technical windows table printed in chat @%s", c.ChatToWrite.UserName)
	tghelper.SendTextMessage(messageText, chatWithUser.ID, c.Bot)

	return nil
}

func (c *ShowWindowSortCommand) createGetRequestWithArguments(
	arguments []string) (*http.Request, error) {
	req, err := http.NewRequest("GET", c.CommandURL, nil)
	if err != nil {
		return nil, fmt.Errorf("http request creation error: %w", err)
	}

	query := req.URL.Query()
	query.Add("sort", arguments[0])

	timeBorders, err := datastructures.NewTimeRange(arguments[1], arguments[2])
	if err != nil {
		return nil, fmt.Errorf("time range creation error: %w", err)
	}

	err = timeBorders.ConvertToServerTime()
	if err != nil {
		return nil, fmt.Errorf("convert to server time error: %w", err)
	}

	query.Add("start", timeBorders.Start)
	query.Add("end", timeBorders.End)
	query.Add("action", arguments[3])

	req.URL.RawQuery = query.Encode()

	fmt.Println("Sent request:", req.URL.String())

	return req, nil
}
