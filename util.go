package main

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
)

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

func main12() {
	var x string = ` (I)8751                       PB20034 EST2 ASPHALT            31/07/20 $469974.60 001000006
		(I)8753                       PB20034 EST3 ASPHALT            15/08/20 $583537.82 001000006
		(I)8682                       PB20034 EST20101 ASPHALT        30/06/20$1353810.33 001000006
		(I)8685                       PB20034 JN20 ASPHALT INDEX      30/06/20 $125440.53 001000006
		`
	result := cleanInvoiceBlob(x)
	fmt.Println(result)
}

func isNullOrEmpty(str string) bool {
	//if str != nil {return false}
	return len(str) == 0
}

func cleanInvoiceBlob(multiLineBlob string) (invoices []Invoice) {
	scanner := bufio.NewScanner(strings.NewReader(multiLineBlob))
	//cleanInvoice = "INVOICE#          DATE        AMOUNT          EFTREF#\n"
	//var invoices []Invoice
	for scanner.Scan() {
		var csvRecord = scanner.Text()
		if !isNullOrEmpty(csvRecord) {
			var s = strings.TrimSpace(csvRecord)
			var invoiceNum = substr(s, 0, 30)
			var date = substr(s, 62, 8)
			var amount = substr(s, 70, 11)
			var eftRef = substr(s, 82, 9)
			var invoice = Invoice{
				InvoiceNumber: invoiceNum,
				Date:          date,
				Amount:        amount,
				Ref:           eftRef,
			}
			invoices = append(invoices, invoice)
			//             fmt.Println("invoicenum=" ,invoiceNum)
			//             fmt.Println("date=" + date)
			//             fmt.Println("amount=" + amount)
			//             fmt.Println("eftRef=" + eftRef)
			//             str := []string{invoiceNum, date, amount, eftRef}
			//             fmt.Println(str)
			//             cleanInvoice = cleanInvoice + strings.Join(str, "    ") + "\n"
		}
	}

	sort.SliceStable(invoices, func(i, j int) bool { // Order by Invoice#
		return invoices[i].InvoiceNumber < invoices[j].InvoiceNumber
	})
	return
}

//(I)8751
//(I)8751                       PB20034 EST2 ASPHALT            31/07/20 $469974.60 001000006
//PB20034 EST2 ASPHALT            31/07/20 $469974.60
