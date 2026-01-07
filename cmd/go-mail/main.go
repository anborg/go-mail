package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"muni/go-mail/internal/config"
	"muni/go-mail/internal/processor"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// Set working directory to executable path
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	if err := os.Chdir(exPath); err != nil {
		log.Fatal(err)
	}

	//read cmdline
	var configFile string
	flag.StringVar(&configFile, "configFile", "config.yml", "Provid config file path,  e.g c:/my/dir/eftconf.yml")
	flag.Parse()
	//Read config
	var conf config.Config
	if err := conf.ReadConfig(configFile); err != nil {
		log.Fatalf("Error reading config file : %s, %v", configFile, err)
	} else {
		log.Println("Config: ", conf)
		log.Println("Check log file for details :", conf.AppConfig.LumberjackLogConfig.Filename)
	}

	//config log //TODO perhaps long to system.out sstem.error and then use infra to process log location/output?
	logconf := conf.AppConfig.LumberjackLogConfig
	
	// Check if log file needs rotation (older than 1 month or > 10MB)
	rotateLogIfNeeded(logconf.Filename)

	log.SetOutput(&lumberjack.Logger{Filename: logconf.Filename, MaxSize: logconf.MaxSize, MaxBackups: logconf.MaxBackups, MaxAge: logconf.MaxAge, Compress: logconf.Compress})
	fileProcessorConf := conf.FileProcessorConfig
	if err := processor.EnsureMandatoryDirsExist(fileProcessorConf.InputDir, fileProcessorConf.DoneDir, fileProcessorConf.ErrorDir); err != nil {
		log.Fatal(err)
	}

	//trigger -- filesToProcses
	files, err := processor.FilesMatch(fileProcessorConf)
	if err != nil {
		log.Println(err)
	}
	log.Println("Files found for processing: ", files)

	for _, inputFileInfo := range files {
		//process input csv file
		if err := processor.Process(inputFileInfo.Path, conf.MailServerConfig); err != nil {
			log.Println(err)
			processor.PostProcess(inputFileInfo, conf.FileProcessorConfig.ErrorDir)
		} else { // on error just move that file so other files in input dir can be processed
			//email eft processing error?
			processor.PostProcess(inputFileInfo, conf.FileProcessorConfig.DoneDir)
		}
	}

} //main

func rotateLogIfNeeded(filename string) {
	if info, err := os.Stat(filename); err == nil {
		// 1 month = 30 days roughly
		isOld := time.Since(info.ModTime()) > 30*24*time.Hour
		isLarge := info.Size() > 2*1024*1024 // 2MB

		if isOld || isLarge {
			ext := filepath.Ext(filename)
			name := strings.TrimSuffix(filename, ext)
			backupName := fmt.Sprintf("%s-%s%s", name, time.Now().Format("20060102-150405"), ext)

			if err := os.Rename(filename, backupName); err != nil {
				log.Printf("Failed to rotate existing log file: %v", err)
			} else {
				log.Printf("Rotated existing log file to: %s", backupName)
			}
		}
	}
}

