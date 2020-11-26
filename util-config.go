package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// func main() {
// 	var conf Config
// 	err := readConfig("config.yml", &conf)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	log.Println(conf)
// }

//Config required to run
type Config struct {
	AppConfig           AppConfig           `yaml:"app"`
	FileProcessorConfig FileProcessorConfig `yaml:"fileProcessor"`
	MailServerConfig    MailServerConfig    `yaml:"mailServer"`
}

// AppConfig App related config - eg logging
type AppConfig struct {
	LogDir              string              `yaml:"logDir"`
	LumberjackLogConfig LumberjackLogConfig `yaml:"lumberJackLogging"`
}

// LumberjackLogConfig Loggin config
type LumberjackLogConfig struct {
	Filename   string `yaml:"filename" json:"Filename,string"`
	MaxSize    int    `yaml:"maxSize" json:"MaxSize,int"`
	MaxBackups int    `yaml:"maxBackups" json:"MaxBackups,int"`
	MaxAge     int    `yaml:"maxAge" json:"MaxAge,int"`
	Compress   bool   `yaml:"compress" json:"Compress,bool"`
}

// FileProcessorConfig Files to pocess. Specify  which/where/when. which (glob), where (dir), when (old/young)
type FileProcessorConfig struct {
	GlobPath         string `yaml:"globPath"`
	OlderThanSeconds int    `yaml:"olderThanSeconds"`
	InputDir         string `yaml:"inputDir"`
	DoneDir          string `yaml:"doneDir"`
	ErrorDir         string `yaml:"errorDir"`
}

// MailServerConfig smtp server details
type MailServerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	CcUser   string `yaml:"ccUser"`
	OpsUser  string `yaml:"opsUser"`
	Password string `yaml:"password"` //not needed for m
}

func (cfg *Config) readConfig(path string) (err error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return err
	}
	return
}
