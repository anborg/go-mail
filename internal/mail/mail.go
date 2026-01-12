package mail

import (
	"crypto/tls"

	gomail "gopkg.in/mail.v2"
)

// Service provides email sending capabilities.
type Service struct {
	conf MailServerConfig
}

// NewService creates a new mail service.
func NewService(conf MailServerConfig) *Service {
	return &Service{conf: conf}
}


// Email sends an email based on EmailInfo.
func (s *Service) Email(mailinfo EmailInfo) error {
	dialer := gomail.NewDialer(s.conf.Host, s.conf.Port, s.conf.User, s.conf.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sendCloser, err := dialer.Dial()
	if err != nil {
		return err
	}
	defer sendCloser.Close()

	message := gomail.NewMessage()
	message.SetHeader("From", mailinfo.From)
	message.SetHeader("To", mailinfo.To)
	message.SetHeader("Cc", mailinfo.Cc)
	message.SetHeader("Subject", mailinfo.Subject)
	message.SetBody("text/html", mailinfo.Body)

	if err := gomail.Send(sendCloser, message); err != nil {
		return err
	}
	return nil
}
