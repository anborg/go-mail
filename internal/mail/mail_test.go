package mail

import (
	"fmt"
	"muni/go-mail/internal/config"
	"strings"
	"testing"
	"time"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
)

func TestErrorEmail(t *testing.T) {
	// 1. Start mock SMTP server
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       false,
		LogServerActivity: false,
	})

	if err := server.Start(); err != nil {
		t.Fatalf("Failed to start mock SMTP server: %v", err)
	}
	defer server.Stop()

	// 2. Prepare config pointing to mock server
	conf := config.MailServerConfig{
		Host:     "127.0.0.1",
		Port:     server.PortNumber(),
		User:     "testuser",
		Password: "testpassword",
		OpsUser:  "ops@example.com",
		CcUser:   "cc@example.com",
	}

	info := EmailInfo{
		From:    "sender@example.com",
		To:      "receiver@example.com",
		Cc:      "cc@example.com",
		Subject: "Test Subject",
		Body:    "<h1>Test Body</h1>",
	}

	// 3. Call the function
	err := ErrorEmail(conf, info)
	if err != nil {
		t.Fatalf("ErrorEmail failed: %v", err)
	}

	// 4. Verify mock server received the message
	messages, err := server.WaitForMessages(1, 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	msg := messages[0]
	// smtpmock V2 uses MailfromRequest() which returns something like "MAIL FROM:<sender@example.com>"
	if !strings.Contains(msg.MailfromRequest(), "sender@example.com") {
		t.Errorf("Expected MAIL FROM to contain sender@example.com, got %s", msg.MailfromRequest())
	}
}


func TestSendErrorAlert(t *testing.T) {
	server := smtpmock.New(smtpmock.ConfigurationAttr{})
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	conf := config.MailServerConfig{
		Host:    "127.0.0.1",
		Port:    server.PortNumber(),
		OpsUser: "ops@example.com",
		CcUser:  "cc@example.com",
	}

	err := SendErrorAlert(conf, "Alert Subject", "Alert Body")
	if err != nil {
		t.Errorf("SendErrorAlert failed: %v", err)
	}

	messages, err := server.WaitForMessages(1, 2*time.Second)
	if err != nil {
		t.Fatalf("Expected 1 message, got error: %v", err)
	}
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}
}

func TestErrorEmail_ConnectionError(t *testing.T) {
	// Point to a port that is likely closed or invalid
	conf := config.MailServerConfig{
		Host: "127.0.0.1",
		Port: 9999,
	}
	info := EmailInfo{From: "s@e.com", To: "r@e.com"}

	err := ErrorEmail(conf, info)
	if err == nil {
		t.Error("Expected error for invalid connection, got nil")
	} else {
		fmt.Printf("Expected connection error: %v\n", err)
	}
}
