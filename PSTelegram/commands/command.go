package commands

import (
	datastructures "PSTelegram/dataStructures"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command interface {
	Execute(*tgbotapi.Message) error
	SetBotApi(*tgbotapi.BotAPI)
	SetChatToWrite(*tgbotapi.Chat)
}

type StructureOfCommand struct {
	CommandURL        string
	ChatToWrite       *tgbotapi.Chat
	Bot               *tgbotapi.BotAPI
	CallName          string
	ExpectedArguments []string
	// ExpectedArguments - слайс строк вида ["i", "u", "s(desc|asc)", "d"], где:
	// - количество элементов слайса - это количество аргументов команды
	// - буква определяет тип аргумента
	//   - i -> integer
	//   - s -> string
	//   - d -> date
	//	 - u -> unsigned int (>0)
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
	for i, expectedArgument := range sc.ExpectedArguments {
		if havePossibleValues(expectedArgument) {
			possibleValues := getPossibleValuesOfArgument(expectedArgument)

			if !isOneOfPossibleValues(arguments[i], possibleValues) {
				return fmt.Errorf("invalid value of argument №%d: possible values: %v",
					i+1, possibleValues)
			}
		}

		switch getExpectedArgumentType(expectedArgument) {
		case "i":
			_, err := strconv.Atoi(arguments[i])
			if err != nil {
				return fmt.Errorf("argument №%d must be number", i+1)
			}

		case "s":
			return nil

		case "d":
			_, err := time.Parse(datastructures.TimeFormat, arguments[i])
			if err != nil {
				return fmt.Errorf("date parsing error (argument №%d) [pattern: 'YYYY-MM-DD HH:MI:SS']: %w", i+1, err)
			}

		case "u":
			value, err := strconv.Atoi(arguments[i])
			if err != nil {
				return fmt.Errorf("argument №%d must be number", i+1)
			}

			if value < 0 {
				return fmt.Errorf("argument №%d can't be negative", i+1)
			}

		}

	}

	return nil
}

func havePossibleValues(expectedArgument string) bool {
	return strings.Contains(expectedArgument, "|")
}

func getExpectedArgumentType(expectedArgument string) string {
	return expectedArgument[:1]
}

func isOneOfPossibleValues(argumentValue string, possibleValues []string) bool {
	for _, possibleValue := range possibleValues {
		if argumentValue == possibleValue {
			return true
		}
	}

	return false
}

func getPossibleValuesOfArgument(expectedArgument string) []string {
	possibleValuesInParentheses := expectedArgument[1:]
	possibleValuesWithSeparator := strings.Trim(possibleValuesInParentheses, "()")
	possibleValues := strings.Split(possibleValuesWithSeparator, "|")

	return possibleValues
}

func (sc *StructureOfCommand) SetBotApi(bot *tgbotapi.BotAPI) {
	sc.Bot = bot
}

func (sc *StructureOfCommand) SetChatToWrite(chat *tgbotapi.Chat) {
	sc.ChatToWrite = chat
}
