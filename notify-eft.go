package main

import (
	"flag"
	"log"

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

	//trigger -- filesTOprocses
	files, err := filesMatch(config.FileProcessorConfig)
	if err != nil {
		log.Println(err)
	}
	log.Println(files)

	for i, inputFileInfo := range files {
		//process input csv file

		eftInfos, err := getEftFromCSV(inputFileInfo.path)
		if err != nil {
			log.Println("Error parsing input file:", i, inputFileInfo.path, err) //Go to next file. Email?
		}
		//send mails
		err1 := batchSendMail(config.MailServerConfig, eftInfos)
		if err1 != nil {
			log.Println("Error sending emails for input file:", i, inputFileInfo.path, err1) //Go to next file. Email?
		}
		log.Println("Processed done:", i, inputFileInfo.path, "Emails sent #: ", len(eftInfos.EftInfos)) //Go to next file. Email?

	}

} //main
