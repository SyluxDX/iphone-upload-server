package utils

import (
	"encoding/json"
	"os"
	"path"
)

// Configurations json
type Configurations struct {
	ServerURL    string `json:"serverUrl"`
	ServerPort   int    `json:"serverPort"`
	UploadFolder string `json:"uploadFolder"`
}

// GetConfigs read and parse configurations json file
func GetConfigs() (Configurations, error) {
	// read file
	fdata, err := os.ReadFile(path.Join(".", "config.json"))
	if err != nil {
		return Configurations{}, err
	}
	// json data
	var config Configurations
	// unmarshall it
	err = json.Unmarshal(fdata, &config)
	if err != nil {
		return Configurations{}, err
	}
	return config, nil
}
