package eftnotify

import (
	"muni/go-mail/internal/config"
	"os"
	"testing"
	"time"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
)

func TestBatchSendMail(t *testing.T) {
	// 1. Start mock SMTP server
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		MultipleMessageReceiving: true,
		LogToStdout:              true,
		LogServerActivity:        true,
	})
	if err := server.Start(); err != nil {
		t.Fatalf("Failed to start mock SMTP server: %v", err)
	}
	defer server.Stop()

	// 2. Prepare mock template
	content := "Supplier: {{.SupplierName}}"
	tmpfile, err := os.CreateTemp("", "test_template_*.gohtml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(content))
	tmpfile.Close()

	// Override global template path for test
	oldPath := EftTemplatePath
	EftTemplatePath = tmpfile.Name()
	defer func() { EftTemplatePath = oldPath }()

	// 3. Prepare mock data
	conf := config.MailServerConfig{
		Host:   "127.0.0.1",
		Port:   server.PortNumber(),
		CcUser: "cc@example.com",
	}

	eftInfos := EftInfos{
		EftInfos: []EftInfo{
			{
				SupplierName: "Supplier 1",
				Email:        "s1@example.com",
				Invoices: []Invoice{
					{InvoiceNumber: "INV1", Amount: "$100"},
				},
			},
			{
				SupplierName: "Supplier 2",
				Email:        "s2@example.com",
				Invoices: []Invoice{
					{InvoiceNumber: "INV2", Amount: "$200"},
				},
			},
		},
	}

	// 4. Call BatchSendMail
	err = BatchSendMail(conf, eftInfos)
	if err != nil {
		t.Fatalf("BatchSendMail failed: %v", err)
	}

	// 5. Verify mock server received both messages
	messages, err := server.WaitForMessages(2, 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to receive 2 messages: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}
}
