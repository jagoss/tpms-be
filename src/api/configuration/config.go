package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	ConfigFileName = "test_config.yml"
	ConfigFilePath = "configfiles"
	DefaultPort    = "8080"
	Prod           = "prod"
	Dev            = "dev"
	Test           = "test"
)

type GeneralConfiguration struct {
	Database   DBConfig `yaml:"database"`
	CVModelURL string   `yaml:"cv_model_url"`
}

type DBConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"db_name"`
}

func (c *GeneralConfiguration) LoadConfiguration(filePath string) error {
	var source []byte
	var err error
	if source, err = ioutil.ReadFile(filePath); err != nil {
		return fmt.Errorf("error reading conf file: %v", err)
	}
	if err = yaml.Unmarshal(source, &c); err != nil {
		return fmt.Errorf("error Unmarshalling YAML configuration: %v", err)
	}
	return nil
}
