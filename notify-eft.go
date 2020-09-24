package main

import (
	"log"
)

func main() {
	//Read config
	var config Config
	if err := readConfig("config.yml", &config); err != nil {
		log.Fatalf("Unmarshal: %v", err)
	} else {
		log.Println(config)
	}
	//process input csv file
	eftInfos, err := getEftFromCSV("inputfile.csv")
	if err != nil {
		log.Fatalf("Error reading csv: %v", err)
	}
	//send mails
	if err := batchSendMail(config, eftInfos); err != nil {
		log.Fatalln(err)
	}
} //main
