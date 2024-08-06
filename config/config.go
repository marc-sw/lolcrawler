package config

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
	"strings"
)

const (
	BASE_DIR            = "lol-tools"
	CRAWLER_DIR         = "crawler"
	CRAWLER_LOG_FILE    = "logs.txt"
	CRAWLER_DATA_FILE   = "data.sqlite"
	CRAWLER_CONFIG_FILE = "config.toml"
)

type MissingConfigField struct {
	message string
}

func NewMissingConfigField(fieldName, filepath string) MissingConfigField {
	return MissingConfigField{message: fmt.Sprintf("missing or empty field '%s' in '%s'", fieldName, filepath)}
}

func (e MissingConfigField) Error() string {
	return e.message
}

type Config struct {
	LogFile string

	RiotApi struct {
		Key    string `toml:"key"`
		Region string `toml:"region"`
	} `toml:"riotapi"`

	Crawler struct {
		StartName  string `toml:"start_name"`
		StartTag   string `toml:"start_tag"`
		DataSource string
	} `toml:"crawler"`
}

func prepareDataStructure(dataDir string) error {
	if _, err := os.Stat(dataDir); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(dataDir, 0750); err != nil {
			return err
		}
	}

	configFile := filepath.Join(dataDir, CRAWLER_CONFIG_FILE)
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(filepath.Join(dataDir, CRAWLER_CONFIG_FILE), []byte(strings.Join([]string{
			"[riotapi]",
			"key = ''",
			"region = ''",
			"",
			"[crawler]",
			"start_name = ''",
			"start_tag = ''",
		}, "\n")), 0660)
	} else {
		return err
	}
}

func Load() (Config, error) {
	config := Config{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}
	dataDir := filepath.Join(homeDir, BASE_DIR, CRAWLER_DIR)
	config.LogFile = filepath.Join(dataDir, CRAWLER_LOG_FILE)
	config.Crawler.DataSource = filepath.Join(dataDir, CRAWLER_DATA_FILE)

	if err = prepareDataStructure(dataDir); err != nil {
		return config, err
	}
	configFilePath := filepath.Join(dataDir, CRAWLER_CONFIG_FILE)
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("count")
		return config, err
	}
	if err = toml.Unmarshal(data, &config); err != nil {
		return config, err
	}
	if config.RiotApi.Key == "" {
		return config, NewMissingConfigField("riotapi key", configFilePath)
	}
	if config.RiotApi.Region == "" {
		return config, NewMissingConfigField("riotapi region", configFilePath)
	}
	if config.Crawler.StartName == "" {
		return config, NewMissingConfigField("crawler start_name", configFilePath)
	}
	if config.Crawler.StartTag == "" {
		return config, NewMissingConfigField("crawler start_tag", configFilePath)
	}
	return config, nil
}
