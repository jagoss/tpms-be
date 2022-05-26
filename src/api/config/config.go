package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	ConfigFileName = "config.yml"
	ConfigFilePath = "configfiles/"
	DefaultPort    = "8080"
)

type GeneralConfiguration struct {
}

func (c *GeneralConfiguration) ReadConfiguration(filePath string) error {
	var source []byte
	var err error
	if source, err = ioutil.ReadFile(filePath); err != nil {
		return fmt.Errorf("error readign conf file: %s", err)
	}
	if err = yaml.Unmarshal(source, &c); err != nil {
		return fmt.Errorf("error Unmarshalling YAML configuration: %s", err)
	}
	return nil
}
