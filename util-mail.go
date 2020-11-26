package main

import (
	"crypto/tls"
	"log"
	"time"

	gomail "gopkg.in/mail.v2"
)

//EmailInfo - convert eft into Email Data obj
type EmailInfo struct {
    From    string
	To      string
	Cc      string
	Subject string
	Body    string
}

func errorEmail(config MailServerConfig, mailinfo EmailInfo) error {
	dialer := gomail.NewDialer(config.Host, config.Port, config.User, config.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sendCloser, err := dialer.Dial()
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", mailinfo.From) //eftapp@
	message.SetHeader("To", mailinfo.To) // for errors this should be adminEmail, opsEmail
	message.SetHeader("Cc", mailinfo.Cc)
	message.SetHeader("Subject", "ERROR: Markham Notification - EFT")
	message.SetBody("text/html", mailinfo.Body) //TODO create body form template

	if err := gomail.Send(sendCloser, message); err != nil {
		//log.Printf("Could not send email to %q: %v", eftinfo.Email, err)
		return err
	}
	//message.Reset()
	return nil
}

func batchSendMail(config MailServerConfig, myarray EftInfos) error {
	// Settings for SMTP server
	dialer := gomail.NewDialer(config.Host, config.Port, config.User, config.Password)
	//log.Println("DialerCof: ", config.Host, config.Port, config.User, config.Password)
	// This is only needed when SSL/TLS certi

	//certificate is not valid on server.
	// In production this should be set to false.
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sendCloser, err := dialer.Dial()
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	for i := 0; i < len(myarray.EftInfos); i++ {
		message.Reset()
		var eftinfo EftInfo = myarray.EftInfos[i]
		eftinfo.TodayDate = time.Now().Format("2006-01-02 15:04:05 Monday")
		bytesHtml, err := ExecEftTemplate(eftinfo)
		if err != nil {
			return (err)
		}
		log.Println(i, eftinfo)
		// Set E-Mail sender
		message.SetHeader("From", config.CcUser)
		message.SetHeader("To", eftinfo.Email)
		message.SetHeader("Cc", config.CcUser)
		//message.SetHeader("Bcc", "office@example.com", "anotheroffice@example.com")
		//message.SetAddressHeader("Reply-To", "noreply@example.com")
		message.SetHeader("Subject", "Markham Notification - EFT")
		message.SetBody("text/html", bytesHtml.String()) //TODO create body form template

		// Attach some file
		//message.Attach("myfile1.pdf")

		if err := gomail.Send(sendCloser, message); err != nil {
			//log.Printf("Could not send email to %q: %v", eftinfo.Email, err)
			return err
		}
		//message.Reset()
	}

	return nil
} //batchSendMail
