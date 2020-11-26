package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func main2() {
	var eftinfo = EftInfo{
		Email:        "demo@gmail.com",
		TodayDate:    "2020-09-01",
		SupplierName: "Jane Inc",
		SupplierId:   "D004",
		//BankAccountNumber: "111",
		TransferDate:   "2019-09-01",
		TransferAmount: "100.54",
		//InvoiceDetail: "big bla",
		Invoices: []Invoice{
			{
				InvoiceNumber: "123",
				Date:          "2020-01-01",
				Amount:        "$12.11",
				Ref:           "00011",
			},
			{
				InvoiceNumber: "222",
				Date:          "2020-01-01",
				Amount:        "$12.11",
				Ref:           "00011",
			},
		},
	}
	// var bytesHtml bytes.Buffer

	bytesHtml, err := ExecEftTemplate(eftinfo)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(bytesHtml.String())
	}
}

// ExecEftTemplate apply info on Template to create html
func ExecEftTemplate(eftinfo EftInfo) (bytes.Buffer, error) { //eftInfo EftInfo
	eftTemplatePath := "templates/DEFAULT.gohtml"
	// 	eftTemplatePath := "templates/SAMPLE_TEMPLATE.txt" //for testing
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
