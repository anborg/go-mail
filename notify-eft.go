package main

import (
	"log"
	"os"
	"io"
)


func main() {
	//prepare log
	f, err := os.OpenFile("notify-eft.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)


	//Read config
	var config Config
	if err := readConfig("config.yml", &config); err != nil {
		log.Fatalf("Unmarshal: %v", err)
	} else {
		log.Println(config)
	}
	//process input csv file
	eftInfos, err := getEftFromCSV("cayinput556.csv")
	if err != nil {
		log.Fatalf("Error reading csv: %v", err)
	}
	//send mails
	if err := batchSendMail(config, eftInfos); err != nil {
		log.Fatalln(err)
	}
} //main
