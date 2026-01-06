package mail

import (
	"crypto/tls"
	"log"
	"time"

	"muni/go-mail/internal/config"
	"muni/go-mail/internal/model"

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

func SendErrorAlert(mailConf config.MailServerConfig, msg string) error {
	emailInfo := EmailInfo{
		From:    mailConf.OpsUser,
		To:      mailConf.OpsUser,
		Cc:      mailConf.CcUser,
		Subject: "Error: Markham Notification - EFT",
		Body:    "Error while processing eft file : \n\n" + msg,
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
	message.SetHeader("Subject", "ERROR: Markham Notification - EFT")
	message.SetBody("text/html", mailinfo.Body)

	if err := gomail.Send(sendCloser, message); err != nil {
		return err
	}
	return nil
}

func BatchSendMail(conf config.MailServerConfig, myarray model.EftInfos) error {
	dialer := gomail.NewDialer(conf.Host, conf.Port, conf.User, conf.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sendCloser, err := dialer.Dial()
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	for i := 0; i < len(myarray.EftInfos); i++ {
		message.Reset()
		var eftinfo model.EftInfo = myarray.EftInfos[i]
		eftinfo.TodayDate = time.Now().Format("2006-01-02 15:04:05 Monday")
		bytesHtml, err := ExecEftTemplate(eftinfo)
		if err != nil {
			return (err)
		}
		log.Println(i, eftinfo)
		message.SetHeader("From", conf.CcUser)
		message.SetHeader("To", eftinfo.Email)
		message.SetHeader("Cc", conf.CcUser)
		message.SetHeader("Subject", "Markham Notification - EFT")
		message.SetBody("text/html", bytesHtml.String())

		if err := gomail.Send(sendCloser, message); err != nil {
			return err
		}
	}

	return nil
}
