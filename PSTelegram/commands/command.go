package commands

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command interface {
	Execute(*tgbotapi.Message) error
	SetExecutor(*tgbotapi.BotAPI)
	SetChatToWrite(*tgbotapi.Chat)
}

type StructureOfCommand struct {
	CommandURL        string
	ChatToWrite       *tgbotapi.Chat
	Bot               *tgbotapi.BotAPI
	CallName          string
	ExpectedArguments []string
	// ExpectedArguments - слайс строк вида ["i", "s", "s(desc|asc)", "d"], где:
	// - количество элементов - это количество аргументов
	// - буква определяет тип аргумента
	//   - i -> integer
	//   - s -> string
	//   - d -> date
	// - если в скобках указаны значения, разделенные "|", это означает,
	// что аргумент может принимать только одно из этих значений.
}

func (sc *StructureOfCommand) ValidateArguments(arguments []string) error {
	err := sc.validateArgumentsCount(arguments)
	if err != nil {
		return fmt.Errorf("arguments count validation error: %w", err)
	}

	err = sc.validateArgumentsValues(arguments)
	if err != nil {
		return fmt.Errorf("arguments values validation error: %w", err)
	}

	return nil
}

func (sc *StructureOfCommand) validateArgumentsCount(arguments []string) error {
	var expectedArgumentsCount = len(sc.ExpectedArguments)

	actualArgumentsCount := len(arguments)
	if actualArgumentsCount != expectedArgumentsCount {
		return fmt.Errorf("actual count of arguments is %d, but expected count is %d",
			actualArgumentsCount,
			expectedArgumentsCount,
		)
	}

	return nil
}

func (sc *StructureOfCommand) validateArgumentsValues(arguments []string) error {
	return nil
}

func (sc *StructureOfCommand) SetExecutor(bot *tgbotapi.BotAPI) {
	sc.Bot = bot
}

func (sc *StructureOfCommand) SetChatToWrite(chat *tgbotapi.Chat) {
	sc.ChatToWrite = chat
}
