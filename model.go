package main

//EftInfos arry of eft info objects
type EftInfos struct {
	EftInfos []EftInfo `json:"eftInfos"`
}

//EftInfo object for one email
type EftInfo struct {
	TodayDate         string `json:"todayDate"`
	SupplierName      string `json:"supplierName"`
	SupplierId        string `json:"supplierId"`
	Email             string `json:"email"`
	BankAccountNumber string `json:"bankAccountNumber"`
	//EFT Payment detail - could be []
	Invoices         []Invoice `json:"invoices"`
	InvoiceDetail    string    `json:"invoiceNumber"`
	TransferDate     string    `json:"transferDate"`
	TransferAmount   string    `json:"transferAmount"`
	PaymentReference string    `json:"paymentReference"`
}

type Invoice struct {
	InvoiceNumber string `json:"invoiceNumber"`
	Date          string `json:"date"`
	Amount        string `json:"amount"`
	Ref           string `json:"amount"`
}
