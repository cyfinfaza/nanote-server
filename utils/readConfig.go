package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App struct { 
		Port			string `yaml:"port"`
		CacheDefault	string `yaml:"cacheDefault"`
	} `yaml:"app"`
	Users map[string]struct {
		User 			string `yaml:"user"`
		Name 			string `yaml:"name"`
		Picture 		string `yaml:"picture"`
		Key 			string `yaml:"key"`
		Cache			bool `yaml:"cache"`
		MediaRoot		string `yaml:"mediaRoot"`
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