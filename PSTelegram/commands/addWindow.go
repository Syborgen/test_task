package commands

import (
	datastructures "PSTelegram/dataStructures"
	"PSTelegram/helper"
	"PSTelegram/tghelper"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const AddWindowCommandCallName = "add_window"

type AddWindowCommandArguments struct {
	ObjectID int    `json:"object_id"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Action   string `json:"action"`
}

type AddWindowCommand struct {
	StructureOfCommand
}

func NewAddWindowCommand() *AddWindowCommand {
	return &AddWindowCommand{
		StructureOfCommand: StructureOfCommand{
			CallName:          AddWindowCommandCallName,
			ExpectedArguments: []string{"u", "d", "d", "s(query|proc)"},
		},
	}
}

func (c *AddWindowCommand) Execute(message *tgbotapi.Message) error {
	arguments := helper.ParseArguments(message.CommandArguments())
	err := c.ValidateArguments(arguments)
	if err != nil {
		return fmt.Errorf("arguments validation error: %w", err)
	}

	argumentsInJson, err := putArgumentsInJson(arguments)
	if err != nil {
		return fmt.Errorf("convertion argument to json error: %w", err)
	}

	res, err := http.Post(c.CommandURL, "application/json", bytes.NewBuffer(argumentsInJson))
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

	chatWithUser := message.Chat
	tghelper.SendTextMessage("New window added.", chatWithUser.ID, c.Bot)

	return nil
}

func putArgumentsInJson(arguments []string) ([]byte, error) {
	objectID, _ := strconv.Atoi(arguments[0])

	duration, err := datastructures.NewTimeRange(arguments[1], arguments[2])
	if err != nil {
		return nil, fmt.Errorf("time range creation error: %w", err)
	}

	err = duration.ConvertToServerTime()
	if err != nil {
		return nil, fmt.Errorf("convert to server time error: %w", err)
	}

	commandArguments := AddWindowCommandArguments{
		ObjectID: objectID,
		Start:    duration.Start,
		End:      duration.End,
		Action:   arguments[3],
	}

	argumentsInJson, err := json.Marshal(commandArguments)
	if err != nil {
		return nil, fmt.Errorf("arguments marshaling error: %w", err)
	}

	return argumentsInJson, nil
}
