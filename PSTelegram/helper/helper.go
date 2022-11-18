package helper

import (
	"regexp"
	"strings"
)

var argumentPattern = regexp.MustCompile(`( |^)('.+?'|\S+)`)

func ParseArguments(argumentsAsString string) []string {
	arguments := argumentPattern.FindAllString(argumentsAsString, -1)
	for i := range arguments {
		arguments[i] = strings.Trim(arguments[i], " '")
	}

	return arguments
}

func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
