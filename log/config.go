package logger

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var (
	Logger   *zap.Logger
	filePath string
	level    string
)

func InitConfig() {
	config, err := LoadConfig("./config/journey.yaml")
	if err != nil {
		panic(err)
	}
	logConfig := config["logging"]
	filePath = logConfig.(map[string]interface{})["file"].(string)
	level = logConfig.(map[string]interface{})["level"].(string)
	NewLogger(level, filePath)
}

func LoadConfig(filename string) (map[string]interface{}, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config *map[string]interface{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return *config, nil
}
