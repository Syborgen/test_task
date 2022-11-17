package main

import (
	"PSTelegram/bothandler"
	cmdpkg "PSTelegram/commands"
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

func initCommands(configFileName string) (map[string]cmdpkg.Command, error) {
	commands := make(map[string]cmdpkg.Command)

	cfg, err := cmdpkg.CreateConfigFromFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("create config error: %w", err)
	}

	createCommand := cmdpkg.NewCreateCommand()
	createCommand.CommandURL = cfg.CreateCommandURL
	commands[cmdpkg.CreateCommandCallName] = createCommand

	addWindowCommand := cmdpkg.NewAddWindowCommand()
	addWindowCommand.CommandURL = cfg.AddWindowCommandURL
	commands[cmdpkg.AddWindowCommandCallName] = addWindowCommand

	showObjectCommand := cmdpkg.NewShowObjectCommand()
	showObjectCommand.CommandURL = cfg.ShowObjectCommandURL
	commands[cmdpkg.ShowObjectCommandCallName] = showObjectCommand

	showWindowCommand := cmdpkg.NewShowWindowCommand()
	showWindowCommand.CommandURL = cfg.ShowWindowCommandURL
	commands[cmdpkg.ShowWindowCommandCallName] = showWindowCommand

	showWindowSortCommand := cmdpkg.NewShowWindowSortCommand()
	showWindowSortCommand.CommandURL = cfg.ShowWindowSortCommandURL
	commands[cmdpkg.ShowWindowSortCommandCallName] = showWindowSortCommand

	showWindowAllCommand := cmdpkg.NewShowWindowAllCommand()
	showWindowAllCommand.CommandURL = cfg.ShowWindowAllCommandURL
	commands[cmdpkg.ShowWindowAllCommandCallName] = showWindowAllCommand

	return commands, nil
}
