package eftnotify

import (
    "crypto/tls"
    "log"
    "time"
"bytes"
	"muni/go-mail/internal/utils"
	"muni/go-mail/internal/config"
	gomail "gopkg.in/mail.v2"
)

func BatchSendMail(conf config.MailServerConfig, myarray EftInfos) error {
    dialer := gomail.NewDialer(conf.Host, conf.Port, conf.User, conf.Password)
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
        bytesHtml, err := generateEmailBody(eftinfo)
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

// GenerateEmailBody apply info on Template to create html
func generateEmailBody(eftinfo EftInfo) (bytes.Buffer, error) {
	eftTemplatePath := "templates/DEFAULT.gohtml"
	// 	eftTemplatePath := "templates/SAMPLE_TEMPLATE.txt" //for testing
	return utils.RenderTemplate(eftTemplatePath, eftinfo)
}