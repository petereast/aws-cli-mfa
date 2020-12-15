package config

import (
  "io/ioutil"
  "awsmfacli/src/filepaths"
  "encoding/json"
)

type Config struct {
	SecretAccessKey string
	AccessKeyId     string
	MfaDeviceArn    string
}


type ConfigReader interface {
	GetConfig() (Config, error)
}

type ConfigGetter struct {
	FilePaths filepaths.FilePaths
}

func (self ConfigGetter) GetConfig() (config Config, err error) {
	configPath, err := self.FilePaths.GetConfigLocation()
	jsonBytes, err := ioutil.ReadFile(configPath)

	err = json.Unmarshal(jsonBytes, &config)

	return
}
