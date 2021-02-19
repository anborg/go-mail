package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func main11() {
	// re := regexp.MustCompile(`\r\r\n`)

	path := "testfiles/AP566005_CRCRLF_present.csv"
	data, err := ioutil.ReadFile(path)
	input := string(data)

	//fmt.Println(string(data))
	if err != nil {
		log.Println("Error opening "+path, err)
		// return eftInfos, err
	}
	// eftInfo, err := getEftFromCSV("testfiles/cayinput566000.csv")
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(eftInfo)
	eftInfo, err := getEftFromCSV(input)
	if err != nil {
		log.Println(err)
	}
	log.Println(eftInfo)

}

func getEftFromCSV(csvString string) (EftInfos, error) {
	re := regexp.MustCompile(`\r\r\n`)
	csvString = re.ReplaceAllString(csvString, "\r\n")
	var eftInfos EftInfos
	r := csv.NewReader(strings.NewReader(csvString))
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
	//log.Println("eftInfos :", eftInfos)
	return eftInfos, nil
}
