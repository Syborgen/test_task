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

type ShowWindowSortCommand struct {
	StructureOfCommand
}

const ShowWindowSortCommandCallName = "show_window_sort"

func NewShowWindowSortCommand() *ShowWindowSortCommand {
	return &ShowWindowSortCommand{
		StructureOfCommand: StructureOfCommand{
			CallName:          ShowWindowSortCommandCallName,
			ExpectedArguments: []string{"s(asc|desc)", "d", "d", "s(query|proc)"},
		},
	}
}

func (c *ShowWindowSortCommand) Execute(message *tgbotapi.Message) error {
	arguments := tghelper.ParseArguments(message.CommandArguments())
	err := c.ValidateArguments(arguments)
	if err != nil {
		return fmt.Errorf("arguments validation error: %w", err)
	}

	req, err := c.createGetRequest(arguments)
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

	var groupedTechWindows []datastructures.GroupedTechWindow
	json.Unmarshal(body, &groupedTechWindows)

	formattedTechWindowsTable := datastructures.CreateGroupedTechWindowTable(groupedTechWindows)
	tghelper.SendTextMessage(formattedTechWindowsTable, c.ChatToWrite.ID, c.Bot)

	chatWithUser := message.Chat
	messageText := fmt.Sprintf("Sorted technical windows table printed in chat @%s", c.ChatToWrite.UserName)
	tghelper.SendTextMessage(messageText, chatWithUser.ID, c.Bot)

	return nil
}

func (c *ShowWindowSortCommand) createGetRequest(arguments []string) (*http.Request, error) {
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
