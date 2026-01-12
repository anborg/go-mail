package eftprocess

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"sort"
	"muni/go-mail/internal/utils"
	"strings"


)

func GetEftInfosFromCSV(csvString string) (EftInfos, error) {
	re := regexp.MustCompile(`\r\r\n`)
	csvString = re.ReplaceAllString(csvString, "\r\n") //CRLF fix
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
			break
		}
		if err != nil {
			return eftInfos, fmt.Errorf("error processing csv record=%d, row=%s, detail=%w", count, strings.Join(row, "|"), err)
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

		count++
	}
	eftInfos.EftInfos = eftInfoArray
	return eftInfos, nil
}



func cleanInvoiceBlob(multiLineBlob string) (invoices []Invoice) {
	scanner := bufio.NewScanner(strings.NewReader(multiLineBlob))
	for scanner.Scan() {
		var csvRecord = scanner.Text()
		if !utils.IsNullOrEmpty(csvRecord) {
			var s = strings.TrimSpace(csvRecord)
			var invoiceNum = utils.Substr(s, 0, 30)
			var date = utils.Substr(s, 62, 8)
			var amount = utils.Substr(s, 70, 11)
			var eftRef = utils.Substr(s, 82, 9)
			var invoice = Invoice{
				InvoiceNumber: invoiceNum,
				Date:          date,
				Amount:        amount,
				Ref:           eftRef,
			}
			invoices = append(invoices, invoice)
		}
	}

	sort.SliceStable(invoices, func(i, j int) bool { // Order by Invoice#
		return invoices[i].InvoiceNumber < invoices[j].InvoiceNumber
	})
	return
}

