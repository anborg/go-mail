package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"muni/go-mail/internal/config"
	"muni/go-mail/internal/eftnotify"
	"muni/go-mail/internal/mail"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogWriter struct {
	target   io.Writer
	hostname string
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	prefix := fmt.Sprintf("%s %s ", time.Now().Format("2006-01-02-15-04-05-000"), l.hostname)
	_, err = l.target.Write([]byte(prefix))
	if err != nil {
		return 0, err
	}
	return l.target.Write(p)
}

func main() {
	// Set working directory to executable path
	if err := chdirToExecutable(); err != nil {
		log.Fatal(err)
	}

	configFile := flag.String("configFile", "config.yml", "Path to config file")
	flag.Parse()

	var conf config.Config
	if err := conf.ReadConfig(*configFile); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	setupLogging(conf.AppConfig.LumberjackLogConfig)

	// Initialize mail service and notifier
	mailService := mail.NewService(conf.MailServerConfig)
	notifier := eftnotify.NewNotifier(mailService, conf.MailServerConfig)

	processor := eftnotify.NewProcessor(conf.FileProcessorConfig, notifier)
	if err := processor.Initialize(); err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}

	files, err := processor.FilesMatch()
	if err != nil {
		log.Printf("Error finding files: %v", err)
	}

	if len(files) == 0 {
		log.Println("No files found for processing.")
		return
	}

	for _, file := range files {
		if err := processor.Process(file.Path); err != nil {
			log.Printf("Processing failed for %s: %v", file.Path, err)
			processor.PostProcess(file, conf.FileProcessorConfig.ErrorDir)
		} else {
			processor.PostProcess(file, conf.FileProcessorConfig.DoneDir)
		}
	}
}

func chdirToExecutable() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	return os.Chdir(filepath.Dir(ex))
}

func setupLogging(logconf config.LumberjackLogConfig) {
	rotateLogIfNeeded(logconf.Filename)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logconf.Filename,
		MaxSize:    logconf.MaxSize,
		MaxBackups: logconf.MaxBackups,
		MaxAge:     logconf.MaxAge,
		Compress:   logconf.Compress,
	}

	hostname, _ := os.Hostname()
	log.SetFlags(0) // Disable default timestamps
	log.SetOutput(&LogWriter{
		target:   lumberjackLogger,
		hostname: hostname,
	})
}

func rotateLogIfNeeded(filename string) {
	info, err := os.Stat(filename)
	if err != nil {
		return
	}

	isOld := time.Since(info.ModTime()) > 30*24*time.Hour
	isLarge := info.Size() > 2*1024*1024 // 2MB

	if isOld || isLarge {
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		backupName := filepath.Join(filepath.Dir(filename), name+"-"+time.Now().Format("20060102-150405")+ext)

		if err := os.Rename(filename, backupName); err != nil {
			log.Printf("Failed to rotate log file: %v", err)
		} else {
			log.Printf("Rotated log file to: %s", backupName)
		}
	}
}
