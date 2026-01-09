package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

//Config required to run
type Config struct {
	AppConfig           AppConfig           `yaml:"app"`
	FileProcessorConfig FileProcessorConfig `yaml:"fileProcessor"`
	MailServerConfig    MailServerConfig    `yaml:"mailServer"`
}

// AppConfig App related config - eg logging
type AppConfig struct {
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
	Password string `yaml:"password"`
}

func (cfg *Config) ReadConfig(path string) (err error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return err
	}
	return cfg.Validate()
}

// Validate checks if the configuration is valid.
func (cfg *Config) Validate() error {
	if cfg.FileProcessorConfig.InputDir == "" {
		return fmt.Errorf("inputDir is mandatory")
	}
	if cfg.FileProcessorConfig.DoneDir == "" {
		return fmt.Errorf("doneDir is mandatory")
	}
	if cfg.FileProcessorConfig.ErrorDir == "" {
		return fmt.Errorf("errorDir is mandatory")
	}
	if cfg.MailServerConfig.Host == "" {
		return fmt.Errorf("mailServer host is mandatory")
	}
	return nil
}
