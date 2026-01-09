package utils

import (
	"os"
	"strings"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	// Mock structures that match what DEFAULT.gohtml expects
	type MockInvoice struct {
		InvoiceNumber string
		Date          string
		Amount        string
		Ref           string
	}
	type MockEftInfo struct {
		TodayDate      string
		SupplierName   string
		SupplierId     string
		Email          string
		Invoices       []MockInvoice
		InvoiceDetail  string
		TransferDate   string
		TransferAmount string
	}

	// Create a temporary template file mimicking a simplified version of DEFAULT.gohtml
	content := `
ATTN: {{.SupplierName}}
Email: {{.Email}}
ID: {{.SupplierId}}
Transfer Date: {{.TransferDate}}
Transfer Amount: {{.TransferAmount}}
Invoices:
{{range .Invoices}}- {{.InvoiceNumber}}: {{.Amount}}
{{end}}
Date: {{.TodayDate}}`

	tmpfile, err := os.CreateTemp("", "test_template_*.gohtml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Mock data mimicking AP566005_CRCRLF_present.csv
	data := MockEftInfo{
		TodayDate:      "2026-01-09",
		SupplierName:   "D CRUPI & SONS LIMITED",
		SupplierId:     "D0004",
		Email:          "pnataraj@markham.ca",
		TransferDate:   "01 DEC 2020",
		TransferAmount: "$286632.55",
		Invoices: []MockInvoice{
			{InvoiceNumber: "(I)8791", Date: "31/08/20", Amount: "$194518.00", Ref: "001000007"},
			{InvoiceNumber: "(I)8792", Date: "31/08/20", Amount: "$92114.55", Ref: "001000007"},
		},
	}

	// Render
	buf, err := RenderTemplate(tmpfile.Name(), data)
	if err != nil {
		t.Errorf("RenderTemplate failed: %v", err)
	}

	result := buf.String()
	expectedSubstrings := []string{
		"ATTN: D CRUPI & SONS LIMITED",
		"Email: pnataraj@markham.ca",
		"ID: D0004",
		"Transfer Date: 01 DEC 2020",
		"Transfer Amount: $286632.55",
		"- (I)8791: $194518.00",
		"- (I)8792: $92114.55",
		"Date: 2026-01-09",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(result, s) {
			t.Errorf("Expected result to contain %q, but it didn't.\nResult:\n%s", s, result)
		}
	}
}


func TestRenderTemplate_FileNotFound(t *testing.T) {
	_, err := RenderTemplate("non_existent_file.gohtml", nil)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
