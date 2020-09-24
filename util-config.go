package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//Config required to run
type Config struct {
	MailServer struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"` //not needed for m
	} `yaml:"mailserver"`
}

func readConfig(path string, cfg interface{}) error {
	yamlFile, err := ioutil.ReadFile(path)
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return err
	}
	return nil
}
