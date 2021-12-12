package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App struct {
		Port string `yaml:"port"`
	} `yaml:"app"`
	Users map[string]struct {
		User            string `yaml:"user"`
		Name            string `yaml:"name"`
		Picture         string `yaml:"picture"`
		Key             string `yaml:"key"`
		CacheEnabled    bool   `yaml:"cache"`
		FSEventsEnabled bool   `yaml:"fsevents"`
		MediaRoot       string `yaml:"mediaRoot"`
	} `yaml:"users"`
}

func ReadConfig(path string) (Config, error) {
	var config Config
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
