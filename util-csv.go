package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// func main() {
// 	getEftFromCSV("cayinput556.csv")
// }

func getEftFromCSV(path string) (EftInfos, error) {
	var eftInfos EftInfos
	log.Println("Opening file " + path)
	csvfile, err := os.Open(path)
	if err != nil {
		log.Println("Error opening "+path, err)
		return eftInfos, err
	}
	defer csvfile.Close()
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
			break //Not an error? Just ignore?
		}
		if err != nil {
			return eftInfos, fmt.Errorf("Error processing csv record=" + strconv.Itoa(count) + ", row=" + strings.Join(row, "|") + ",  detail=" + err.Error())
		}
		var efinfo EftInfo
		efinfo.Email = row[1]
		efinfo.SupplierId = row[2]
		efinfo.SupplierName = row[3]
		// var address := row[4]
		// var contact := row[5]
		efinfo.TransferAmount = row[6]
		efinfo.TransferDate = row[7]
		efinfo.Invoices = cleanInvoiceBlob(row[9])
		// efinfo.BankAccountNumber = row[2]

		eftInfoArray = append(eftInfoArray, efinfo)

		log.Println("eftInfo "+strconv.Itoa(count)+":", efinfo.SupplierName)
		count++
	}
	eftInfos.EftInfos = eftInfoArray
	// log.Println("eftInfos L",eftInfos)
	return eftInfos, nil
}
