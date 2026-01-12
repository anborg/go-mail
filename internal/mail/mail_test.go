package mail

import (
	"fmt"
	"strings"
	"testing"
	"time"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
)

func TestEmail(t *testing.T) {
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
	conf := MailServerConfig{
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
	svc := NewService(conf)
	err := svc.Email(info)
	if err != nil {
		t.Fatalf("Email failed: %v", err)
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


func TestEmail_ConnectionError(t *testing.T) {
	// Point to a port that is likely closed or invalid
	conf := MailServerConfig{
		Host: "127.0.0.1",
		Port: 9999,
	}
	info := EmailInfo{From: "s@e.com", To: "r@e.com"}

	svc := NewService(conf)
	err := svc.Email(info)
	if err == nil {
		t.Error("Expected error for invalid connection, got nil")
	} else {
		fmt.Printf("Expected connection error: %v\n", err)
	}
}
