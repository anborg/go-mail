package mail

import (
	"os"
	"strings"
	"testing"

	"muni/go-mail/internal/model"
)

func TestExecEftTemplate(t *testing.T) {
	// Change CWD to project root to find templates/DEFAULT.gohtml
	// We are in internal/mail, so we need to go up two levels
	wd, _ := os.Getwd()
	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(wd)

	var eftinfo = model.EftInfo{
		Email:          "demo@gmail.com",
		TodayDate:      "2020-09-01",
		SupplierName:   "Jane Inc",
		SupplierId:     "D004",
		TransferDate:   "2019-09-01",
		TransferAmount: "100.54",
		Invoices: []model.Invoice{
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

	bytesHtml, err := ExecEftTemplate(eftinfo)
	if err != nil {
		t.Fatalf("ExecEftTemplate failed: %v", err)
	}

	if bytesHtml.Len() == 0 {
		t.Error("Expected non-empty HTML output")
	}
	
	// Optional: Check if output contains expected values
	output := bytesHtml.String()
	ensureContains(t, output, "Jane Inc")
	ensureContains(t, output, "D004")
	ensureContains(t, output, "$12.11")
}

func ensureContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("Expected output to contain %q", substr)
	}
}
