package main

import (
	"crypto/tls"
	"log"
	"time"

	gomail "gopkg.in/mail.v2"
)

//EmailInfo - convert eft into Email Data obj
type EmailInfo struct {
	To      string
	Cc      string
	Subject string
	Body    string
}

func batchSendMail(config MailServerConfig, myarray EftInfos) error {
	// Settings for SMTP server
	dialer := gomail.NewDialer(config.Host, config.Port, config.User, config.Password)

	// This is only needed when SSL/TLS certi

	//certificate is not valid on server.
	// In production this should be set to false.
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	s, err := dialer.Dial()
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	for i := 0; i < len(myarray.EftInfos); i++ {
		m.Reset()
		var eftinfo EftInfo = myarray.EftInfos[i]
		eftinfo.TodayDate = time.Now().Format("2006-01-02 15:04:05 Monday")
		bytesHtml, err := ExecEftTemplate(eftinfo)
		if err != nil {
			return (err)
		}
		log.Println(i, eftinfo)
		// Set E-Mail sender
		m.SetHeader("From", config.User)
		m.SetHeader("To", eftinfo.Email)
		m.SetHeader("Cc", config.User)
		//m.SetHeader("Bcc", "office@example.com", "anotheroffice@example.com")
		//m.SetAddressHeader("Reply-To", "noreply@example.com")
		m.SetHeader("Subject", "Markham Notification - EFT")
		m.SetBody("text/html", bytesHtml.String()) //TODO create body form template

		// Attach some file
		//m.Attach("myfile1.pdf")

		if err := gomail.Send(s, m); err != nil {
			//log.Printf("Could not send email to %q: %v", eftinfo.Email, err)
			return err
		}
		//m.Reset()
	}

	return nil
} //batchSendMail
