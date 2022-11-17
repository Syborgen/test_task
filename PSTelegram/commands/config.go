package commands

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	CreateCommandURL         string `json:"create_command_url"`
	AddWindowCommandURL      string `json:"add_window_command_url"`
	ShowObjectCommandURL     string `json:"show_object_command_url"`
	ShowWindowCommandURL     string `json:"show_window_command_url"`
	ShowWindowSortCommandURL string `json:"show_window_sort_command_url"`
	ShowWindowAllCommandURL  string `json:"show_window_all_command_url"`
}

func CreateConfigFromFile(fileName string) (cfg config, err error) {
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
