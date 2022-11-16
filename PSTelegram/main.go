package main

import (
	"PSTelegram/bothandler"
	commandspkg "PSTelegram/commands"
	"fmt"
)

const configFileName = "config.json"

// sudo docker build --tag docker-pstelegram .
func main() {
	commands, err := initCommands(configFileName)
	if err != nil {
		fmt.Printf("Commands initialization error: %v\n", err)
		return
	}

	botHandler, err := bothandler.New(configFileName, commands)
	if err != nil {
		fmt.Printf("BotHandler creation error: %v\n", err)
		return
	}

	botHandler.Start()
}

func initCommands(configFileName string) (map[string]commandspkg.Command, error) {
	commands := make(map[string]commandspkg.Command)

	cfg, err := commandspkg.CreateConfigFromFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("create config error: %w", err)
	}

	createCommand := commandspkg.NewCreateCommand()
	createCommand.CommandURL = cfg.CreateCommandURL
	commands[commandspkg.CreateCommandCallName] = createCommand

	addWindowCommand := commandspkg.NewAddWindowCommand()
	addWindowCommand.CommandURL = cfg.AddWindowCommandURL
	commands[commandspkg.AddWindowCommandCallName] = addWindowCommand

	showObjectCommand := commandspkg.NewShowObjectCommand()
	showObjectCommand.CommandURL = cfg.ShowObjectCommandURL
	commands[commandspkg.ShowObjectCommandCallName] = showObjectCommand

	showWindowCommand := commandspkg.NewShowWindowCommand()
	showWindowCommand.CommandURL = cfg.ShowWindowCommandURL
	commands[commandspkg.ShowWindowCommandCallName] = showWindowCommand

	showWindowSortCommand := commandspkg.NewShowWindowSortCommand()
	showWindowSortCommand.CommandURL = cfg.ShowWindowSortCommandURL
	commands[commandspkg.ShowWindowSortCommandCallName] = showWindowSortCommand
	return commands, nil
}
