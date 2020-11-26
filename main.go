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
	var config Config
	if err := config.readConfig(configFile); err != nil {
		log.Fatalf("Error reading config file :", configFile, err)
	} else {
		log.Println("Config: ", config)
		log.Println("Check log file for details :", config.AppConfig.LumberjackLogConfig.Filename)
	}

	//config log //TODO perhaps long to system.out sstem.error and then use infra to process log location/output?
	logconf := config.AppConfig.LumberjackLogConfig
	log.SetOutput(&lumberjack.Logger{Filename: logconf.Filename, MaxSize: logconf.MaxSize, MaxBackups: logconf.MaxBackups, MaxAge: logconf.MaxAge, Compress: logconf.Compress})
	fileProcessorConf := config.FileProcessorConfig
	if err := ensureMandatoryDirsExist(fileProcessorConf.InputDir, fileProcessorConf.DoneDir, fileProcessorConf.ErrorDir); err != nil {
		log.Fatal(err)
	}

	//trigger -- filesToProcses
	files, err := filesMatch(fileProcessorConf)
	if err != nil {
		log.Println(err)
	}
	log.Println("Files found for processing: ", files)

	for _, inputFileInfo := range files {
		//process input csv file
		if err := process(inputFileInfo.path, config.MailServerConfig); err != nil {
			log.Println(err)
			postProcess(inputFileInfo, config.FileProcessorConfig.ErrorDir)
		} else { // on error just move that file so other files in input dir can be processed
			//email eft processing error?
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
		var errStr = "Error parsing input file:" + filePath + err.Error()
		emailInfo := EmailInfo{
			To:      mailConf.OpsUser,
			Cc:      mailConf.CcUser,
			Subject: "Error: Markham Notification - EFT",
			Body:    "Error while processing eft file : \n" + errStr,
		}
		_ = errorEmail(mailConf, emailInfo)
		log.Println(errStr) //Go to next file. Email?
		return err
	}
	//send mails
	err1 := batchSendMail(mailConf, eftInfos)
	if err1 != nil {
		log.Println("Error sending emails for input file:", filePath, err1) //Go to next file. Email?
		return err1
	}
	log.Println("Processed done:", filePath, "Emails sent #: ", len(eftInfos.EftInfos)) //Go to next file. Email?

	return nil
}
