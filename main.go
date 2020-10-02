package main

import (
	"flag"
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	//read cmdline
	var configFile string
	flag.StringVar(&configFile, "configFile", "config.yml", "Provid config file path,  e.g c:/my/dir/eftconf.yml")
	flag.Parse()
	//Read config
	// configFile :=
	var config Config
	if err := config.readConfig(configFile); err != nil {
		log.Fatalf("Error reading config file :", configFile, err)
	} else {
		log.Println(config)
	}

	//config log
	logconf := config.AppConfig.LumberjackLogConfig
	log.SetOutput(&lumberjack.Logger{Filename: logconf.Filename, MaxSize: logconf.MaxSize, MaxBackups: logconf.MaxBackups, MaxAge: logconf.MaxAge, Compress: logconf.Compress})

	//trigger -- filesToProcses
	files, err := filesMatch(config.FileProcessorConfig)
	if err != nil {
		log.Println(err)
	}
	log.Println(files)

	for _, inputFileInfo := range files {
		//process input csv file
		if err := process(inputFileInfo.path, config.MailServerConfig); err != nil {
			log.Println(err)
			postProcess(inputFileInfo, config.FileProcessorConfig.ErrorDir)
		} else {
			postProcess(inputFileInfo, config.FileProcessorConfig.DoneDir)
		}
	}

} //main
func postProcess(fileInfo InputFileInfo, targetPath string) {
	currentfileName := fileInfo.info.Name()
	newFileName := currentfileName + time.Now().Format("2020-01-31_154560.555")
	newfullName := targetPath + newFileName
	if e := os.Rename(fileInfo.path, newfullName); e != nil {
		log.Fatal("Error moving processed file to target dir: ", newfullName, e)
	} //
}

func process(filePath string, mailConf MailServerConfig) error {

	eftInfos, err := getEftFromCSV(filePath)
	if err != nil {
		log.Println("Error parsing input file:", filePath, err) //Go to next file. Email?
		return err
	}
	//send mails
	err1 := batchSendMail(mailConf, eftInfos)
	if err1 != nil {
		log.Println("Error sending emails for input file:", filePath, err1) //Go to next file. Email?
		return err1
	} else { //success
		log.Println("Processed done:", filePath, "Emails sent #: ", len(eftInfos.EftInfos)) //Go to next file. Email?
	}
	return nil
}
