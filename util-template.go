package main

import (
	"bytes"
	"text/template"
)

// func main() {
// 	var eftinfo = EftInfo{Email: "demo@gmail.com", TodayDate: "2020-09-01", CustomerName: "Jane Inc", BankAccountNumber: "111", InvoiceNumber: "2342", TransferDate: "2019-09-01", TransferAmount: "100.54"}
// 	// var bytesHtml bytes.Buffer

// 	bytesHtml, err := ExecEftTemplate(eftinfo)
// 	if err != nil {
// 		panic(err)
// 	} else {
// 		log.Println(bytesHtml.String())
// 	}
// }

func ExecEftTemplate(eftinfo EftInfo) (bytes.Buffer, error) { //eftInfo EftInfo
	return execTempate("templates/EFT_EMAIL_TEMPLATE.gohtml", eftinfo)
}
func execReceipt(eftinfo EftInfo) (bytes.Buffer, error) { //eftInfo EftInfo
	return execTempate("templates/EFT_EMAIL_TEMPLATE.gohtml", eftinfo)
}

func execTempate(path string, data interface{}) (bytes.Buffer, error) {
	var bytesHtml bytes.Buffer
	tpl, err := template.ParseFiles(path)
	if err != nil {
		return bytesHtml, err
	}

	if err := tpl.Execute(&bytesHtml, data); err != nil {
		return bytesHtml, err
	}
	return bytesHtml, nil
}

//EftInfos arry of eft info objects
type EftInfos struct {
	EftInfos []EftInfo `json:"eftInfos"`
}

//EftInfo object for one email
type EftInfo struct {
	TodayDate         string `json:"todayDate"`
	CustomerName      string `json:"customerName"`
	Email             string `json:"email"`
	BankAccountNumber string `json:"bankAccountNumber"`
	//EFT Payment detail - could be []
	InvoiceNumber    string `json:"invoiceNumber"`
	TransferDate     string `json:"transferDate"`
	TransferAmount   string `json:"transferAmount"`
	PaymentReference string `json:"paymentReference"`
}
