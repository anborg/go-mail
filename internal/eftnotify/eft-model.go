package eftnotify

import "os"

// EftInfos represents a collection of EFT information.
type EftInfos struct {
	EftInfos []EftInfo `json:"eftInfos"`
}

// EftInfo holds information for a single email notification.
type EftInfo struct {
	TodayDate      string    `json:"todayDate"`
	SupplierName   string    `json:"supplierName"`
	SupplierId     string    `json:"supplierId"`
	Email          string    `json:"email"`
	Invoices       []Invoice `json:"invoices"`
	InvoiceDetail  string    `json:"invoiceNumber"`
	TransferDate   string    `json:"transferDate"`
	TransferAmount string    `json:"transferAmount"`
}

// Invoice represents a single invoice within an EFT payment.
type Invoice struct {
	InvoiceNumber string `json:"invoiceNumber"`
	Date          string `json:"date"`
	Amount        string `json:"amount"`
	Ref           string `json:"ref"`
}

// InputFileInfo wraps path and os.FileInfo for processing.
type InputFileInfo struct {
	Path string
	Info os.FileInfo
}
