package eftnotify

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetEftFromCSV(t *testing.T) {
	// Locate the test file in the 'testfiles' directory, relative to this test file.
	// This test file is in internal/parser, so testfiles is ../../testfiles
	testFilePath := filepath.Join("..", "..", "testfiles", "cayinput566000.csv")

	// Read the file
	content, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", testFilePath, err)
	}

	// Parse CSV
	eftInfos, err := GetEftInfosFromCSV(string(content))
	if err != nil {
		t.Fatalf("GetEftFromCSV returned error: %v", err)
	}

	// Expected: 2 records (D CRUPI and AECOM)
	if len(eftInfos.EftInfos) != 2 {
		t.Errorf("Expected 2 EftInfo records, got %d", len(eftInfos.EftInfos))
	}

	// Verify first record (D CRUPI)
	if len(eftInfos.EftInfos) > 0 {
		info1 := eftInfos.EftInfos[0]
		if info1.SupplierId != "D0004" {
			t.Errorf("Record 1: Expected SupplierId 'D0004', got %q", info1.SupplierId)
		}
		if info1.TransferAmount != "$2532763.28" {
			t.Errorf("Record 1: Expected TransferAmount '$2532763.28', got %q", info1.TransferAmount)
		}
		
		// Verify invoices for first record
		// The blob contains 4 invoices.
		if len(info1.Invoices) != 4 {
			t.Errorf("Record 1: Expected 4 invoices, got %d", len(info1.Invoices))
		}
		
		// Check one invoice details (e.g. sorted order check or specific item)
		// We know from previous manual check the order: 8682, 8685, 8751, 8753
		expectedInvNumPrefix := "(I)8682"
		if len(info1.Invoices) > 0 {
			firstInv := info1.Invoices[0]
			if !strings.HasPrefix(firstInv.InvoiceNumber, expectedInvNumPrefix) {
				t.Errorf("Record 1 Invoice 0: Expected InvoiceNumber starting with %q, got %q", expectedInvNumPrefix, firstInv.InvoiceNumber)
			}
		}
	}

	// Verify second record (AECOM)
	if len(eftInfos.EftInfos) > 1 {
		info2 := eftInfos.EftInfos[1]
		if info2.SupplierId != "19185" {
			t.Errorf("Record 2: Expected SupplierId '19185', got %q", info2.SupplierId)
		}
		// Verify invoices for second record
		// Blob has 1 invoice: (I)8753
		if len(info2.Invoices) != 1 {
			t.Errorf("Record 2: Expected 1 invoice, got %d", len(info2.Invoices))
		}
	}
}

func TestGetEftFromCSV_AP566005(t *testing.T) {
	// Locate the test file
	testFilePath := filepath.Join("..", "..", "testfiles", "AP566005_CRCRLF_present.csv")

	// Read the file
	content, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", testFilePath, err)
	}

	// Parse CSV
	eftInfos, err := GetEftInfosFromCSV(string(content))
	if err != nil {
		t.Fatalf("GetEftFromCSV returned error: %v", err)
	}

	// Expected: 3 records
	// 1. D0004
	// 2. 12889
	// 3. M1221
	if len(eftInfos.EftInfos) != 3 {
		t.Errorf("Expected 3 EftInfo records, got %d", len(eftInfos.EftInfos))
	}

	// Verify Record 1 (D0004)
	if len(eftInfos.EftInfos) > 0 {
		info := eftInfos.EftInfos[0]
		if info.SupplierId != "D0004" {
			t.Errorf("Record 1: Expected SupplierId 'D0004', got %q", info.SupplierId)
		}
		if len(info.Invoices) != 2 {
			t.Errorf("Record 1: Expected 2 invoices, got %d", len(info.Invoices))
		}
		// First invoice (I)8791 - but check if sorted?
		// 8791 < 8792. So order should be preserved.
		if len(info.Invoices) > 0 {
			if !strings.HasPrefix(info.Invoices[0].InvoiceNumber, "(I)8791") {
				t.Errorf("Record 1 Invoice 0: Expected prefix '(I)8791', got %q", info.Invoices[0].InvoiceNumber)
			}
		}
	}

	// Verify Record 2 (12889)
	if len(eftInfos.EftInfos) > 1 {
		info := eftInfos.EftInfos[1]
		if info.SupplierId != "12889" {
			t.Errorf("Record 2: Expected SupplierId '12889', got %q", info.SupplierId)
		}
		if len(info.Invoices) != 2 {
			t.Errorf("Record 2: Expected 2 invoices, got %d", len(info.Invoices))
		}
	}

	// Verify Record 3 (M1221)
	if len(eftInfos.EftInfos) > 2 {
		info := eftInfos.EftInfos[2]
		if info.SupplierId != "M1221" {
			t.Errorf("Record 3: Expected SupplierId 'M1221', got %q", info.SupplierId)
		}
		// Invoices: 6 lines in the cell.
		if len(info.Invoices) != 6 {
			t.Errorf("Record 3: Expected 6 invoices, got %d", len(info.Invoices))
		}
	}
}
