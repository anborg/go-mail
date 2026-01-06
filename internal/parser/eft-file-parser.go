package parser

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"muni/go-mail/internal/model"
)

func GetEftFromCSV(csvString string) (model.EftInfos, error) {
	re := regexp.MustCompile(`\r\r\n`)
	csvString = re.ReplaceAllString(csvString, "\r\n") //CRLF fix
	var eftInfos model.EftInfos
	r := csv.NewReader(strings.NewReader(csvString))
	r.FieldsPerRecord = -1 // optional
	r.TrimLeadingSpace = true
	r.Read() //skip header line

	var eftInfoArray []model.EftInfo
	count := 1
	for {
		// Read each record from csv
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return eftInfos, fmt.Errorf("Error processing csv record=" + strconv.Itoa(count) + ", row=" + strings.Join(row, "|") + ",  detail=" + err.Error())
		}
		var efinfo model.EftInfo
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

func substr(input string, start int, length int) string {
	asRunes := []rune(input)
	if start >= len(asRunes) {
		return ""
	}
	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func isNullOrEmpty(str string) bool {
	return len(str) == 0
}

func cleanInvoiceBlob(multiLineBlob string) (invoices []model.Invoice) {
	scanner := bufio.NewScanner(strings.NewReader(multiLineBlob))
	for scanner.Scan() {
		var csvRecord = scanner.Text()
		if !isNullOrEmpty(csvRecord) {
			var s = strings.TrimSpace(csvRecord)
			var invoiceNum = substr(s, 0, 30)
			var date = substr(s, 62, 8)
			var amount = substr(s, 70, 11)
			var eftRef = substr(s, 82, 9)
			var invoice = model.Invoice{
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

