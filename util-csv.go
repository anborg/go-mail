package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func getEftFromCSV(path string) (EftInfos, error) {
	var eftInfos EftInfos
	csvfile, err := os.Open("inputfile.csv")
	if err != nil {
		return eftInfos, err
	}
	r := csv.NewReader(csvfile)
	r.Read() //skip header line

	var eftInfoArray []EftInfo

	for {
		// Read each record from csv
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return eftInfos, err
		}
		var efinfo EftInfo
		efinfo.Email = row[0]
		efinfo.CustomerName = row[1]
		efinfo.BankAccountNumber = row[2]
		efinfo.InvoiceNumber = row[3]
		efinfo.TransferAmount = row[4]
		efinfo.TransferDate = row[5]
		eftInfoArray = append(eftInfoArray, efinfo)

	}
	eftInfos.EftInfos = eftInfoArray
	log.Println(eftInfos)
	return eftInfos, nil
}
