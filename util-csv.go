package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

// func main() {
// 	getEftFromCSV("cayinput556.csv")
// }

func getEftFromCSV(path string) (EftInfos, error) {
	var eftInfos EftInfos
	csvfile, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return eftInfos, err
	}
	r := csv.NewReader(csvfile)
	r.FieldsPerRecord = -1 // optional
	r.TrimLeadingSpace = true
	r.Read() //skip header line

	var eftInfoArray []EftInfo
	count := 1
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
		efinfo.Email = row[1]
		efinfo.SupplierId = row[2]
		efinfo.SupplierName = row[3]
		// var address := row[4]
		// var contact := row[5]
		efinfo.TransferAmount = row[6]
		efinfo.TransferDate = row[7]
		efinfo.InvoiceDetail = row[9]
		// efinfo.BankAccountNumber = row[2]

		eftInfoArray = append(eftInfoArray, efinfo)

		log.Println("eftInfo "+strconv.Itoa(count)+":", efinfo.SupplierName)
		count++
	}
	eftInfos.EftInfos = eftInfoArray
	// log.Println("eftInfos L",eftInfos)
	return eftInfos, nil
}
