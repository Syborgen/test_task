package commands

import (
	"PSTelegram/tghelper"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CreateCommandArguments struct {
	Objects int `json:"objects"`
	Windows int `json:"windows"`
}

type CreateCommand struct {
	StructureOfCommand
}

const CreateCommandCallName = "create"

func NewCreateCommand() *CreateCommand {
	return &CreateCommand{
		StructureOfCommand: StructureOfCommand{
			CallName:          CreateCommandCallName,
			ExpectedArguments: []string{"u", "u"},
		},
	}
}

func (c *CreateCommand) Execute(message *tgbotapi.Message) error {
	arguments := tghelper.ParseArguments(message.CommandArguments())
	err := c.ValidateArguments(arguments)
	if err != nil {
		return fmt.Errorf("arguments validation error: %w", err)
	}

	objectsCount, _ := strconv.Atoi(arguments[0])
	windowsCount, _ := strconv.Atoi(arguments[1])
	commandArguments := CreateCommandArguments{
		Objects: objectsCount,
		Windows: windowsCount,
	}

	jsonArguments, err := json.Marshal(commandArguments)
	if err != nil {
		return fmt.Errorf("arguments marshaling error: %w", err)
	}

	res, err := http.Post(c.CommandURL, "application/json", bytes.NewBuffer(jsonArguments))
	if err != nil {
		return fmt.Errorf("http request error: %w", err)
	}

	defer res.Body.Close()

	chatWithUser := message.Chat
	tghelper.SendTextMessage("Data generated.", chatWithUser.ID, c.Bot)

	return nil
}
