package mail

import (
	"crypto/tls"
	"muni/go-mail/internal/config"

	gomail "gopkg.in/mail.v2"
)



func SendErrorAlert(mailConf config.MailServerConfig, subject string, body string) error {
	emailInfo := EmailInfo{
		From:    mailConf.OpsUser,
		To:      mailConf.OpsUser,
		Cc:      mailConf.CcUser,
		Subject: subject,
		Body:    body,
	}
	err := ErrorEmail(mailConf, emailInfo)
	return err
}

func ErrorEmail(conf config.MailServerConfig, mailinfo EmailInfo) error {
	dialer := gomail.NewDialer(conf.Host, conf.Port, conf.User, conf.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sendCloser, err := dialer.Dial()
	if err != nil {
		return err
	}

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


