package eftnotify

import (
	"bytes"
	"errors"
	"muni/go-mail/internal/config"
	"muni/go-mail/internal/mail"
	"muni/go-mail/internal/utils"
	"sync"
	"time"
)

var EftTemplatePath = "templates/DEFAULT.gohtml"

// Notifier defines the interface for all email communication for the eftnotify package.
type Notifier interface {
	NotifyAll(EftInfos) error
	Error(subject, body string) error
}

// EftNotifier handles email communication using the mail service.
type EftNotifier struct {
	mailService *mail.Service
	conf        config.MailServerConfig
}

// NewNotifier creates a new EftNotifier.
func NewNotifier(mailService *mail.Service, conf config.MailServerConfig) Notifier {
	return &EftNotifier{
		mailService: mailService,
		conf:        conf,
	}
}

// Error sends an administrative error alert (routed through mail service).
func (n *EftNotifier) Error(subject, body string) error {
	return n.mailService.SendErrorAlert(subject, body)
}

// NotifyAll sends emails to multiple suppliers concurrently.
func (n *EftNotifier) NotifyAll(myarray EftInfos) error {
	var wg sync.WaitGroup
	errs_channel := make(chan error, len(myarray.EftInfos))
	maxWorkers := 5
	eft_jobs_channel := make(chan EftInfo, len(myarray.EftInfos))

	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for eftinfo := range eft_jobs_channel {
				if err := n.notify(eftinfo); err != nil {
					errs_channel <- err
				}
			}
		}()
	}

	nowStr := time.Now().Format("2006-01-02 15:04:05 Monday")
	for i := 0; i < len(myarray.EftInfos); i++ {
		eftinfo := myarray.EftInfos[i]
		eftinfo.TodayDate = nowStr
		eft_jobs_channel <- eftinfo
	}
	close(eft_jobs_channel)

	wg.Wait()
	close(errs_channel)

	var allErrs []error
	for err := range errs_channel {
		if err != nil {
			allErrs = append(allErrs, err)
		}
	}
	if len(allErrs) > 0 {
		return errors.Join(allErrs...)
	}
	return nil
}

// notify sends a single EFT notification email.
func (n *EftNotifier) notify(eftinfo EftInfo) error {
	bytesHtml, err := generateEmailBody(eftinfo)
	if err != nil {
		return err
	}

	email := mail.EmailInfo{
		From:    n.conf.CcUser,
		To:      eftinfo.Email,
		Cc:      n.conf.CcUser,
		Subject: "Markham Notification - EFT",
		Body:    bytesHtml.String(),
	}
	
	return n.mailService.Email(email)
}

// generateEmailBody applies info on Template to create html
func generateEmailBody(eftinfo EftInfo) (bytes.Buffer, error) {
	return utils.RenderTemplate(EftTemplatePath, eftinfo)
}