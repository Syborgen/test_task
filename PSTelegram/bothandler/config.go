package bothandler

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	BotToken        string `json:"bot_token"`
	ChatToWriteName string `json:"chat_to_write_name"`
}

func createConfigFromFile(fileName string) (cfg config, err error) {
	err = cfg.fillFromFile(fileName)
	return
}

func (cfg *config) fillFromFile(fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	err = json.Unmarshal(file, cfg)
	if err != nil {
		return fmt.Errorf("json unmarshal error: %w", err)
	}

	return nil
}
