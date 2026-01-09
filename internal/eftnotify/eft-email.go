package eftnotify

import (
	"bytes"
	"crypto/tls"
	"log"
	"muni/go-mail/internal/config"
	"muni/go-mail/internal/utils"
	"sync"
	"time"

	gomail "gopkg.in/mail.v2"
)

var EftTemplatePath = "templates/DEFAULT.gohtml"

// BatchSendMail sends emails to multiple suppliers concurrently.
func BatchSendMail(conf config.MailServerConfig, myarray EftInfos) error {
	dialer := gomail.NewDialer(conf.Host, conf.Port, conf.User, conf.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	var wg sync.WaitGroup
	errs := make(chan error, len(myarray.EftInfos))
	
	// Use a worker pool to limit concurrency if needed. 
	// For now, we'll just use a simple limit or go full parallel if the batch is small.
	maxWorkers := 5
	jobs := make(chan EftInfo, len(myarray.EftInfos))

	// Start workers
	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for eftinfo := range jobs {
				if err := sendSingleMail(dialer, conf, eftinfo); err != nil {
					errs <- err
				}
			}
		}()
	}

	// Send jobs
	for i := 0; i < len(myarray.EftInfos); i++ {
		eftinfo := myarray.EftInfos[i]
		eftinfo.TodayDate = time.Now().Format("2006-01-02 15:04:05 Monday")
		jobs <- eftinfo
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	close(errs)

	// Collect first error if any
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func sendSingleMail(dialer *gomail.Dialer, conf config.MailServerConfig, eftinfo EftInfo) error {
	message := gomail.NewMessage()
	bytesHtml, err := generateEmailBody(eftinfo)
	if err != nil {
		return err
	}
	
	log.Println("Sending email to:", eftinfo.Email)
	message.SetHeader("From", conf.CcUser)
	message.SetHeader("To", eftinfo.Email)
	message.SetHeader("Cc", conf.CcUser)
	message.SetHeader("Subject", "Markham Notification - EFT")
	message.SetBody("text/html", bytesHtml.String())

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

// generateEmailBody apply info on Template to create html
func generateEmailBody(eftinfo EftInfo) (bytes.Buffer, error) {
	return utils.RenderTemplate(EftTemplatePath, eftinfo)
}