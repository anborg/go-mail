package config

import (
    "testing"
)

func TestReadConfig(t *testing.T) {
    configFilePath := "config-demo.yml"
    var conf Config
    
    // Test: Read configuration from file
    err := conf.ReadConfig(configFilePath)
    if err != nil {
        t.Fatalf("Failed to read config file: %v", err)
    }

    // Assert: App Config
    expectedLogFile := "logs/notify-eft/eft.log"
    if conf.AppConfig.LumberjackLogConfig.Filename != expectedLogFile {
        t.Errorf("Expected Log Filename to be '%s', got '%s'", expectedLogFile, conf.AppConfig.LumberjackLogConfig.Filename)
    }

    if conf.AppConfig.LumberjackLogConfig.MaxSize != 500 {
        t.Errorf("Expected MaxSize to be 500, got %d", conf.AppConfig.LumberjackLogConfig.MaxSize)
    }

    if !conf.AppConfig.LumberjackLogConfig.Compress {
        t.Errorf("Expected Compress to be true, got false")
    }

    // Assert: File Processor Config
    expectedGlob := "AP566*"
    if conf.FileProcessorConfig.GlobPath != expectedGlob {
        t.Errorf("Expected GlobPath to be '%s', got '%s'", expectedGlob, conf.FileProcessorConfig.GlobPath)
    }

    if conf.FileProcessorConfig.OlderThanSeconds != 2 {
        t.Errorf("Expected OlderThanSeconds to be 2, got %d", conf.FileProcessorConfig.OlderThanSeconds)
    }
    
    expectedInputDir := "C:\\data\\temp\\go-mail\\input\\"
    if conf.FileProcessorConfig.InputDir != expectedInputDir {
         t.Errorf("Expected InputDir to be '%s', got '%s'", expectedInputDir, conf.FileProcessorConfig.InputDir)
    }

    // Assert: Mail Server Config
    expectedOpsUser := "appsalert@google.com"
    if conf.MailServerConfig.OpsUser != expectedOpsUser {
        t.Errorf("Expected OpsUser to be '%s', got '%s'", expectedOpsUser, conf.MailServerConfig.OpsUser)
    }
    
     if conf.MailServerConfig.Port != 1025 {
        t.Errorf("Expected Port to be 1025, got %d", conf.MailServerConfig.Port)
    }
}
