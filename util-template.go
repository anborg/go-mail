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

// ExecEftTemplate apply info on Template to create html
func ExecEftTemplate(eftinfo EftInfo) (bytes.Buffer, error) { //eftInfo EftInfo
	// eftTemplatePath := "templates/EFT_EMAIL_TEMPLATE.gohtml"
	eftTemplatePath := "templates/DEFAULT.gohtml"
	return execTempate(eftTemplatePath, eftinfo)
}
func execReceipt(eftTemplatePath string, eftinfo EftInfo) (bytes.Buffer, error) { //eftInfo EftInfo
	return execTempate(eftTemplatePath, eftinfo)
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
